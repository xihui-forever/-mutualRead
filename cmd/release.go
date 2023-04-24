//go:build release

package main

import "github.com/darabuchi/log"

func init() {
	log.SetOutput(log.GetOutputWriterHourly("/tmp/mutualRead/", 3))
}
