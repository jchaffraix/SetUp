package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var verbose bool

func runCommandInteractively(args []string) ([]byte, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if verbose {
		fmt.Println("Running `" + strings.Join(args, " ") + "`: " + string(output))
	}
	return output, err
}

// TODO: This belongs to a git package.
func getCurrentBranch() string {
	args := []string{"git", "branch", "--show-current"}
	output, err := runCommandInteractively(args)
	// TODO: Handle detached HEAD?
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func checkoutMaster() error {
	args := []string{"git", "checkout", "master"}
	_, err := runCommandInteractively(args)
	return err
}

func updateMasterIfNeeded() error {
	// TODO: Checkout last refs and check its timestamp.
	args := []string{"git", "pull", "-r"}
	_, err := runCommandInteractively(args)
	return err
}

func rebaseOntoMaster(branch string) error {
	args := []string{"git", "checkout", branch}
	_, err := runCommandInteractively(args)
	if err != nil {
		return err
	}
	args = []string{"git", "rebase", "master"}
	_, err = runCommandInteractively(args)
	return err
}

func main() {
	flag.BoolVar(&verbose, "verbose", false, "toggle verbose mode")
	flag.Parse()

	branch := getCurrentBranch()
	if branch == "" {
		fmt.Println("💣 Couldn't determine current branch... Aborting!")
		return
	}
	fmt.Println("Git is on branch: " + branch)
	isOnMaster := branch == "master"
	if !isOnMaster {
		checkoutMaster()
	}
	updateMasterIfNeeded()
	if !isOnMaster {
		// TODO: Do I want to rebase every local branches?
		rebaseOntoMaster(branch)
	}
}
