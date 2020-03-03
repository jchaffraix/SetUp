package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Global constants
var deps = []string{"git", "tmux", "zsh"}
var githubURL = "git@github.com:jchaffraix/SetUp.git"

// TODO: I should use go-git to fetch git.

func installSoftwareDeps() error {
	fmt.Println("✨ Installing deps")
	switch runtime.GOOS {
	case "linux":
		args := []string{"sudo", "apt-get", "install"}
		args = append(args, deps...)
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	case "darwin":
		args := []string{"brew", "install"}
		args = append(args, deps...)
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	case "Windows":
		return errors.New("Can't install missing deps on Windows.")
	default:
		return errors.New("Unknown OS: " + runtime.GOOS)
	}
	panic("Missing return in installSoftwareDeps, os = " + runtime.GOOS)
}

func main() {
	if err := installSoftwareDeps(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
