class GoAudio {
  static _audioContext;

  static async init() {
    this._audioContext = new (window.AudioContext || window.webkitAudioContext)(
      {
        latencyHint: "interactive",
        sampleRate: 44100,
      }
    );

    await this._audioContext.resume();
    await this._audioContext.audioWorklet.addModule("./audio_worklet.js");

    // Goshan
    // const go = new Go();
    const b = await fetch("./main.wasm");
    const bytes = await b.blob();
    const ab = await bytes.arrayBuffer();

    const gg = await fetch("./go.js");
    const ggt = await gg.text();

    /*const wasmModule = await WebAssembly.instantiate(
      await bytes.arrayBuffer(),
      go.importObject
    );*/

    const player = new AudioWorkletNode(
      this._audioContext,
      "my-audio-processor"
    );
    player.port.postMessage({
      type: "init",
      ab,
      ggt,
    });
    player.connect(this._audioContext.destination);
  }
}

window.GoAudio = GoAudio;
