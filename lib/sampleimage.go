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

// CLI はio.WriterとFlagSetの管理を行う
type CLI struct {
	OutStream, ErrStream io.Writer
}

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

func newFlag(outStream io.Writer) *pflag.FlagSet {
	flag := pflag.NewFlagSet("sampleimage", pflag.ContinueOnError)
	flag.IntVarP(&width, "width", "W", 300, "image width")
	flag.IntVarP(&height, "height", "H", 300, "image height")
	flag.StringVarP(&text, "text", "t", "SAMPLE", "image text")
	flag.StringVar(&bg, "bg", "gray", "image background color")
	flag.SetOutput(outStream)
	flag.Usage = func() {
		fmt.Fprintf(outStream, "Usage: sampleimage [file] [options]\n\n")
		fmt.Fprintf(outStream, "Available extentions:\n  jpeg, jpg, png\n")
		fmt.Fprintf(outStream, "Available background colors:\n  %s\n", strings.Join(colorNames(), ", "))
		fmt.Fprintf(outStream, "Options:\n")
		flag.PrintDefaults()
	}

	return flag
}

// Run はsampleimageコマンドのメイン処理
func (c *CLI) Run(args []string) int {
	flag := newFlag(c.OutStream)
	if err := flag.Parse(args[1:]); err != nil {
		if err == pflag.ErrHelp {
			return exitCodeOK
		}

		fmt.Fprintln(c.ErrStream, err)
		return exitCodeErr
	}

	path := flag.Arg(0)
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
	if len(path) == 0 {
		return fmt.Errorf("Please specify an output path")
	}
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
