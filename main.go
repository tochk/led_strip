package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"

	_ "github.com/kidoman/embd/host/rpi"
	"github.com/tochk/led_strip/templates"
)

const ledMax = 4095

func Cie1931(b float64) float32 {
	b /= ledMax
	b *= 100
	if b <= 8 {
		return float32(b / 902.3) * ledMax
	}
	x := (float64(b) + 16) / 116
	x = x * x * x
	return float32(x) * ledMax
}


type ledData struct {
	white float64
	red   float64
	green float64
	blue  float64
}

var (
	userLed, alertLed, emptyLed, fullLed ledData
)

func alertHandler(w http.ResponseWriter, r *http.Request) {
	go alert()
}

func (data *ledData) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w, templates.IndexPage(data.white, data.red, data.green, data.blue))
	} else {
		r.ParseForm()
		var err error
		data.white, err = strconv.ParseFloat(r.PostForm.Get("white"), 64)
		if err != nil {
			log.Error(err)
		}
		data.red, err = strconv.ParseFloat(r.PostForm.Get("red"), 64)
		if err != nil {
			log.Error(err)
		}
		data.green, err = strconv.ParseFloat(r.PostForm.Get("green"), 64)
		if err != nil {
			log.Error(err)
		}
		data.blue, err = strconv.ParseFloat(r.PostForm.Get("blue"), 64)
		if err != nil {
			log.Error(err)
		}
		data.apply()
		http.Redirect(w, r, "/", 302)
	}
}

func (data *ledData) apply() {
	log.Infof("Apply: %+v", data)
	bus := embd.NewI2CBus(1)
	pwm := pca9685.New(bus, 0x40)
	for i := 0; i < 4; i++ {
		err := pwm.SetPwm(0+(i*4), 0, int(Cie1931(data.white)))
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(1+(i*4), 0, int(Cie1931(data.red)))
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(2+(i*4), 0, int(Cie1931(data.green)))
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(3+(i*4), 0, int(Cie1931(data.blue)))
		if err != nil {
			log.Error(err)
		}
	}
}

func alert() {
	for i := 0; i < 3; i++ {
		emptyLed.apply()
		time.Sleep(time.Millisecond * 300)
		alertLed.apply()
		time.Sleep(time.Millisecond * 300)
	}
	userLed.apply()
}

func main() {
	userLed = ledData{
		white: 4095,
		red:   4095,
		green: 4095,
		blue:  4095,
	}
	emptyLed = ledData{
		white: 0,
		red:   0,
		green: 0,
		blue:  0,
	}
	fullLed = ledData{
		white: 4095,
		red:   4095,
		green: 4095,
		blue:  4095,
	}
	alertLed = ledData{
		white: 700,
		red:   4095,
		green: 0,
		blue:  0,
	}
	fullLed.apply()
	http.HandleFunc("/", userLed.indexHandler)
	http.HandleFunc("/alert/", alertHandler)
	http.ListenAndServe(":5601", nil)
}
