//go:build !release

package main

import (
	"github.com/darabuchi/log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}
