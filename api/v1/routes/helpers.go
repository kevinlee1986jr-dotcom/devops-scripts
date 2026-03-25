package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(logrus.InfoLevel)
}

// ExecuteCommand executes a shell command and returns the output and error.
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var stdout strings.Builder
	var stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Infof("Executing command: %s %s", command, strings.Join(args, " "))

	err := cmd.Run()

	if err != nil {
		logger.Errorf("Command failed: %s %s, Error: %v, Stdout: %s, Stderr: %s", command, strings.Join(args, " "), err, stdout.String(), stderr.String())
		return stdout.String(), fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
	}

	logger.Debugf("Command succeeded: %s %s, Stdout: %s", command, strings.Join(args, " "), stdout.String())

	return stdout.String(), nil
}

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir() && err == nil
}

// DirectoryExists checks if a directory exists
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir() && err == nil
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// WaitForProcessCompletion waits for a process to complete with a timeout.
func WaitForProcessCompletion(cmd *exec.Cmd, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		// Process timed out, kill it
		if cmd.Process != nil {
			logger.Warnf("Process timed out, killing pid %d", cmd.Process.Pid)
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil { // Kill process group
				logger.Errorf("Failed to kill process group: %v", err)
				if err := cmd.Process.Kill(); err != nil { // Fallback to killing the process directly
					logger.Errorf("Failed to kill process: %v", err)
					return fmt.Errorf("process timed out and could not be killed: %w", err)
				}
			}
			return fmt.Errorf("process timed out")
		} else {
			return fmt.Errorf("process timed out but no process found to kill")
		}
	case err := <-done:
		if err != nil {
			return fmt.Errorf("process exited with error: %w", err)
		}
		return nil
	}
}