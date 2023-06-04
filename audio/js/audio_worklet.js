class NumberHelper {
  static lerp(start, end, amt) {
    return (1 - amt) * start + amt * end;
  }
}

class SamplePlayer {
  main = undefined;
  ch = "";
  loop = false;

  constructor(main, ch) {
    this.main = main;
    this.ch = ch;
  }

  do() {
    let ch = this.ch;

    if (this.main.queue[ch]) {
      let sampleName = this.main.queue[ch].sampleName;
      if (!this.main.sampleList[sampleName]) return 0;

      let offset = this.main.queue[ch].offset;
      let volume = this.main.queue[ch].volume;
      let v = this.main.sampleList[sampleName][~~offset];
      let pitch = this.main.queue[ch].pitch;
      this.main.queue[ch].offset += pitch;

      if (v === undefined) {
        if (this.loop) {
          this.main.queue[ch].offset = 0;
        } else {
          delete this.main.queue[ch];
        }
      }

      return v === undefined ? 0 : v * volume;

      /*outList[4].push(v === undefined ? 0 : v * volume);
      if (v === undefined) {
        delete this.queue["sfx0"];
      } else {
        this.queue["sfx0"].offset += pitch;
      }*/
    } else {
      // outList[4].push(0);
      return 0;
    }
  }
}

class MyAudioProcessor extends AudioWorkletProcessor {
  sex = 0;
  sampleList = {};
  queue = {};
  sfx0 = new SamplePlayer(this, "sfx0");
  sfx1 = new SamplePlayer(this, "sfx1");
  sfx2 = new SamplePlayer(this, "sfx2");
  sfx3 = new SamplePlayer(this, "sfx3");
  bgm = new SamplePlayer(this, "bgm");
  masterVolume = 1;

  constructor() {
    super();

    this.port.onmessage = this.onmessage.bind(this);
    this.bgm.loop = true;
  }

  async onmessage(e) {
    if (e.data.type === "init") {
      eval(e.data.synthText);
      this.ch0 = new AudioChannel(e.data.sampleRate);
      this.ch1 = new AudioChannel(e.data.sampleRate);
      this.ch2 = new AudioChannel(e.data.sampleRate);
      this.ch3 = new AudioChannel(e.data.sampleRate);

      // this.sfx0 = new SamplePlayer(this, "sfx0");

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
    if (e.data.type === "setMasterVolume") {
      this.masterVolume = e.data.volume;
    }
    if (e.data.type === "playSample") {
      // Find free channel
      if (e.data.channel === "sfx?") {
        for (let i = 0; i < 4; i++) {
          if (!this.queue["sfx" + i]) {
            e.data.channel = "sfx" + i;
            break;
          }
        }

        // Not found any free channel
        if (e.data.channel === "sfx?") {
          e.data.channel = "sfx" + ~~(Math.random() * 4);
        }
      }
      this.queue[e.data.channel] = {
        offset: 0,
        sampleName: e.data.name,
        pitch: e.data.pitch === undefined ? 1 : e.data.pitch,
        volume: e.data.volume === undefined ? 1 : e.data.volume,
      };
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

    // Init
    for (let i = 0; i < channel.length; ++i) {
      channel[i] = 0;
    }

    // Sfx
    for (let i = 0; i < channel.length; ++i) {
      channel[i] += this.sfx0.do();
      channel[i] += this.sfx1.do();
      channel[i] += this.sfx2.do();
      channel[i] += this.sfx3.do();
      channel[i] += this.bgm.do();
      channel[i] *= this.masterVolume;
    }

    /*for (let i = 0; i < channel.length; ++i) {
      if (this.queue["sfx0"]) {
        let sampleName = this.queue["sfx0"].sampleName;
        let offset = this.queue["sfx0"].offset;
        let volume = this.queue["sfx0"].volume;
        let v = this.sampleList[sampleName][~~offset];
        let pitch = this.queue["sfx0"].pitch;

        outList[4].push(v === undefined ? 0 : v * volume);
        if (v === undefined) {
          delete this.queue["sfx0"];
        } else {
          this.queue["sfx0"].offset += pitch;
        }
      } else {
        outList[4].push(0);
      }
    }*/

    // Combine channels
    for (let i = 0; i < channel.length; ++i) {
      // channel[i] = outList[4][i];
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
