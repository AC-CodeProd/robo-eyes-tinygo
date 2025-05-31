package main

import (
	"image/color"
	"machine"
	roboeyestinygo "robo-eyes-tinygo"
	"time"

	"tinygo.org/x/drivers/sh1106"
)

// SH1106Display implements DisplayInterface for SH1106
type SH1106Display struct {
	device *sh1106.Device
}

func (d *SH1106Display) ClearBuffer() {
	d.device.ClearBuffer()
}

func (d *SH1106Display) ClearDisplay() {
	d.device.ClearDisplay()
}

func (d *SH1106Display) Display() error {
	return d.device.Display()
}

func (d *SH1106Display) SetPixel(x, y int16, c color.RGBA) {
	d.device.SetPixel(int16(x), int16(y), c)
}

func (d *SH1106Display) Size() (int16, int16) {
	return d.device.Size()
}

const (
	width  = 128
	height = 64
)

func main() {
	// led := machine.LED
	// led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	time.Sleep(time.Millisecond * 100)
	machine.I2C1.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})
	time.Sleep(time.Second * 1)
	device := sh1106.NewI2C(machine.I2C1)

	// Configure display (specific to SH1106 driver)
	device.Configure(sh1106.Config{
		Width:    width,
		Height:   height,
		Address:  sh1106.Address, // Adresse I2C
		VccState: sh1106.SWITCHCAPVCC,
	})

	// Create adapter
	adapter := &SH1106Display{device: &device}

	// // Initialize eyes
	eyes := roboeyestinygo.RoboEyes{}
	eyes.Begin(adapter, width, height, 50) // 128x64 OLED @ 50 FPS
	// eyes.Debug()
	// Set expressions
	// eyes.SetDirection(roboeyestinygo.DirCenter)
	// eyes.SetMood(roboeyestinygo.Mood)
	// eyes.CloseEyes(true, false)
	eyes.SetAutoBlinkerWithInterval(true, 3, 2) // Blink every 3-5 seconds
	eyes.SetIdleModeWithInterval(true, 2, 2)    // Start idle animation cycle (eyes looking in random directions) -> turn on/off, set interval between each eye repositioning in full seconds, set range for random time interval variation in full seconds
	// eyes.SetCuriosity(true)
	// eyes.SetIdleMode(true)
	// eyes.SetCyclops(true)

	for {
		eyes.Update()
		time.Sleep(20 * time.Millisecond)
	}
}
