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

class GoRenderWasm {
  static wasmTime = [];
  static jsTime = [];
  static calculateTime = [];

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
    console.log(wasmModule.instance);

    // Init webgl render
    await GoRender.init();

    let start;

    const step = (timestamp) => {
      if (start === undefined) start = timestamp;
      let delta = (timestamp - start) / 1000;
      if (delta <= 0) delta = 1 / 1000;

      this.beforeFrame(delta);

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

      this.afterFrame(delta);

      start = timestamp;
      window.requestAnimationFrame(step);
    };

    window.requestAnimationFrame(step);

    // Timers
    const avg = (x) => x.reduce((a, b) => a + b) / x.length;
    setInterval(() => {
      document.getElementById("wasm").innerHTML = `render calculate: ${avg(
        this.wasmTime
      ).toFixed(2)}`;
      document.getElementById("js").innerHTML = `render draw: ${avg(
        this.jsTime
      ).toFixed(2)}`;
      document.getElementById("calculate").innerHTML = `game calculate: ${avg(
        this.calculateTime
      ).toFixed(2)}`;

      document.getElementById("mem").innerHTML = `mem: ${(
        memory.byteLength / 1048576
      ).toFixed(3)} mb`;

      this.wasmTime.length = 0;
      this.jsTime.length = 0;
      this.calculateTime.length = 0;
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