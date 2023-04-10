class AudioChannel {
  public waveType = 0;
  public frequency = 440;
  public newFrequency = 440;
  public sampleRate = 44100;
  public volume = 0.0;
  public newVolume = 0.0;
  public dutyCycle = 0;

  public lastSin = 0;
  public phaseOffset = 0;

  constructor(sampleRate: number) {
    this.sampleRate = sampleRate;
  }

  public doSin(t: number): number {
    let v = Math.sin(2 * Math.PI * this.frequency * (t / this.sampleRate));
    let isUp = v > this.lastSin;

    if (this.frequency != this.newFrequency) {
      if (Math.abs(v) < 0.1) {
        this.frequency = this.newFrequency;

        // Find next phase
        for (let i = 0; i < 4096; i++) {
          let v2 = Math.sin(
            2 * Math.PI * this.newFrequency * ((t + this.phaseOffset) / 44100)
          );
          let v3 = Math.sin(
            2 *
              Math.PI *
              this.newFrequency *
              ((t + this.phaseOffset + 1) / 44100)
          );

          if (v3 > v2 && isUp && Math.abs(v2) < 0.1) {
            break;
          }
          if (v3 < v2 && !isUp && Math.abs(v2) < 0.1) {
            break;
          }

          this.phaseOffset += 1;
        }
      }
    }

    this.lastSin = v;
    return v;
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
