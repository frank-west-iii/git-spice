package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"go.abhg.dev/gs/internal/git"
	"go.abhg.dev/gs/internal/state"
	"go.abhg.dev/gs/internal/text"
)

type branchOntoCmd struct {
	Branch string `help:"Branch to move" placeholder:"NAME"`
	Onto   string `arg:"" optional:"" help:"Destination branch"`
}

func (*branchOntoCmd) Help() string {
	return text.Dedent(`
		Transplants the commits of a branch on top of another branch
		without picking up any changes from the old base.
		The base for the branch will be updated to the new branch.
	`)
}

func (cmd *branchOntoCmd) Run(ctx context.Context, log *log.Logger, opts *globalOptions) error {
	repo, err := git.Open(ctx, ".", git.OpenOptions{
		Log: log,
	})
	if err != nil {
		return fmt.Errorf("open repository: %w", err)
	}

	if cmd.Branch == "" {
		currentBranch, err := repo.CurrentBranch(ctx)
		if err != nil {
			return fmt.Errorf("get current branch: %w", err)
		}
		cmd.Branch = currentBranch
	}

	// TODO: fuzzy prompt for Onto if unset
	if cmd.Onto == "" {
		return fmt.Errorf("destination branch required")
	}

	ontoHash, err := repo.PeelToCommit(ctx, cmd.Onto)
	if err != nil {
		return fmt.Errorf("peel to commit: %w", err)
	}

	store, err := ensureStore(ctx, repo, log, opts)
	if err != nil {
		return err
	}

	if cmd.Branch == store.Trunk() {
		return fmt.Errorf("cannot move trunk")
	}

	b, err := store.Lookup(ctx, cmd.Branch)
	if err != nil {
		if errors.Is(err, state.ErrNotExist) {
			return fmt.Errorf("branch not tracked: %s", cmd.Branch)
		}
		return fmt.Errorf("get branch: %w", err)
	}

	if b.Base == cmd.Onto {
		log.Infof("%s: already on %s", cmd.Branch, cmd.Onto)
		return nil
	}

	if err := repo.Rebase(ctx, git.RebaseRequest{
		Branch:    cmd.Branch,
		Upstream:  b.Base, // TODO: use fork point?
		Onto:      cmd.Onto,
		Autostash: true,
		Quiet:     true,
	}); err != nil {
		return fmt.Errorf("rebase: %w", err)
	}

	// TODO: handle conflicts/partial rebase

	err = store.Update(ctx, &state.UpdateRequest{
		Upserts: []state.UpsertRequest{
			{
				Name:     cmd.Branch,
				Base:     cmd.Onto,
				BaseHash: ontoHash,
			},
		},
		Message: fmt.Sprintf("move %s onto %s", cmd.Branch, cmd.Onto),
	})
	if err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	return nil
}