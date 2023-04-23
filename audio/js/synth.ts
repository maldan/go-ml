class AudioChannel {
  public waveType = 0;
  public frequency = 440;
  public newFrequency = 440;
  public sampleRate = 44100;
  public volume = 0.0;
  public newVolume = 0.0;
  public dutyCycle = 0;

  public lastSin = 0;
  public phase = 0;

  public echoBuffer = new Float32Array(8192);
  public echoBuffer2 = new Float32Array(8192);
  public echoId = 0;
  public echoMode = 0;
  public noiseTable = new Float32Array(8192);
  public noiseId = 0;

  constructor(sampleRate: number) {
    this.sampleRate = sampleRate;
    for (let i = 0; i < this.noiseTable.length; i++) {
      this.noiseTable[i] = -1 + Math.random() * 2;
    }
  }

  public doSin(): number {
    let v = Math.sin(
      2 * Math.PI * this.frequency * (this.phase / this.sampleRate)
    );
    let isUp = v > this.lastSin;

    if (this.frequency != this.newFrequency) {
      if (Math.abs(v) < 0.1) {
        this.frequency = this.newFrequency;

        // Find next phase
        for (let i = 0; i < 4096; i++) {
          let v2 = Math.sin(
            2 * Math.PI * this.newFrequency * (this.phase / this.sampleRate)
          );
          let v3 = Math.sin(
            2 *
              Math.PI *
              this.newFrequency *
              ((this.phase + 1) / this.sampleRate)
          );

          if (v3 > v2 && isUp && Math.abs(v2) < 0.1) {
            break;
          }
          if (v3 < v2 && !isUp && Math.abs(v2) < 0.1) {
            break;
          }

          this.phase += 1;
        }
      }
    }

    this.lastSin = v;
    return v;
  }

  public doSquare(): number {
    let x = this.doSin();
    if (x > this.dutyCycle) {
      return 1;
    } else {
      return -1;
    }
  }

  public doTriangle(): number {
    let ft = this.frequency * (this.phase / this.sampleRate);
    let v = 4 * Math.abs(ft - Math.floor(ft + 1 / 2)) - 1;

    if (this.frequency != this.newFrequency) {
      if (v < -0.9) {
        this.frequency = this.newFrequency;

        // Find next phase
        for (let i = 0; i < 4096; i++) {
          ft = this.frequency * (this.phase / this.sampleRate);
          let v2 = 4 * Math.abs(ft - Math.floor(ft + 1 / 2)) - 1;
          if (v2 < -0.9) break;
          this.phase += 1;
        }
      }
    }
    return v;
  }

  public doSaw(): number {
    if (this.frequency != this.newFrequency) {
      this.frequency = this.newFrequency;
    }

    let f = this.frequency;
    let tt = this.phase / this.sampleRate;
    return 2 * (tt % (1 / f)) * f - 1;
  }

  public doNoise(): number {
    if (this.frequency != this.newFrequency) {
      this.frequency = this.newFrequency;
    }

    const step = this.frequency / 440;

    this.noiseId += step;
    if (this.noiseId > this.noiseTable.length - 1) {
      this.noiseId = 0;
    }
    return this.noiseTable[~~this.noiseId];
  }

  private lerp(start, end, amt) {
    return (1 - amt) * start + amt * end;
  }

  public do(): number {
    this.phase += 1;

    if (this.volume != this.newVolume) {
      let dir = this.newVolume - this.volume > 0 ? 1 : -1;
      this.volume += 0.0005 * dir;
      if (Math.abs(this.volume - this.newVolume) < 0.01) {
        this.volume = this.newVolume;
      }
    }

    // if (Math.abs(this.volume) < 0.005) return 0;

    let v = 0;

    if (this.waveType == 0) v = this.doSin() * this.volume;
    if (this.waveType == 1) v = this.doSquare() * this.volume;
    if (this.waveType == 2) v = this.doTriangle() * this.volume;
    if (this.waveType == 3) v = this.doSaw() * this.volume;
    if (this.waveType == 4) v = this.doNoise() * this.volume;

    if (this.echoMode === 0) {
      this.echoBuffer[this.echoId] = v * 0.35;
      this.echoId += 1;
      if (this.echoId > this.echoBuffer.length - 1) {
        this.echoId = 0;
        this.echoMode = 1;
      }
    } else {
      this.echoBuffer2[this.echoId] = v * 0.35;
      this.echoId += 1;
      if (this.echoId > this.echoBuffer2.length - 1) {
        this.echoId = 0;
        this.echoMode = 0;
      }
    }

    return v + this.echoBuffer[this.echoId];
  }
}

// @ts-ignore
globalThis.AudioChannel = AudioChannel;
