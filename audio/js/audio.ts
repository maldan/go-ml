const NoteFrequency: number[] = [
  ...[
    16.35, 17.32, 18.35, 19.45, 20.6, 21.83, 23.12, 24.5, 25.96, 27.5, 29.14,
    30.87,
  ], // 0
  ...[
    32.7, 34.65, 36.71, 38.89, 41.2, 43.65, 46.25, 49, 51.91, 55, 58.27, 61.74,
  ], // 1
  ...[
    65.41, 69.3, 73.42, 77.78, 82.41, 87.31, 92.5, 98, 103.83, 110, 116.54,
    123.47,
  ], // 2
  ...[
    130.81, 138.59, 146.83, 155.56, 164.81, 174.61, 185, 196, 207.65, 220,
    233.08, 246.94,
  ], // 3
  ...[
    261.63, 277.18, 293.66, 311.13, 329.63, 349.23, 369.99, 392, 415.3, 440,
    466.16, 493.88,
  ], // 4
  ...[
    523.25, 554.37, 587.33, 622.25, 659.25, 698.46, 739.99, 783.99, 830.61, 880,
    932.33, 987.77,
  ], // 5
  ...[
    1046.5, 1108.73, 1174.66, 1244.51, 1318.51, 1396.91, 1479.98, 1567.98,
    1661.22, 1760, 1864.66, 1975.53,
  ], // 6
  ...[
    2093, 2217.46, 2349.32, 2489.02, 2637.02, 2793.83, 2959.96, 3135.96,
    3322.44, 3520, 3729.31,
  ], // 7
  ...[
    3951.07, 4186.01, 4434.92, 4698.63, 4978.03, 5274.04, 5587.65, 5919.91,
    6271.93, 6644.88, 7040, 7458.62, 7902.13,
  ], // 8
];

class MegaAudio {
  static _audioContext: AudioContext;
  static _player: AudioWorkletNode;
  static _analyzer: AnalyserNode;

  static async init() {
    // @ts-ignore
    this._audioContext = new (window.AudioContext || window.webkitAudioContext)(
      {
        latencyHint: "interactive",
        sampleRate: 44100,
      }
    );

    await this._audioContext.audioWorklet.addModule(
      "./js/audio/audio_worklet.js"
    );

    // Load synth
    this._player = new AudioWorkletNode(
      this._audioContext,
      "my-audio-processor"
    );
    this._player.port.postMessage({
      type: "init",
      sampleRate: this._audioContext.sampleRate,
      synthText: await (await fetch("./js/audio/synth.js")).text(),
      samplePlayerText: await (
        await fetch("./js/audio/sample_player.js")
      ).text(),
    });

    this._analyzer = this._audioContext.createAnalyser();
    this._player.connect(this._analyzer);
    this._analyzer.connect(this._audioContext.destination);
  }

  static capture(): Uint8Array {
    const dataArray = new Uint8Array(this._analyzer.frequencyBinCount);
    this._analyzer.getByteTimeDomainData(dataArray);
    return dataArray;
  }

  static sendData(data: any) {
    this._player.port.postMessage({
      ...data,
      type: "setChannelData",
    });
  }

  static async loadSample(name: string, url: string) {
    let buffer = await (await fetch(url)).arrayBuffer();
    const audioBuffer = await this._audioContext.decodeAudioData(buffer);

    this._player.port.postMessage({
      name,
      data: audioBuffer.getChannelData(0),
      type: "loadSample",
    });
  }

  static playSample(
    name: string,
    channel: string,
    volume: number,
    pitch: number
  ) {
    this._player.port.postMessage({
      name,
      channel,
      pitch,
      volume,
      type: "playSample",
    });
  }
}

// @ts-ignore
window.MegaAudio = MegaAudio;
