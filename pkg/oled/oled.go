package oled

import (
	"fmt"
	"image"
	"os"

	"github.com/golang/freetype/truetype"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
)

type Oled struct {
	oled   *i2c.SSD1306Driver
	height int
	width  int
}

func CreateOled() *Oled {
	board := raspi.NewAdaptor()
	oled := Oled{
		oled:   i2c.NewSSD1306Driver(board),
		height: 64,
		width:  128,
	}
	oled.oled.Start()
	return &oled
}

func (s *Oled) Display(text string) error {
	s.oled.Clear()
	img := s.createImage(text)
	return s.oled.ShowImage(img)
}

func (s *Oled) createImage(text string) image.Image {
	img := image.NewGray(image.Rect(0, 0, s.width, s.height))
	textSize := 20
	opt := truetype.Options{
		Size: float64(textSize),
	}
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(s.width) - dr.MeasureString(text)) / 2
	dr.Dot.Y = (fixed.I(s.height)) / 2
	dr.DrawString(text)
	return img
}
