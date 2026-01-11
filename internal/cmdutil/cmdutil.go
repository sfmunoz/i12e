package cmdutil

import (
	"bufio"
	"io"
	"os/exec"
	"sync"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "cmdutil")

func logOutput(out io.ReadCloser, prefix string) {
	s := bufio.NewScanner(out)
	for s.Scan() {
		line := s.Text()
		log.Info(prefix + line)
	}
	if err := s.Err(); err != nil {
		log.Error("'scanner.Err()' is not nil", "err", err)
	}
}

func RunCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		logOutput(stdout, "o> ")
	}()
	go func() {
		defer wg.Done()
		logOutput(stderr, "e> ")
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
