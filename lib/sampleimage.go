package sampleimage

import (
	"fmt"
	"github.com/spf13/pflag"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	exitCodeOK int = iota
	exitCodeErr
)

var (
	width  int
	height int
	bg     string
)

// CLI はio.WriterとFlagSetの管理を行う
type CLI struct {
	OutStream, ErrStream io.Writer
	Flag                 *pflag.FlagSet
}

func (c *CLI) flagSettings(args []string) {
	c.Flag.IntVarP(&width, "width", "W", 100, "image width")
	c.Flag.IntVarP(&height, "height", "H", 100, "image height")
	c.Flag.StringVarP(&bg, "bg", "c", "gray", "image background color.")
	c.Flag.SetOutput(c.OutStream)
	c.Flag.Usage = func() {
		fmt.Fprintf(c.OutStream, "Usage: sampleimage [file] [options]\n\n")
		fmt.Fprintf(c.OutStream, "Available extentions:\n  jpeg, jpg, png\n")
		fmt.Fprintf(c.OutStream, "Available background colors:\n  %s\n", strings.Join(colorNames(), ", "))
		fmt.Fprintf(c.OutStream, "Options:\n")
		c.Flag.PrintDefaults()
	}
}

// Run はsampleimageコマンドのメイン処理
func (c *CLI) Run(args []string) int {
	c.flagSettings(args)
	if err := c.Flag.Parse(args[1:]); err != nil {
		if err == pflag.ErrHelp {
			return exitCodeOK
		}

		fmt.Fprintln(c.ErrStream, err)
		return exitCodeErr
	}
	path := args[1]
	if isExists(path) {
		fmt.Fprintf(c.ErrStream, "%s already exits.\n", path)
		return exitCodeErr
	}
	if dir := filepath.Dir(path); !isExists(dir) {
		fmt.Fprintf(c.ErrStream, "%s does not exit.\n", dir)
		return exitCodeErr
	}
	ext := filepath.Ext(path)
	if isInValidExt(ext) {
		fmt.Fprintf(c.ErrStream, "%s is invalid extention.\n", ext)
		return exitCodeErr
	}
	if isInValidBg(bg) {
		fmt.Fprintf(c.ErrStream, "%s is invalid color.\n", bg)
		return exitCodeErr
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	setBgColor(img, bg)

	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(c.ErrStream, "Error: create image file failed.")
		return exitCodeErr
	}
	defer file.Close()
	if err := encode(file, img, ext); err != nil {
		fmt.Fprintln(c.ErrStream, "Error: encode failed.")
		return exitCodeErr
	}
	return exitCodeOK
}

func encode(f *os.File, img *image.RGBA, ext string) error {
	switch ext {
	case ".jpeg", ".jpg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	case ".png":
		return png.Encode(f, img)
	default:
		return fmt.Errorf("Error: encode failed")
	}
}

func setBgColor(img *image.RGBA, bg string) {
	bgColor := color.RGBA{bgColors[bg]["r"], bgColors[bg]["g"], bgColors[bg]["b"], bgColors[bg]["a"]}
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, bgColor)
		}
	}
}

func isInValidExt(ext string) bool {
	return ext != ".jpg" && ext != ".jpeg" && ext != ".png"
}

func isInValidBg(bg string) bool {
	_, ok := bgColors[bg]
	return !ok
}

func isExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func colorNames() []string {
	names := make([]string, 0, len(bgColors))
	for key := range bgColors {
		names = append(names, key)
	}
	return names
}
