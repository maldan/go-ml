class MyAudioProcessor extends AudioWorkletProcessor {
  constructor() {
    super();

    this.port.onmessage = this.onmessage.bind(this);
  }

  async onmessage(e) {
    console.log(e.data);

    if (e.data.type === "init") {
      console.log(new TextEncoder());
      globalThis.crypto = 123;
      globalThis.performance = 123;

      eval(e.data.ggt);

      const go = new Go();
      const wasmModule = await WebAssembly.instantiate(
        e.data.ab,
        go.importObject
      );
      console.log(wasmModule);

      /*const go = new Go();
      const b = await fetch("./main.wasm");
      const bytes = await b.blob();
      const wasmModule = await WebAssembly.instantiate(
        await bytes.arrayBuffer(),
        go.importObject
      );
      go.run(wasmModule.instance);*/
    }
  }

  process(inputList, outputList, parameters) {
    const output = outputList[0];
    const channel = output[0];

    for (let i = 0; i < channel.length; ++i) {
      channel[i] = Math.random();
    }

    return true;
  }
}

registerProcessor("my-audio-processor", MyAudioProcessor);
