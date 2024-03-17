package main

import (
	"bufio"
	"errors"
	"fmt"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Configurable constants (flags).
var setupPath string
var verbose bool

// Global constants
var deps = []string{"git", "tmux", "vim", "zsh"}
const githubURL string = "git@github.com:jchaffraix/SetUp.git"
const vimSensibleURL string = "https://tpope.io/vim/sensible.git"

func runCommandInteractively(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func installConfigFile(homePath string, relFilePath []string) error {
	file := relFilePath[len(relFilePath)-1]
	relFilePathStr := filepath.Join(relFilePath...)
	destinationPath := filepath.Join(homePath, "." + file)
	fmt.Println("Installing: ", destinationPath)
	if pathExists(destinationPath) {
		// File exists, give users options.
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("File exist " + destinationPath + ": Overwrite/Skip/Exit [ose]: ")
			opt, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			switch opt[0] {
			case 'e':
				fallthrough
			case 'E':
				fmt.Println("💥 Exiting...")
				os.Exit(-1)
			case 's':
				fallthrough
			case 'S':
				fmt.Println("🚨 Skipping " + destinationPath)
				return nil
			case 'o':
				fallthrough
			case 'O':
				err = os.Rename(destinationPath, destinationPath + ".bak")
				if err != nil {
					return err
				}
				return os.Symlink(relFilePathStr, destinationPath)
			default:
				fmt.Println("🚩 Unknown input. Try again")
			}
		}
	}
	return os.Symlink(relFilePathStr, destinationPath)
}

func cloneConfig(homePath, relPath string) error {
	clonePath := filepath.Join(homePath, relPath)
	_, err := os.Lstat(clonePath)
	if err == nil {
		// TODO: We should probably give an option instead of assuming it is fine.
		fmt.Println("🚨 " + clonePath + " is not empty, skipping git clone.")
		return nil
	}
	args := []string{"git", "clone", githubURL, clonePath}
	if err := runCommandInteractively(args); err != nil {
		return err
	}
	return nil
}

// This script doesn't use go-git as it imports a lot of extra cruft.
func installConfigFiles(homePath, relPath string) error {
	fmt.Println("⚙️  Cloning the configs")
	err := cloneConfig(homePath, relPath)
	if err != nil {
		return err
	}

	fmt.Println("🚀 Installing the configs")
	relConfigPath := []string{relPath, "Configs"}
	absFilePath := append([]string{homePath}, relConfigPath...)
	dir, err := os.Open(filepath.Join(absFilePath...))
	if err != nil {
		return err
	}
	// -1 makes Readdir read all the entries in the directory.
	fileInfos, err := dir.Readdir(-1)
	dir.Close()
	if err != nil {
		return err
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue;
		}
		// |fileInfo| is just the last part of the path.
		// So we reconstruct the path relative to |homePath|.
		file := append(relConfigPath, fileInfo.Name())
		installConfigFile(homePath, file)
	}
	return nil
}

func installSoftwareDeps() error {
	fmt.Println("✨ Installing deps")
	switch runtime.GOOS {
	case "linux":
		args := []string{"sudo", "apt-get", "install", "-y"}
		args = append(args, deps...)
		return runCommandInteractively(args)
	case "darwin":
		args := []string{"brew", "install"}
		args = append(args, deps...)
		return runCommandInteractively(args)
	case "Windows":
		return errors.New("Can't install missing deps on Windows.")
	default:
		return errors.New("Unknown OS: " + runtime.GOOS)
	}
	panic("Missing return in installSoftwareDeps, os = " + runtime.GOOS)
}

func installVimPlugins(homePath, setupPath string) error {
	path := []string{homePath, ".vim", "pack", "tpope", "start"}
	pathStr := filepath.Join(path...)
	// Sanity check: Does the path exists?
	// This is required as MkdirAll doesn't return an error
	// if it does, but git will.
	if pathExists(pathStr) {
		fmt.Println("Skipping vim-sensible as path exists")
		return nil
	}
	if err := os.MkdirAll(pathStr, 0755); err != nil {
		return err
	}
	args := []string{"git", "clone", vimSensibleURL, pathStr}
	if err := runCommandInteractively(args); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.BoolVar(&verbose, "verbose", false, "toggle verbose mode")
	flag.StringVar(&setupPath, "setup_path", "Projects/SetUp", "installation path")
	flag.Parse()

	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get $HOME directory: " + err.Error() + "\n")
		os.Exit(1)
	}

	if err := installSoftwareDeps(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error() + "\n")
		os.Exit(1)
	}
	if err := installConfigFiles(homePath, setupPath); err != nil {
		fmt.Fprintf(os.Stderr, err.Error() + "\n")
		os.Exit(1)
	}
	if err := installVimPlugins(homePath, setupPath); err != nil {
		fmt.Fprintf(os.Stderr, err.Error() + "\n")
		os.Exit(1)
	}
	fmt.Println("✅ Install successful")
}
