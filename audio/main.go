package maudio

import (
	"reflect"
	"unsafe"
)

type AudioEffect struct {
	FromF float32
	ToF   float32
	Time  float32
}

type AudioState struct {
	SampleRate float32

	T      float32
	Buffer []float32
	Volume float32

	CurrentLength int

	Effect AudioEffect
	Ch0    AudioChannel
	Ch1    AudioChannel
}

var State AudioState = AudioState{}
var externalState = map[string]any{
	"bufferPointer": 0,
	"bufferLength":  0,
	"sampleRate":    0,
}

func Init(sampleRate float32) {
	State.SampleRate = sampleRate

	// Buffer size is one second
	State.Buffer = make([]float32, int(sampleRate*2))

	State.Ch0.SampleRate = sampleRate
	State.Ch0.Frequency = 440
	State.Ch0.Volume = 1.0
}

func Tick(samples int) {
	for i := 0; i < samples; i++ {
		// State.Ch0.Frequency = 440 * float32(math.Sin(float64(i)/512.0)*0.5)
		// State.Ch0.SquareOffset = float32(math.Sin(float64(i)/512.0) * 0.5)
		ch0 := State.Ch0.Do(float32(i) + State.T)
		ch1 := State.Ch1.Do(float32(i) + State.T)
		State.Buffer[i] = ch0 + ch1
	}
	State.T += float32(samples)
	State.CurrentLength += samples
}

func GetState() map[string]any {
	bufferHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.Buffer))
	length := State.CurrentLength
	State.CurrentLength = 0

	externalState["bufferPointer"] = bufferHeader.Data
	externalState["bufferLength"] = length
	externalState["sampleRate"] = State.SampleRate

	return externalState
}
