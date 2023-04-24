declare const Go: any;

class GoRenderWasm {
  static wasmTime: number[] = [];
  static jsTime: number[] = [];
  static calculateTime: number[] = [];
  static soundTime: number[] = [];

  static stats = {
    requestAnimationFramePerSecond: 0,
    avgDelta: [] as number[],
  };

  static async init(wasmData: ArrayBuffer) {
    const go = new Go();
    const wasmModule = await WebAssembly.instantiate(wasmData, go.importObject);
    let memory = new ArrayBuffer(0);
    go.run(wasmModule.instance);

    // Init webgl render
    await GoRender.init();

    let start = 0;

    (window as any).exportWasmMemory = () => {
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

    const step = (timestamp: number) => {
      if (start === undefined) start = timestamp;
      let delta = (timestamp - start) / 1000;
      if (delta <= 0) delta = 1 / 1000;

      this.stats.avgDelta.push(delta);

      // Calculate scene in golang
      let pp = performance.now();
      (window as any).go.gameTick(delta);
      this.calculateTime.push(performance.now() - pp);

      pp = performance.now();
      (window as any).go.renderFrame();
      this.wasmTime.push(performance.now() - pp);

      // Get state
      let state = (window as any).go.renderState();

      if (wasmModule.instance.exports.mem) {
        // @ts-ignore
        memory = wasmModule.instance.exports.mem.buffer;
      } else {
        // @ts-ignore
        memory = wasmModule.instance.exports.memory.buffer;
      }

      // Send golang data to webgl render
      GoRender.layerList.forEach((layer) => {
        layer.state = state;

        if (wasmModule.instance.exports.mem) {
          layer.setWasmData(
            // @ts-ignore
            wasmModule.instance.exports.mem.buffer,
            state[layer.name + "Layer"]
          );
        } else {
          layer.setWasmData(
            // @ts-ignore
            wasmModule.instance.exports.memory.buffer,
            state[layer.name + "Layer"]
          );
        }
      });

      // Draw scene
      pp = performance.now();
      GoRender.draw();
      this.jsTime.push(performance.now() - pp);

      this.stats.requestAnimationFramePerSecond += 1;

      // Request next frame
      start = timestamp;
      window.requestAnimationFrame(step);

      // Reset Mouse click
      // @ts-ignore
      for (let i = 0; i < 4; i++) window.go.setMouseClick(i, false);
    };

    window.requestAnimationFrame(step);

    // Timers
    const avg = (x: number[]) =>
      x.reduce((a: number, b: number) => a + b) / x.length;
    setInterval(() => {
      const gameCalculate = avg(this.calculateTime);
      const renderCalculate = avg(this.wasmTime);
      const renderDraw = avg(this.jsTime);
      const avgDelta = avg(this.stats.avgDelta);
      // const soundTime = avg(this.soundTime);

      const stats = document.getElementById("stats");
      if (stats) {
        stats.innerHTML = `
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
      }

      this.stats.avgDelta.length = 0;
      this.stats.requestAnimationFramePerSecond = 0;
      this.wasmTime.length = 0;
      this.jsTime.length = 0;
      this.calculateTime.length = 0;
      this.soundTime.length = 0;
    }, 1000);
  }
}
