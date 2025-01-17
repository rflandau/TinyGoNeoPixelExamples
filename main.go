/*
Examples of NeoPixel LED manipulation via TinyGo.

# Performed on an Adafruit CircuitPlayground Express, but should work on any device with NeoPixels.

The initialization code is based on a gist by Chris Amico: https://gist.github.com/eyeseast/9566cbeabd587654c7ea647706731217

I just created a couple examples of how to make the LEDs do things.
*/
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const LEDCOUNT int = 10 // The express has ten independent LEDs in a ring; change this to match your device's count

var (
	RED   color.RGBA = color.RGBA{R: 0xFF}
	GREEN color.RGBA = color.RGBA{G: 0xFF}
	BLUE  color.RGBA = color.RGBA{B: 0xFF}
	OFF   color.RGBA = color.RGBA{}
)

var (
	neo   machine.Pin   // pin for interacting with the neopixels
	delay time.Duration = 500 * time.Millisecond
)

func init() {
	// per the machine reference https://tinygo.org/docs/reference/microcontrollers/machine/circuitplay-express/
	neo = machine.NEOPIXELS
}

func main() {
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// the LEDs are manipulate via ws2812
	ws := ws2812.NewWS2812(neo)

	// UNCOMMENT ONE OF THE EFFECT SUBROUTINES
	SlowFillClear(ws)
	//Spinning(ws, GREEN)
}

//#region effect subroutines

// Sets the LEDs to static red, green, and blue
func RedGreenBlue(ws ws2812.Device) {
	var leds [LEDCOUNT]color.RGBA
	for i := range leds {
		if i%3 == 0 {
			leds[i] = color.RGBA{R: 0xff, G: 0x00, B: 0x00}
		} else if i%3 == 1 {
			leds[i] = color.RGBA{R: 0x00, G: 0xff, B: 0x00}
		} else {
			leds[i] = color.RGBA{R: 0x00, G: 0x00, B: 0xff}
		}
	}

	ws.WriteColors(leds[:])
}

// Slowly lights each LED until the final one is lit, at which point all will be wiped and the fill starts again
func SlowFillClear(ws ws2812.Device) {
	var (
		leds   [LEDCOUNT]color.RGBA
		curLED int = 0
	)
	for {
		if curLED == LEDCOUNT {
			AllOff(leds[:])
			curLED = 0
		} else {
			// light the current LED
			leds[curLED] = color.RGBA{R: 0x00, G: 0x00, B: 0xff}
			curLED += 1
		}
		ws.WriteColors(leds[:])
		time.Sleep(delay)
	}
}

// Spins a single LED around the ring
func Spinning(ws ws2812.Device, spinColor color.RGBA) {
	var (
		// assumes len(leds) == LEDCOUNT
		leds   [LEDCOUNT]color.RGBA
		curLED int = 0
	)

	for {
		// turn off the previous LED
		leds[decrWrap(curLED, LEDCOUNT)] = OFF
		// set the current LED
		leds[curLED] = spinColor
		// write the changes
		ws.WriteColors(leds[:])
		// increment for next loop
		curLED = incrWrap(curLED, LEDCOUNT)

		time.Sleep(delay)
	}
}

//#endregion effect subroutines

// Increments the given value around the limit (0 <= value < limit)
func incrWrap(value int, limit int) int {
	if value >= limit-1 {
		return 0
	}
	return value + 1
}

// Decrements the given value around the limit (0 <= value < limit)
func decrWrap(value int, limit int) int {
	if value <= 0 {
		return limit - 1
	}
	return value - 1
}

var emptyLEDs [LEDCOUNT]color.RGBA

// Turn off all LEDs.
// Overwrites the LED array with an array of zero'd color.RGBA{}s
func AllOff(leds []color.RGBA) {
	for i := range leds {
		leds[i] = OFF
	}
}
