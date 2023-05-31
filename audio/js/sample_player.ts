class AudioSampleChannel {
  public volume = 0.0;
  public sampleRate;

  constructor(sampleRate: number) {
    this.sampleRate = sampleRate;
  }

  public do(): number {
    return 0;
  }
}

// @ts-ignore
globalThis.AudioSampleChannel = AudioSampleChannel;
