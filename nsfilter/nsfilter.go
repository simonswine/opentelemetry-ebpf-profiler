// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package nsfilter provides utilities for working with PID namespaces.
package nsfilter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetCurrentPIDNamespace returns the PID namespace inode number for the current process.
// It reads /proc/self/ns/pid and parses the inode number from the format "pid:[inode]".
func GetCurrentPIDNamespace() (uint64, error) {
	return getPIDNamespaceFromPath("/proc/self/ns/pid")
}

// GetPIDNamespace returns the PID namespace inode number for the specified PID.
// It reads /proc/<pid>/ns/pid and parses the inode number from the format "pid:[inode]".
func GetPIDNamespace(pid int) (uint64, error) {
	path := fmt.Sprintf("/proc/%d/ns/pid", pid)
	return getPIDNamespaceFromPath(path)
}

// getPIDNamespaceFromPath reads a namespace file and extracts the inode number.
// The namespace file is a symbolic link with format "pid:[inode]".
func getPIDNamespaceFromPath(path string) (uint64, error) {
	// Read the symlink target which has format "pid:[inode]"
	link, err := os.Readlink(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read namespace symlink %s: %w", path, err)
	}

	// Parse the inode number from "pid:[inode]" format
	// Expected format: "pid:[123456]"
	if !strings.HasPrefix(link, "pid:[") || !strings.HasSuffix(link, "]") {
		return 0, fmt.Errorf("unexpected namespace format: %s", link)
	}

	// Extract the inode number between '[' and ']'
	inodeStr := link[5 : len(link)-1]
	inode, err := strconv.ParseUint(inodeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse namespace inode from %s: %w", link, err)
	}

	return inode, nil
}
