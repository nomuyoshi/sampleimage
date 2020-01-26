package main

import (
	"os"

	sampleimage "github.com/nomuyoshi/sampleimage/lib"
)

func main() {
	cli := &sampleimage.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	os.Exit(cli.Run(os.Args))
}
