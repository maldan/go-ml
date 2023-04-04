class GoRender {
  static _gl = null;

  static shaderSource = {};
  static layerList = [];

  static async loadShaderSourceCode(url) {
    const x = await fetch(url);
    this.shaderSource[url] = await x.text();
  }

  static async init() {
    const canvas = document.querySelector("#glcanvas");
    this._gl = canvas.getContext("webgl", { antialias: false });
    if (this._gl === null) throw new Error("WebGL is not supported");

    // Load shaders
    const shaderList = ["matrix.glsl"];
    shaderList.push(
      ...["main", "point", "line"]
        .map((x) => [`${x}.vertex.glsl`, `${x}.fragment.glsl`])
        .flat()
    );
    for (let i = 0; i < shaderList.length; i++) {
      await this.loadShaderSourceCode(`./shader/${shaderList[i]}`);
    }

    // Inject library
    shaderList.slice(1).forEach((x) => {
      this.shaderSource[`./shader/${x}`] = this.shaderSource[
        `./shader/${x}`
      ].replace("// LIB", this.shaderSource["./shader/matrix.glsl"]);
    });

    // Main texture
    const texture = loadTexture(this._gl, "./texture.png");

    // Compile shaders
    this.layerList = [
      new GoRenderLayer("main", this._gl),
      new GoRenderPointLayer("point", this._gl),
      new GoRenderLineLayer("line", this._gl),
    ].map((x) => {
      x.init(
        this.shaderSource[`./shader/${x.name}.vertex.glsl`],
        this.shaderSource[`./shader/${x.name}.fragment.glsl`]
      );
      x.texture = texture;
      return x;
    });
  }

  static draw() {
    this._gl.clearColor(0.0, 0.0, 0.0, 1.0);
    this._gl.clearDepth(1.0);
    this._gl.enable(this._gl.DEPTH_TEST);

    this._gl.enable(this._gl.CULL_FACE);
    this._gl.cullFace(this._gl.BACK);

    this._gl.depthFunc(this._gl.LEQUAL);
    this._gl.clear(this._gl.COLOR_BUFFER_BIT | this._gl.DEPTH_BUFFER_BIT);

    this.layerList.forEach((layer) => {
      layer.draw();
    });
  }
}

class GoSound {
  static _audioContext;
  static _time = -1;
  static _bigBuffer = [];
  // static _privateData = [];

  static init(wasmInstance, getState, tick) {
    this._audioContext = new (window.AudioContext || window.webkitAudioContext)(
      {
        latencyHint: "interactive",
        sampleRate: 44100 / 4,
      }
    );

    const script_processor = this._audioContext.createScriptProcessor(
      256,
      0,
      1
    );
    script_processor.onaudioprocess = (event) => {
      const dst = event.outputBuffer;
      const dst_l = dst.getChannelData(0);

      if (this._bigBuffer.length < 256) {
        tick(32 / 1000);
        this.capture(wasmInstance, getState);
      }

      const arr = this._bigBuffer.slice(0, 256);
      for (let i = 0; i < arr.length; i++) {
        dst_l[i] = arr[i];
      }
      this._bigBuffer.splice(0, 256);
    };
    script_processor.connect(this._audioContext.destination);

    // tick(16 / 1000);
    setInterval(() => {
      /* const state = getState();
      let memory = null;
      if (wasmInstance.exports.mem) memory = wasmInstance.exports.mem.buffer;
      if (wasmInstance.exports.memory)
        memory = wasmInstance.exports.memory.buffer;

      const float32Array = new Float32Array(memory);
      const sampleData = float32Array.subarray(
        state.bufferPointer / 4,
        state.bufferPointer / 4 + state.bufferPosition
      );*/
      tick(32 / 1000);
      this.capture(wasmInstance, getState);
      // this.play(wasmInstance, getState);
    }, 32);

    /*this._audioBuffer = this._audioContext.createBuffer(
      1,
      this._sampleRate,
      this._sampleRate
    );

    this._source = this._audioContext.createBufferSource();
    this._source.buffer = this._audioBuffer;
    this._source.loop = true;

    this._source.addEventListener("ended", () => {
      console.log("gas");
    });

    document.addEventListener("click", () => {
      this._source.connect(this._audioContext.destination);
    });*/

    /*document.addEventListener("click", () => {
      this.isPlay = false;
      this.play([0, 1, 2]);
    });*/
    /*tick();
    const soundState = goWasmSoundState();

    const float32Array = new Float32Array(memory);
    const audioBufferData = float32Array.subarray(
        soundState.bufferPointer / 4,
        soundState.bufferPointer / 4 + soundState.bufferLength
    );

    GoSound.play(audioBufferData);*/
    // this.play();
  }

  static capture(wasmInstance, getState) {
    const state = getState();
    let memory = null;
    if (wasmInstance.exports.mem) memory = wasmInstance.exports.mem.buffer;
    if (wasmInstance.exports.memory)
      memory = wasmInstance.exports.memory.buffer;

    const float32Array = new Float32Array(memory);
    const sampleData = float32Array.subarray(
      state.bufferPointer / 4,
      state.bufferPointer / 4 + state.bufferPosition
    );
    // console.log(sampleData.length);
    this._bigBuffer.push(...sampleData);

    return sampleData;
  }

  /*static play(wasmInstance, getState) {
    if (this._time === -1) this._time = this._audioContext.currentTime;
    const state = getState();
    let memory = null;
    if (wasmInstance.exports.mem) memory = wasmInstance.exports.mem.buffer;
    if (wasmInstance.exports.memory)
      memory = wasmInstance.exports.memory.buffer;

    const float32Array = new Float32Array(memory);
    const sampleData = float32Array.subarray(
      state.bufferPointer / 4,
      state.bufferPointer / 4 + state.bufferPosition
    );

    // Create buffer
    const audioBuffer = this._audioContext.createBuffer(
      1,
      sampleData.length + 100,
      state.sampleRate
    );

    // Copy sound
    audioBuffer.getChannelData(0).set(sampleData);

    const source = this._audioContext.createBufferSource();
    source.buffer = audioBuffer;
    source.connect(this._audioContext.destination);
    source.start(this._time, 0, audioBuffer.duration);

    this._time += source.buffer.duration;
  }*/
}

class GoRenderWasm {
  static wasmTime = [];
  static jsTime = [];
  static calculateTime = [];
  static soundTime = [];

  static stats = {
    requestAnimationFramePerSecond: 0,
    avgDelta: [],
  };

  static beforeFrame = () => {};
  static afterFrame = () => {};

  static async init(wasmModuleUrl) {
    const go = new Go();
    const b = await fetch(wasmModuleUrl);
    const bytes = await b.blob();
    const wasmModule = await WebAssembly.instantiate(
      await bytes.arrayBuffer(),
      go.importObject
    );
    let memory = new ArrayBuffer(0);
    go.run(wasmModule.instance);

    // Init webgl render
    await GoRender.init();

    // Init sound
    GoSound.init(wasmModule.instance, goWasmSoundState, goWasmSoundTick);

    let start;
    let audioTick = 0;

    const step = (timestamp) => {
      if (start === undefined) start = timestamp;
      let delta = (timestamp - start) / 1000;
      if (delta <= 0) delta = 1 / 1000;

      this.stats.avgDelta.push(delta);

      // Calculate scene in golang
      let pp = performance.now();
      goWasmGameTick(delta);
      this.calculateTime.push(performance.now() - pp);

      pp = performance.now();
      goWasmRenderFrame();
      this.wasmTime.push(performance.now() - pp);

      // Get state
      let state = goWasmRenderState();

      // Send golang data to webgl render
      GoRender.layerList.forEach((layer) => {
        layer.state = state;

        if (wasmModule.instance.exports.mem) {
          memory = wasmModule.instance.exports.mem.buffer;
          layer.setWasmData(
            wasmModule.instance.exports.mem.buffer,
            state[layer.name + "Layer"]
          );
        } else {
          memory = wasmModule.instance.exports.memory.buffer;
          layer.setWasmData(
            wasmModule.instance.exports.memory.buffer,
            state[layer.name + "Layer"]
          );
        }
      });

      // Draw scene
      pp = performance.now();
      GoRender.draw();
      this.jsTime.push(performance.now() - pp);

      /*pp = performance.now();
      if (!GoSound.isPlay) {
        if (goWasmSoundTick) {
          goWasmSoundTick();
        }
        if (goWasmSoundState) {
          const soundState = goWasmSoundState();

          const float32Array = new Float32Array(memory);
          const audioBufferData = float32Array.subarray(
            soundState.bufferPointer / 4,
            soundState.bufferPointer / 4 + soundState.bufferLength
          );

          GoSound.play(audioBufferData);
        }
      }
      this.soundTime.push(performance.now() - pp);*/

      // GoSound._wasmMemory = memory;

      this.stats.requestAnimationFramePerSecond += 1;

      /*audioTick += 1;
      if (audioTick > 5) {
        audioTick = 0;

        GoSound.play(wasmModule.instance, goWasmSoundState);
      }*/

      // Request next frame
      start = timestamp;
      window.requestAnimationFrame(step);
    };

    window.requestAnimationFrame(step);

    // Timers
    const avg = (x) => x.reduce((a, b) => a + b) / x.length;
    setInterval(() => {
      const gameCalculate = avg(this.calculateTime);
      const renderCalculate = avg(this.wasmTime);
      const renderDraw = avg(this.jsTime);
      const avgDelta = avg(this.stats.avgDelta);
      // const soundTime = avg(this.soundTime);

      document.getElementById("stats").innerHTML = `
        <div>game calculate: ${gameCalculate.toFixed(2)}</div>
        <div>render calculate: ${renderCalculate.toFixed(2)}</div>
        <div>render draw: ${renderDraw.toFixed(2)}</div> 
        
        <div>total: ${(gameCalculate + renderCalculate + renderDraw).toFixed(
          2
        )}</div>
        <div>mem usage: ${(memory.byteLength / 1048576).toFixed(3)} mb</div>
        
        <div>fps: ${this.stats.requestAnimationFramePerSecond}</div>
        <div>delta: ${avgDelta.toFixed(4)}</div>
      `;

      this.stats.avgDelta.length = 0;
      this.stats.requestAnimationFramePerSecond = 0;
      this.wasmTime.length = 0;
      this.jsTime.length = 0;
      this.calculateTime.length = 0;
      this.soundTime.length = 0;
    }, 1000);
  }
}

function loadTexture(gl, url) {
  const texture = gl.createTexture();
  gl.bindTexture(gl.TEXTURE_2D, texture);

  const level = 0;
  const internalFormat = gl.RGBA;
  const width = 1;
  const height = 1;
  const border = 0;
  const srcFormat = gl.RGBA;
  const srcType = gl.UNSIGNED_BYTE;
  const pixel = new Uint8Array([0, 0, 255, 255]); // opaque blue
  gl.texImage2D(
    gl.TEXTURE_2D,
    level,
    internalFormat,
    width,
    height,
    border,
    srcFormat,
    srcType,
    pixel
  );

  const image = new Image();
  image.onload = () => {
    gl.bindTexture(gl.TEXTURE_2D, texture);
    gl.texImage2D(
      gl.TEXTURE_2D,
      level,
      internalFormat,
      srcFormat,
      srcType,
      image
    );
    gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);

    // gl.generateMipmap(gl.TEXTURE_2D);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
  };
  image.src = url;

  return texture;
}
