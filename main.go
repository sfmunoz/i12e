package main

import (
	"os"
	"slices"
	"strings"

	"github.com/sfmunoz/i12e/cmd"
)

var (
	Version   = "???"
	CommitSHA = "???"
	BuildTime = "???"
)

// https://github.com/sfmunoz/i12e/issues/239
func addOptBinToPath() {
	sep := string(os.PathListSeparator)
	optBin := "/opt/bin"
	path := strings.TrimSpace(os.Getenv("PATH"))
	if len(path) < 1 {
		os.Setenv("PATH", optBin)
		return
	}
	items := strings.Split(path, sep)
	if slices.Contains(items, optBin) {
		return
	}
	itemsNew := append([]string{optBin}, items...)
	pathNew := strings.Join(itemsNew, ":")
	os.Setenv("PATH", pathNew)
}

func main() {
	addOptBinToPath()
	cmd.Execute(Version, CommitSHA, BuildTime)
}
