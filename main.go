package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	repoURL, branch := flag.Arg(0), flag.Arg(1)

	repo, err := git.PlainOpen(".")
	if err == git.ErrRepositoryNotExists {
		repo, err = git.PlainInit(".", false)
	}
	handleErr(err, "error opening/initializing repo")

	remote, err := repo.Remote("origin")
	if err == git.ErrRemoteNotFound {
		remote, err = repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{repoURL}})
	} else if err == nil && remote.Config().URLs[0] != repoURL {
		err = fmt.Errorf("remote URL mismatch: %s != %s", remote.Config().URLs[0], repoURL)
	}
	handleErr(err, "error getting remote")

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", remote.Config().Name)),
		},
		Force:    true,
		Progress: os.Stderr,
	})
	if err == git.NoErrAlreadyUpToDate {
		err = nil
	}
	handleErr(err, "error fetching")

	wt, err := repo.Worktree()
	handleErr(err, "error getting worktree")

	branchRef := plumbing.NewRemoteReferenceName("origin", branch)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRef),
		Force:  true,
	})
	handleErr(err, "error checking out branch")

	err = wt.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: plumbing.NewHash("origin/" + branch),
	})
	handleErr(err, "error resetting to remote")
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <repository-url> <branch-name>\n", os.Args[0])
		flag.PrintDefaults()
	}
}
