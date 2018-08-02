package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
)

func main() {
	pwm := pca9685.New(embd.NewI2CBus(0), 0)
	err := pwm.SetPwm(1, 2048, 2048)
	if err != nil {
		log.Error(err)
	}
}
