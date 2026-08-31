package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
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

func currentVersion() (string, error) {
	out, err := runCommandInteractively([]string{"code", "--version"})
	if err != nil {
		return "", err
	}
	lines := bytes.Split(out, []byte{'\n'})
	return string(lines[0]), nil
}

func fetchLatestVersion() (string, error) {
	req, err := http.NewRequest("GET", "https://update.code.visualstudio.com/api/releases/stable", nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if verbose {
		fmt.Println("Got body=" + string(body))
	}

	var versions []string
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", errors.New("No versions")
	}
	return versions[0], nil
}

func downloadVSCode(latestVersion string) (string, error) {
	url := "https://update.code.visualstudio.com/" + latestVersion + "/linux-deb-x64/stable"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", "code*.deb")
	if err != nil {
		return "", err
	}
	// We can't defer os.Remove here as we want to return it to the caller.
	// This means that any error path until the end of the function must manually call it!
	_, err = f.Write(body)
	if err != nil {
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func installVSCode(fName string) error {
	_, err := runCommandInteractively([]string{"sudo", "dpkg", "-i", fName})
	return err
}

func main() {
	flag.BoolVar(&verbose, "verbose", false, "toggle verbose mode")
	flag.Parse()

	curVersion, err := currentVersion()
	if err != nil {
		fmt.Println("Error getting current version (" + err.Error() + ")")
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Currently installed version: " + curVersion)
	}

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		fmt.Println("Error when fetching the latest version (" + err.Error() + ")")
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Latest version: " + latestVersion)
	}

	if strings.Compare(curVersion, latestVersion) > 1 {
		fName, err := downloadVSCode(latestVersion)
		defer os.Remove(fName)
		if err != nil {
			fmt.Println("Error when downloading the latest version (" + err.Error() + ")")
			os.Exit(1)
		}
		if err = installVSCode(fName); err != nil {
			fmt.Println("Error installing the latest version (" + err.Error() + ")")
			os.Exit(1)
		}
	} else {
		fmt.Println("No new version to fetch...")
	}
}
