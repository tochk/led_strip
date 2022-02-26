# led_strip

Led strip controller based on Raspberry Pi & pca9685 with web interface.
With this interface you can view and adjust brightness (red, green, blue, white) .

## Installation

Just write in your shell:
```shell
make build
```

## Connections

pca9685 connection to Raspberry Pi:
- vcc - GPIO 1
- sda - GPIO 3
- scl - GPIO 5
- gnd - GPIO 9
- oe - GPIO 11

You can connect power via USB, or using GPIO:
- rpi power (5v) - 2
- rpi power gnd - 6
