package common

import (
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
)

const procPath = "/proc"

func ReadProcFile(name string) ([]string, error) {
	content, err := ioutil.ReadFile(procFileName(name))
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

func procFileName(name string) string {
	return filepath.Join(procPath, name)
}

func NumToFix2(num float64) float64 {
	return math.Round(num*100) / 100
}
