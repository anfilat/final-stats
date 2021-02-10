package common

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

const procPath = "/proc"

func ReadProcFile(name string) ([]string, error) {
	content, err := ioutil.ReadFile(procFileName(name))
	if err != nil {
		return nil, err
	}

	return SplitLines(string(content)), nil
}

func SplitLines(content string) []string {
	return strings.Split(strings.TrimSpace(content), "\n")
}

func procFileName(name string) string {
	return filepath.Join(procPath, name)
}
