package main

import (
	maudio "github.com/maldan/go-ml/audio"
)

func main() {
	/*content, _ := ml_json.FromFile[map[string]any]("test.json")
	fmt.Printf("%v\n", content)
	content["x"] = "gas"
	ml_json.ToFile("test.json", content)
	*/
	maudio.Init(44100 / 4)

	for i := 0; i < 90; i++ {
		maudio.Tick(0.0112)
	}

	// ml_file.New("sas.raw").Write(maudio.State.Buffer)
}
