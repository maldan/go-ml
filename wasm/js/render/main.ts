class GoRenderWasm {
  static wasmTime = [];
  static jsTime = [];
  static calculateTime = [];
  static soundTime = [];

  static stats = {
    requestAnimationFramePerSecond: 0,
    avgDelta: [],
  };

  static async init(wasmData: ArrayBuffer) {
    const go = new Go();
    //const b = await fetch(wasmModuleUrl);
    //const bytes = await b.blob();
    const wasmModule = await WebAssembly.instantiate(wasmData, go.importObject);
    let memory = new ArrayBuffer(0);
    go.run(wasmModule.instance);

    // Init webgl render
    await GoRender.init();

    // Init sound
    // GoSound.init(wasmModule.instance);

    let start;
    let audioTick = 0;

    window.exportWasmMemory = () => {
      const file = new Blob([memory], { type: "bin" });

      const a = document.createElement("a"),
        url = URL.createObjectURL(file);
      a.href = url;
      a.download = "memory.bin";
      document.body.appendChild(a);
      a.click();
      setTimeout(function () {
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
      }, 0);
    };

    const step = (timestamp) => {
      if (start === undefined) start = timestamp;
      let delta = (timestamp - start) / 1000;
      if (delta <= 0) delta = 1 / 1000;

      this.stats.avgDelta.push(delta);

      // Calculate scene in golang
      let pp = performance.now();
      window.go.gameTick(delta);
      this.calculateTime.push(performance.now() - pp);

      pp = performance.now();
      window.go.renderFrame();
      this.wasmTime.push(performance.now() - pp);

      // Get state
      let state = window.go.renderState();

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

      // Reset Mouse click
      // @ts-ignore
      for (let i = 0; i < 4; i++) window.go.setMouseClick(i, false);
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
