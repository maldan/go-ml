package maudio

import (
	ml_number "github.com/maldan/go-ml/util/number"
	"math"
)

const (
	WaveSin      uint32 = 1
	WaveSquare   uint32 = 2
	WaveTriangle uint32 = 4
	WaveNoise    uint32 = 8
)

type AudioChannel struct {
	WaveType     uint32
	Frequency    float32
	SampleRate   float32
	Volume       float32
	SquareOffset float32
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

/*func (c *AudioChannel) DoSaw(t float32) float32 {
	f1 := 3.1415926535
	f2 := (2 * 3.1415926535 * c.Frequency) * (t / c.SampleRate)
	f2s := math.Asin(math.Sin(float64(f2)))
	return float32(f1 * f2s)
}
*/
func (c *AudioChannel) DoTriangle(t float32) float32 {
	f1 := 2 / 3.1415926535
	f2 := (2 * 3.1415926535 * c.Frequency) * (t / c.SampleRate)
	f2s := math.Asin(math.Sin(float64(f2)))
	return float32(f1 * f2s)
}

func (c *AudioChannel) DoSquare(t float32) float32 {
	v := c.DoSin(t)
	if v >= c.SquareOffset {
		return 1
	} else {
		return -1
	}
}

func (c *AudioChannel) Do(t float32) float32 {
	if c.Volume <= 0 {
		return 0
	}

	if c.WaveType&WaveTriangle == WaveTriangle {
		return c.DoTriangle(t) * c.Volume
	}
	if c.WaveType&WaveNoise == WaveNoise {
		return c.DoNoise(t) * c.Volume
	}
	if c.WaveType&WaveSquare == WaveSquare {
		return c.DoSquare(t) * c.Volume
	}
	if c.WaveType&WaveSin == WaveSin {
		return c.DoSin(t) * c.Volume
	}

	return 0
}
