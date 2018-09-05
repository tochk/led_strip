package controller

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"

	log "github.com/Sirupsen/logrus"
	_ "github.com/kidoman/embd/host/rpi"
)

var server struct {
	pwm *pca9685.PCA9685
	bus embd.I2CBus
}

const (
	ledMax = 4095
	White  = iota
	Red
	Green
	Blue
)

type Led struct {
	white float64
	red   float64
	green float64
	blue  float64
}

func init() {
	server.bus = embd.NewI2CBus(1)
	server.pwm = pca9685.New(server.bus, 0x40)
}

func cie1931(b float64) float32 {
	if b <= 8 {
		return float32(b/902.3) * ledMax
	}
	x := (float64(b) + 16) / 116
	x = x * x * x
	return float32(x) * ledMax
}

var (
	AlertLed = Led{
		white: 50,
		red:   100,
	}
	EmptyLed = Led{}
	FullLed  = Led{
		white: 100,
		red:   100,
		green: 100,
		blue:  100,
	}
)

func (led *Led) Apply() {
	for i := 0; i < 4; i++ {
		err := server.pwm.SetPwm(0+(i*4), 0, int(cie1931(led.white)))
		if err != nil {
			log.Error(err)
		}
		err = server.pwm.SetPwm(1+(i*4), 0, int(cie1931(led.red)))
		if err != nil {
			log.Error(err)
		}
		err = server.pwm.SetPwm(2+(i*4), 0, int(cie1931(led.green)))
		if err != nil {
			log.Error(err)
		}
		err = server.pwm.SetPwm(3+(i*4), 0, int(cie1931(led.blue)))
		if err != nil {
			log.Error(err)
		}
	}
}

func New(white, red, green, blue float64) (led Led) {
	led.SetBrightness(White, white)
	led.SetBrightness(Red, red)
	led.SetBrightness(Green, green)
	led.SetBrightness(Blue, blue)
	return
}

func (led *Led) SetBrightness(color int, value float64) {
	if value > 100 || value < 0 {
		value = 0
	}
	switch color {
	case White:
		led.white = value
	case Red:
		led.red = value
	case Green:
		led.green = value
	case Blue:
		led.blue = value
	}
}

func (led *Led) GetBrightness(color int) float64 {
	switch color {
	case White:
		return led.white
	case Red:
		return led.red
	case Green:
		return led.green
	case Blue:
		return led.blue
	default:
		return 0
	}
}
