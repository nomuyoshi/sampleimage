package sampleimage

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	JpegPath = "./testdata/test.jpeg"
	JpgPath  = "./testdata/test.jpg"
	PngPath  = "./testdata/test.png"
)

func TestMain(m *testing.M) {
	code := m.Run()

	removeTestImages()
	os.Exit(code)
}

func removeTestImages() {
	testImages := [3]string{JpegPath, JpgPath, PngPath}
	for _, path := range testImages {
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}
}

func TestRun(t *testing.T) {
	t.Run("Failed: empty path", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage --bg=red", " ")
		gotCode := cli.Run(args)
		gotOutput := buffer.String()
		wantOutput := "Please specify an output path\n"

		if gotCode != exitCodeErr {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeErr, gotCode)
		}
		if gotOutput != wantOutput {
			t.Errorf("unexpected output. want: %s, got: %s", wantOutput, gotOutput)
		}
	})
	t.Run("Failed: invalid bg", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage ./testdata/test.jpg --bg=sliver", " ")
		gotCode := cli.Run(args)
		gotOutput := buffer.String()
		wantOutput := "sliver is invalid color\n"

		if gotCode != exitCodeErr {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeErr, gotCode)
		}
		if gotOutput != wantOutput {
			t.Errorf("unexpected output. want: %s, got: %s", wantOutput, gotOutput)
		}
	})
	t.Run("Failed: invalid extentions", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage ./testdata/test.psd --width=200 --height=200 --bg=black", " ")
		gotCode := cli.Run(args)
		gotOutput := buffer.String()
		wantOutput := ".psd is invalid extention\n"

		if gotCode != exitCodeErr {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeErr, gotCode)
		}
		if gotOutput != wantOutput {
			t.Errorf("unexpected output. want: %s, got: %s", wantOutput, gotOutput)
		}
	})
	t.Run("Success: jpg", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage "+JpgPath+" --width=200 --height=200 --bg=black", " ")
		gotCode := cli.Run(args)
		if gotCode != exitCodeOK {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeOK, gotCode)
		}
	})
	t.Run("Success: jpeg", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage "+JpegPath+" --width=200 --height=200 --bg=black", " ")
		gotCode := cli.Run(args)
		if gotCode != exitCodeOK {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeOK, gotCode)
		}
	})
	t.Run("Success: png", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage "+PngPath+" --width=200 --height=200 --bg=black", " ")
		gotCode := cli.Run(args)
		if gotCode != exitCodeOK {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeOK, gotCode)
		}
	})
	t.Run("Success: Show usage", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cli := &CLI{OutStream: buffer, ErrStream: buffer}
		args := strings.Split("sampleimage --help", " ")
		wantOutput := "Usage: sampleimage [file] [options]"
		gotCode := cli.Run(args)
		gotOutput := buffer.String()

		if gotCode != exitCodeOK {
			t.Errorf("unexpected exit code. want: %d, got: %d.", exitCodeOK, gotCode)
		}

		if !strings.HasPrefix(gotOutput, wantOutput) {
			t.Errorf("unexpected output. want: %s, got: %s", wantOutput, gotOutput)
		}
	})
}
