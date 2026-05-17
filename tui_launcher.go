package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

const tuiBinaryName = "go-automate-tui"

// isInteractiveTerminal reports whether stdin is connected to a terminal.
func isInteractiveTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// findTUIBinary looks for the TUI binary, first alongside the running
// executable, then via PATH lookup.
func findTUIBinary() (string, bool) {
	// Check alongside the running binary
	exe, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exe), tuiBinaryName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
	}

	// Fall back to PATH
	path, err := exec.LookPath(tuiBinaryName)
	if err == nil {
		return path, true
	}

	return "", false
}

// launchTUI finds and runs the TUI binary with inherited stdio.
// It does not return on success (the TUI replaces this process's stdio).
func launchTUI() error {
	tuiPath, found := findTUIBinary()
	if !found {
		return fmt.Errorf("%s binary not found (looked beside executable and in PATH)", tuiBinaryName)
	}

	log.Infof("Launching TUI: %s", tuiPath)

	cmd := exec.Command(tuiPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	os.Exit(0)
	return nil // unreachable
}

// maybeLaunchTUI checks if the session is interactive and a TUI binary
// is available. Returns true if the TUI was launched (caller should exit).
func maybeLaunchTUI() bool {
	if !isInteractiveTerminal() {
		return false
	}

	tuiPath, found := findTUIBinary()
	if !found {
		log.Infof("Interactive terminal detected but %s not found, continuing with CLI", tuiBinaryName)
		return false
	}

	log.Infof("Interactive terminal detected, launching TUI: %s", tuiPath)

	cmd := exec.Command(tuiPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		log.Warnf("TUI failed: %v, falling back to CLI mode", err)
		return false
	}

	os.Exit(0)
	return true // unreachable
}

// terminalEmulators lists terminal emulators to try, in preference order.
var terminalEmulators = []struct {
	name string
	args func(tuiPath string) []string
}{
	{"foot", func(tuiPath string) []string { return []string{"foot", "-e", tuiPath} }},
	{"alacritty", func(tuiPath string) []string { return []string{"alacritty", "-e", tuiPath} }},
	{"kitty", func(tuiPath string) []string { return []string{"kitty", tuiPath} }},
	{"ghostty", func(tuiPath string) []string { return []string{"ghostty", "-e", tuiPath} }},
	{"xterm", func(tuiPath string) []string { return []string{"xterm", "-e", tuiPath} }},
}

// launchTUIInTerminal opens the TUI binary inside a terminal emulator.
// It checks $TERMINAL first, then falls back to known terminal emulators.
func launchTUIInTerminal(tuiPath string) {
	// Try $TERMINAL environment variable first
	if terminal := os.Getenv("TERMINAL"); terminal != "" {
		if termPath, err := exec.LookPath(terminal); err == nil {
			cmd := exec.Command(termPath, "-e", tuiPath)
			if err := cmd.Start(); err != nil {
				log.Warnf("Failed to launch TUI in $TERMINAL (%s): %v", terminal, err)
			} else {
				log.Infof("Launched TUI in $TERMINAL: %s", terminal)
				return
			}
		}
	}

	// Try known terminal emulators
	for _, te := range terminalEmulators {
		if _, err := exec.LookPath(te.name); err != nil {
			continue
		}
		args := te.args(tuiPath)
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Start(); err != nil {
			log.Warnf("Failed to launch TUI in %s: %v", te.name, err)
			continue
		}
		log.Infof("Launched TUI in %s", te.name)
		return
	}

	log.Warn("No terminal emulator found to launch TUI")
}
