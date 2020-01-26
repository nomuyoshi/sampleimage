package sampleimage

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/spf13/pflag"
	"golang.org/x/image/font/gofont/gobold"
)

const (
	exitCodeOK int = iota
	exitCodeErr
)

var (
	width  int
	height int
	bg     string
	text   string
)

// CLI はio.WriterとFlagSetの管理を行う
type CLI struct {
	OutStream, ErrStream io.Writer
	Flag                 *pflag.FlagSet
}

func (c *CLI) flagSettings(args []string) {
	c.Flag.IntVarP(&width, "width", "W", 300, "image width")
	c.Flag.IntVarP(&height, "height", "H", 300, "image height")
	c.Flag.StringVarP(&text, "text", "t", "SAMPLE", "image text")
	c.Flag.StringVar(&bg, "bg", "gray", "image background color")
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
	ext := filepath.Ext(path)
	if err := validateArgs(path, ext); err != nil {
		fmt.Fprintln(c.ErrStream, err)
		return exitCodeErr
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	size := 20.0
	setBg(img)
	addText(img, 10, int(float64(height)/2-size/2), size)

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

func setBg(img *image.RGBA) {
	bgColor := color.RGBA{bgColors[bg]["r"], bgColors[bg]["g"], bgColors[bg]["b"], bgColors[bg]["a"]}
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, bgColor)
		}
	}
}

func validateArgs(path string, ext string) error {
	// 出力先に同名のファイルが存在するか
	if isExists(path) {
		return fmt.Errorf("%s already exits", path)
	}
	// 出力先のディレクトリが存在するか
	if dir := filepath.Dir(path); !isExists(dir) {
		return fmt.Errorf("%s does not exit", dir)
	}

	// 対応している拡張子か
	if isInValidExt(ext) {
		return fmt.Errorf("%s is invalid extention", ext)
	}

	// 対応している背景色か
	if isInValidBg() {
		return fmt.Errorf("%s is invalid color", bg)
	}

	return nil
}

func isInValidExt(ext string) bool {
	return ext != ".jpg" && ext != ".jpeg" && ext != ".png"
}

func isInValidBg() bool {
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

func addText(img *image.RGBA, x, y int, size float64) error {
	ft, _ := freetype.ParseFont(gobold.TTF)
	c := freetype.NewContext()
	c.SetDst(img)
	c.SetDPI(72.0)
	c.SetFont(ft)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetSrc(image.White)

	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))
	_, err := c.DrawString(text, pt)
	return err
}
