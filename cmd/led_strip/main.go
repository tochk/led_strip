package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tochk/led_strip/internal/app/controller"
	"github.com/tochk/led_strip/internal/app/templates"
)

var currentLed controller.Led

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w, templates.IndexPage(currentLed.GetBrightness(controller.White),
			currentLed.GetBrightness(controller.Red),
			currentLed.GetBrightness(controller.Green),
			currentLed.GetBrightness(controller.Blue)))
	} else {
		r.ParseForm()
		var err error
		white, err := strconv.ParseFloat(r.PostForm.Get("white"), 64)
		if err != nil {
			log.Error(err)
		}
		currentLed.SetBrightness(controller.White, white)

		red, err := strconv.ParseFloat(r.PostForm.Get("red"), 64)
		if err != nil {
			log.Error(err)
		}
		currentLed.SetBrightness(controller.Red, red)

		green, err := strconv.ParseFloat(r.PostForm.Get("green"), 64)
		if err != nil {
			log.Error(err)
		}
		currentLed.SetBrightness(controller.Green, green)

		blue, err := strconv.ParseFloat(r.PostForm.Get("blue"), 64)
		if err != nil {
			log.Error(err)
		}
		currentLed.SetBrightness(controller.Blue, blue)

		currentLed.Apply()
		http.Redirect(w, r, "/", 302)
	}
}

func alert() {
	if time.Now().Hour() < 8 || time.Now().Hour() > 20 {
		for i := 0; i < 5; i++ {
			controller.EmptyLed.Apply()
			time.Sleep(time.Millisecond * 300)
			controller.AlertLed.Apply()
			time.Sleep(time.Millisecond * 300)
		}
		currentLed.Apply()
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	go alert()
}

func main() {
	controller.FullLed.Apply()
	currentLed = controller.New(100, 100, 100, 100)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/alert/", alertHandler)

	http.ListenAndServe(":5601", nil)
}
