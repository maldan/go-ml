package maudio

import (
	ml_number "github.com/maldan/go-ml/util/number"
	"math"
)

type AudioChannel struct {
	WaveType   string
	Frequency  float32
	SampleRate float32
	Volume     float32
}

func (c *AudioChannel) DoNoise(t float32) float32 {
	return ml_number.RandFloat32(-1, 1)
}

func (c *AudioChannel) DoSin(t float32) float32 {
	v := 2 * 3.1415926535 * c.Frequency
	f := t / c.SampleRate
	final := float32(math.Sin(float64(v * f)))
	return final
}

func (c *AudioChannel) DoTriangle(t float32) float32 {
	p := float64(3.1415926535/2) * float64(c.Frequency/c.SampleRate)

	tt := float64(math.Mod(float64(t), 2*p))

	y := (1.0 / p) * (p - math.Abs(tt-p))

	// final := math.Abs(hh-1)*1 - 1
	return float32(y)
}

func (c *AudioChannel) DoSquare(t float32) float32 {
	v := c.DoSin(t)
	if v > 0 {
		return 1
	} else {
		return -1
	}
}

func (c *AudioChannel) Do(t float32) float32 {
	if c.Volume <= 0 {
		return 0
	}
	if c.WaveType == "triangle" {
		return c.DoTriangle(t) * c.Volume
	}
	if c.WaveType == "noise" {
		return c.DoNoise(t) * c.Volume
	}
	if c.WaveType == "square" {
		return c.DoSquare(t) * c.Volume
	}
	return c.DoSin(t) * c.Volume
}
