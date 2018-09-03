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

type ledData struct {
	white int
	red   int
	green int
	blue  int
}

var (
	userLed, alertLed, emptyLed, fullLed ledData
)

func (data *ledData) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w, templates.IndexPage(data.white, data.red, data.green, data.blue))
	} else {
		r.ParseForm()
		var err error
		data.white, err = strconv.Atoi(r.PostForm.Get("white"))
		if err != nil {
			log.Error(err)
		}
		data.red, err = strconv.Atoi(r.PostForm.Get("red"))
		if err != nil {
			log.Error(err)
		}
		data.green, err = strconv.Atoi(r.PostForm.Get("green"))
		if err != nil {
			log.Error(err)
		}
		data.blue, err = strconv.Atoi(r.PostForm.Get("blue"))
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
		err := pwm.SetPwm(0+(i*4), 0, data.white)
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(1+(i*4), 0, data.red)
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(2+(i*4), 0, data.green)
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(3+(i*4), 0, data.blue)
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
	userLed = ledData{}
	emptyLed = ledData{}
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
	http.ListenAndServe(":5601", nil)
}
