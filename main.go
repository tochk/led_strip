package main

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	_ "github.com/kidoman/embd/host/rpi"
	log "github.com/sirupsen/logrus"
)

func main() {
	bus := embd.NewI2CBus(1)
	pwm := pca9685.New(bus, 0x40)
	for i := 0; i < 3; i++ {
		err := pwm.SetPwm(0 + (i * 4), 0, 4000)
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(1 + (i * 4), 0, 4000)
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(2 + (i * 4), 0, 4000)
		if err != nil {
			log.Error(err)
		}
		err = pwm.SetPwm(3 + (i * 4), 0, 4000)
		if err != nil {
			log.Error(err)
		}
	}
}
