package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

const (
	// LabelNodeRoleMaster specifies that a node is a master
	LabelNodeRoleMaster = "node-role.kubernetes.io/master"
)

// GetHostname returns OS's hostname if 'hostnameOverride' is empty; otherwise, return 'hostnameOverride'.
func GetHostname(hostnameOverride string) (string, error) {
	hostName := hostnameOverride
	if len(hostName) == 0 {
		nodeName, err := os.Hostname()
		if err != nil {
			return "", fmt.Errorf("couldn't determine hostname: %v", err)
		}
		hostName = nodeName
	}

	// Trim whitespaces first to avoid getting an empty hostname
	// For linux, the hostname is read from file /proc/sys/kernel/hostname directly
	hostName = strings.TrimSpace(hostName)
	if len(hostName) == 0 {
		return "", fmt.Errorf("empty hostname is invalid")
	}
	return strings.ToLower(hostName), nil
}

// ListK8sNodes returns k8s nodes base on node labels
func ListK8sNodes(kubectlPath, labels string) ([]string, error) {
	commandSlice := []string{
		kubectlPath,
		"get",
		"no",
		"--no-headers",
		"-ocustom-columns=:.metadata.name",
	}
	if labels != "" {
		commandSlice = append(commandSlice, "-l", labels)
	}
	commandStr := strings.Join(commandSlice, " ")
	cmd := exec.Command("/bin/sh", "-c", commandStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("exec: [%s] failed, output: %s, error: %v", commandStr, string(output), err)
	}

	nodes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(nodes) == 0 {
		return nil, fmt.Errorf("get k8s nodes is empty")
	}
	glog.Infof("get k8s nodes success: %s, labels: %s", nodes, labels)
	return nodes, nil
}

// EnsureDirectoryExist create directory if does not exist
func EnsureDirectoryExist(dirName string) error {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, os.ModePerm)
		if errDir != nil {
			return fmt.Errorf("create dir %s failed. err: %v", dirName, err)
		}
		return nil
	}

	if src.Mode().IsRegular() {
		return fmt.Errorf("%s already exist as a file", dirName)
	}

	return nil
}
