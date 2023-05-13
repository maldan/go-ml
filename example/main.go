package main

import (
	"fmt"
	"time"
)

func main() {
	/*content, _ := ml_json.FromFile[map[string]any]("test.json")
	fmt.Printf("%v\n", content)
	content["x"] = "gas"
	ml_json.ToFile("test.json", content)
	*/
	/*maudio.Init(44100 / 4)

	for i := 0; i < 90; i++ {
		maudio.Tick(0.0112)
	}*/

	t := time.Now()
	a := 0
	b := 1
	for i := 0; i < 1_000_000; i++ {
		a += b
	}
	fmt.Printf("%v\n", time.Since(t))
	fmt.Printf("%v\n", a)

	// ml_file.New("sas.raw").Write(maudio.State.Buffer)
}
