package maudio

import (
	"math"
	"reflect"
	"unsafe"
)

type AudioState struct {
	SampleRate int
	Position   int
	Length     int

	P             int
	Buffer        []float32
	DurationError float32
	Frequency     float32
	Volume        float32

	Time float32
}

var State AudioState = AudioState{}

func Init(sampleRate int) {
	State.SampleRate = sampleRate

	// Buffer size is one second
	State.Buffer = make([]float32, sampleRate*10)
}

func Tick(delta float32) {
	// State.DurationError =
	sampleAmountF := float32(State.SampleRate) * delta
	sampleAmountI := int(sampleAmountF)
	State.DurationError += sampleAmountF - float32(sampleAmountI)

	if State.DurationError >= 1 {
		State.DurationError -= 1
		sampleAmountI += 1
	}

	// Generate data
	for i := 0; i < sampleAmountI; i++ {
		val := math.Sin((float64(State.P) / float64(State.SampleRate)) * math.Pi * 2.0 * float64(State.Frequency))
		State.Buffer[State.Position] = float32(val) * State.Volume
		// State.Buffer[State.Position] = ml_number.RandFloat32(-0.1, 0.1)

		// Offset position
		State.P += 1
		State.Position += 1
		if State.Position >= len(State.Buffer)-1 {
			State.Position = 0
		}
	}

	State.Length += sampleAmountI

	State.Time += delta

	if State.Volume > 0 {
		State.Volume -= delta
	}
}

func GetState() map[string]any {
	bufferHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.Buffer))
	position := State.Position
	State.Position = 0

	return map[string]any{
		"bufferPointer":  bufferHeader.Data,
		"bufferPosition": position,
		//"bufferLength":  len(State.Buffer),

		"time":       State.Time,
		"sampleRate": State.SampleRate,
	}
}
