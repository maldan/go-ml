class NumberHelper {
  static lerp(start, end, amt) {
    return (1 - amt) * start + amt * end;
  }
}

class MyAudioProcessor extends AudioWorkletProcessor {
  sex = 0;
  echoBuffer = [];
  echoBuffer2 = [];
  echoId = 0;
  echoMode = 0;

  constructor() {
    super();

    this.port.onmessage = this.onmessage.bind(this);
    this.phase = 0;
    this.echoBuffer = new Float32Array(8192);
    this.echoBuffer2 = new Float32Array(8192);
  }

  async onmessage(e) {
    if (e.data.type === "init") {
      eval(e.data.synthText);
      this.ch0 = new AudioChannel(e.data.sampleRate);
      this.ch1 = new AudioChannel(e.data.sampleRate);
    }
    if (e.data.type === "setChannelData") {
      for (let i = 0; i < e.data.channel.length; i++) {
        const ch = e.data.channel[i].id;
        const data = e.data.channel[i].data;

        if (data.setFrequency !== undefined)
          this["ch" + ch].newFrequency = data.setFrequency;
        if (data.setVolume !== undefined)
          this["ch" + ch].newVolume = data.setVolume;
        if (data.setWaveType !== undefined)
          this["ch" + ch].waveType = data.setWaveType;
        if (data.setDutyCycle !== undefined)
          this["ch" + ch].dutyCycle = data.setDutyCycle;
      }
    }
  }

  process(inputList, outputList, parameters) {
    const output = outputList[0];
    const channel = output[0];

    // let echoBuffer = this.echoBuffer2;
    // if (this.echoMode === 1) echoBuffer = this.echoBuffer;

    const outList = [[], []];
    for (let j = 0; j < 2; j++) {
      const ch = this["ch" + j];

      for (let i = 0; i < channel.length; ++i) {
        outList[j].push(ch.do(this.phase));

        // channel[i] = this.ch0.do(this.phase) + this.echoBuffer[this.echoId];

        //channel[i] = this.ch0.do(this.phase);

        // Add echo
        /*if (this.echoMode === 0) {
          this.echoBuffer[this.echoId] = channel[i] * 0.35;
          this.echoId += 1;
          if (this.echoId > this.echoBuffer.length - 1) {
            this.echoId = 0;
            this.echoMode = 1;
          }
        } else {
          this.echoBuffer2[this.echoId] = channel[i] * 0.35;
          this.echoId += 1;
          if (this.echoId > this.echoBuffer2.length - 1) {
            this.echoId = 0;
            this.echoMode = 0;
          }
        }*/

        this.phase += 1;
        /*this.phase += this.ch0.phaseOffset;
        this.ch0.phaseOffset = 0;*/
        ch.volume = NumberHelper.lerp(
          ch.volume,
          ch.newVolume,
          i / channel.length
        );
      }
    }

    // Combine channels
    for (let i = 0; i < channel.length; ++i) {
      channel[i] = outList[0][i] + outList[1][i];
    }

    return true;
  }
}

registerProcessor("my-audio-processor", MyAudioProcessor);
