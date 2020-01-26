package main

import (
	"os"

	sampleimage "github.com/nomuyoshi/sampleimage/lib"
	"github.com/spf13/pflag"
)

func main() {
	cli := &sampleimage.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		Flag:      pflag.NewFlagSet("sampleimage", pflag.ContinueOnError),
	}

	os.Exit(cli.Run(os.Args))
}
