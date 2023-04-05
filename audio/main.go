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

	Effect AudioEffect

	Ch0 AudioChannel
	Ch1 AudioChannel
}

var State AudioState = AudioState{}

func Init(sampleRate float32) {
	State.SampleRate = sampleRate

	// Buffer size is one second
	State.Buffer = make([]float32, int(sampleRate*2))

	State.Ch0.SampleRate = sampleRate
	State.Ch0.Frequency = 440
	State.Ch0.Volume = 1.0
}

func SetEffect(fromF float32, toF float32, time float32) {
	State.Effect.FromF = fromF
	State.Effect.ToF = toF
	State.Effect.Time = time
}

func Tick(samples int) {
	for i := 0; i < samples; i++ {
		ch0 := State.Ch0.Do(float32(i) + State.T)
		ch1 := State.Ch1.Do(float32(i) + State.T)
		State.Buffer[i] = (ch0 + ch1)
	}
	State.T += float32(samples)

	/*// State.DurationError =
	sampleAmountF := float32(State.SampleRate) * delta
	sampleAmountI := int(sampleAmountF)
	State.DurationError += sampleAmountF - float32(sampleAmountI)

	if State.DurationError >= 1 {
		State.DurationError -= 1
		sampleAmountI += 1
	}

	if State.Effect.Time > 0 {
		State.Effect.Time -= delta
		State.Frequency = ml_number.Lerp(State.Effect.FromF, State.Effect.ToF, State.Effect.Time)
	}

	// Generate data
	for i := 0; i < sampleAmountI; i++ {
		val := math.Sin((float64(State.P) / float64(State.SampleRate)) * math.Pi * 2.0 * float64(State.Frequency))
		if val > 0 {
			val = 1
		} else {
			val = -1
		}
		State.Buffer[State.Position] = float32(val) * State.Volume

		// State.Buffer[State.Position] = ml_number.RandFloat32(-0.1, 0.1)

		// Offset position
		State.P += 1
		State.Position += 1
		if State.Position >= len(State.Buffer)-1 {
			State.Position = 0
		}
	}

	if State.Volume > 0 {
		State.Volume -= delta
	} else {
		State.Volume = 0
	}*/
}

func GetState() map[string]any {
	bufferHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.Buffer))

	return map[string]any{
		"bufferPointer": bufferHeader.Data,
		"sampleRate":    State.SampleRate,
	}
}
