package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/ncruces/rethinkraw/internal/config"
	"golang.org/x/sync/semaphore"
)

var semDNGConverter = semaphore.NewWeighted(3)

func testDNGConverter() error {
	_, err := os.Stat(config.DngConverter)
	return err
}

func runDNGConverter(input, output string, side int, exp *exportSettings) error {
	err := os.RemoveAll(output)
	if err != nil {
		return err
	}

	dir := filepath.Dir(output)
	output = filepath.Base(output)

	opts := []string{}
	if exp != nil && exp.DNG {
		if exp.Preview != "" {
			opts = append(opts, "-"+exp.Preview)
		}
		if exp.Lossy {
			opts = append(opts, "-lossy")
		}
		if exp.Embed {
			opts = append(opts, "-e")
		}
	} else {
		if side > 0 {
			opts = append(opts, "-lossy", "-side", strconv.Itoa(side))
		}
		opts = append(opts, "-p2")
	}
	opts = append(opts, "-d", dir, "-o", output, input)

	if err := semDNGConverter.Acquire(context.TODO(), 1); err != nil {
		return err
	}
	defer semDNGConverter.Release(1)

	log.Print("dng converter...")
	cmd := exec.Command(config.DngConverter, opts...)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("DNG Converter: %w", err)
	}
	return nil
}
