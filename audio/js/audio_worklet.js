class NumberHelper {
  static lerp(start, end, amt) {
    return (1 - amt) * start + amt * end;
  }
}

class MyAudioProcessor extends AudioWorkletProcessor {
  sex = 0;
  sampleList = {};
  queue = {};

  constructor() {
    super();

    this.port.onmessage = this.onmessage.bind(this);
  }

  async onmessage(e) {
    if (e.data.type === "init") {
      eval(e.data.synthText);
      this.ch0 = new AudioChannel(e.data.sampleRate);
      this.ch1 = new AudioChannel(e.data.sampleRate);
      this.ch2 = new AudioChannel(e.data.sampleRate);
      this.ch3 = new AudioChannel(e.data.sampleRate);

      /* this.sfx0 = new AudioSampleChannel(e.data.sampleRate);
      this.sfx1 = new AudioSampleChannel(e.data.sampleRate);
      this.bgm = new AudioSampleChannel(e.data.sampleRate);*/
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
    if (e.data.type === "loadSample") {
      this.sampleList[e.data.name] = e.data.data;
    }
    if (e.data.type === "playSample") {
      this.queue[e.data.channel] = { offset: 0, sampleName: e.data.name };
      // this.sampleList[e.data.name] = e.data.data;
    }
  }

  process(inputList, outputList, parameters) {
    const output = outputList[0];
    const channel = output[0];

    // let echoBuffer = this.echoBuffer2;
    // if (this.echoMode === 1) echoBuffer = this.echoBuffer;

    const outList = [[], [], [], [], [], [], [], [], []];
    for (let j = 0; j < 4; j++) {
      const ch = this["ch" + j];
      // console.log(this);
      for (let i = 0; i < channel.length; ++i) {
        // outList[j].push(ch.do());
        // channel[i] = this.ch0.do(this.phase) + this.echoBuffer[this.echoId];
        // channel[i] = this.ch0.do(this.phase);
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
      }
    }

    // Sfx
    for (let i = 0; i < channel.length; ++i) {
      if (this.queue["sfx0"]) {
        let sampleName = this.queue["sfx0"].sampleName;
        let offset = this.queue["sfx0"].offset;
        let v = this.sampleList[sampleName][offset];
        outList[4].push(v === undefined ? 0 : v);
        if (v === undefined) {
          delete this.queue["sfx0"];
        } else {
          this.queue["sfx0"].offset += 1;
        }
      } else {
        outList[4].push(0);
      }
      //outList[4].push(this.sfx0.do());
      //outList[5].push(this.sfx1.do());
      //outList[6].push(this.bgm.do());
    }

    // Combine channels
    for (let i = 0; i < channel.length; ++i) {
      channel[i] = outList[4][i];
      /*channel[i] =
        outList[0][i] +
        outList[1][i] +
        outList[2][i] +
        outList[3][i] +
        outList[4][i] +
        outList[5][i] +
        outList[6][i];*/
    }

    return true;
  }
}

registerProcessor("my-audio-processor", MyAudioProcessor);
