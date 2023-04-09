class AudioChannel {
  public waveType = 0;
  public frequency = 440;
  public realFrequency = 440;
  public sampleRate = 44100;
  public volume = 0.25;
  public dutyCycle = 0;

  constructor(sampleRate: number) {
    this.sampleRate = sampleRate;
  }

  public doSin(t: number): number {
    return Math.sin(2 * Math.PI * this.realFrequency * (t / this.sampleRate));
  }

  public doSquare(t: number): number {
    let x = this.doSin(t);
    if (x > this.dutyCycle) {
      return 1;
    } else {
      return -1;
    }
  }

  public do(t: number): number {
    if (this.waveType == 1) {
      return this.doSquare(t) * this.volume;
    }
    return this.doSin(t) * this.volume;
  }
}

// @ts-ignore
globalThis.AudioChannel = AudioChannel;
