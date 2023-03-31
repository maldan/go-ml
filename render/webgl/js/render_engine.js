class GoRender {
  static _gl = null;

  static data = {
    vertexList: [],
    indexList: [],
    positionList: [],
    projectionMatrix: [],
    pointList: [],
  };
  static programInfo = {
    mainShader: undefined,
    pointShader: undefined,
    attribLocations: {
      pointPosition: 0,
      vertexPosition: 0,
      colorPosition: 0,
      modelPosition: 0,
    },
    uniformLocations: {
      projectionMatrix: 0,
    },
    buffer: {
      point: 0,
      vertexPosition: 0,
      index: 0,
      modelPosition: 0,
    },
  };
  static shaderSource = {};

  static setWasmData(memory, state) {
    let byteArray = new Uint8Array(memory);
    let shortArray = new Uint16Array(memory);
    let float32Array = new Float32Array(memory);

    GoRender.data.positionList = float32Array.subarray(
      state.positionArrayPointer / 4,
      state.positionArrayPointer / 4 + state.vertexArrayLength
    );
    GoRender.data.vertexList = float32Array.subarray(
      state.vertexArrayPointer / 4,
      state.vertexArrayPointer / 4 + state.vertexArrayLength
    );
    GoRender.data.indexList = shortArray.subarray(
      state.indexArrayPointer / 2,
      state.indexArrayPointer / 2 + state.indexArrayLength
    );
    GoRender.data.projectionMatrix = float32Array.subarray(
      state.projectionMatrixPointer / 4,
      state.projectionMatrixPointer / 4 + 16
    );
    GoRender.data.pointList = float32Array.subarray(
      state.pointArrayPointer / 4,
      state.pointArrayPointer / 4 + state.pointArrayLength
    );
  }

  static loadShader(type, source) {
    const shader = this._gl.createShader(type);
    this._gl.shaderSource(shader, source);
    this._gl.compileShader(shader);

    if (!this._gl.getShaderParameter(shader, this._gl.COMPILE_STATUS)) {
      const info = this._gl.getShaderInfoLog(shader);
      this._gl.deleteShader(shader);
      throw new Error(`An error occurred compiling the shaders: ${info}`);
    }

    return shader;
  }

  static compileShader(vertex, fragment) {
    const vertexShader = this.loadShader(this._gl.VERTEX_SHADER, vertex);
    const fragmentShader = this.loadShader(this._gl.FRAGMENT_SHADER, fragment);

    const shaderProgram = this._gl.createProgram();
    this._gl.attachShader(shaderProgram, vertexShader);
    this._gl.attachShader(shaderProgram, fragmentShader);
    this._gl.linkProgram(shaderProgram);

    if (!this._gl.getProgramParameter(shaderProgram, this._gl.LINK_STATUS)) {
      const info = this._gl.getProgramInfoLog(shaderProgram);
      console.log(info);
      throw new Error(`Unable to initialize the shader program: ${info}`);
    }

    return shaderProgram;
  }

  static async loadShaderSourceCode(url) {
    const x = await fetch(url);
    this.shaderSource[url] = await x.text();
  }

  static async init() {
    const canvas = document.querySelector("#glcanvas");
    this._gl = canvas.getContext("webgl");

    if (this._gl === null) throw new Error("WebGL is not supported");

    // Load shaders
    const shaderList = [
      "matrix.glsl",
      "main.vertex.glsl",
      "main.fragment.glsl",
      "point.vertex.glsl",
      "point.fragment.glsl",
    ];
    for (let i = 0; i < shaderList.length; i++) {
      await this.loadShaderSourceCode(`./shader/${shaderList[i]}`);
    }

    // Inject library
    shaderList.slice(1).forEach((x) => {
      this.shaderSource[`./shader/${x}`] = this.shaderSource[
        `./shader/${x}`
      ].replace("// LIB", this.shaderSource["./shader/matrix.glsl"]);
    });

    // Compile shaders
    const mainShader = this.compileShader(
      this.shaderSource["./shader/main.vertex.glsl"],
      this.shaderSource["./shader/main.fragment.glsl"]
    );
    const pointShader = this.compileShader(
      this.shaderSource["./shader/point.vertex.glsl"],
      this.shaderSource["./shader/point.fragment.glsl"]
    );

    this.programInfo = {
      mainShader,
      pointShader,
      attribLocations: {
        pointPosition: this._gl.getAttribLocation(
          pointShader,
          "aVertexPosition"
        ),
        vertexPosition: this._gl.getAttribLocation(
          mainShader,
          "aVertexPosition"
        ),
        colorPosition: this._gl.getAttribLocation(mainShader, "aVertexColor"),
        modelPosition: this._gl.getAttribLocation(mainShader, "aModelPosition"),
      },
      uniformLocations: {
        projectionMatrix: this._gl.getUniformLocation(
          mainShader,
          "uProjectionMatrix"
        ),
        projectionMatrix2: this._gl.getUniformLocation(
          pointShader,
          "uProjectionMatrix"
        ),
      },
      buffer: {
        vertexPosition: this._gl.createBuffer(),
        index: this._gl.createBuffer(),
        modelPosition: this._gl.createBuffer(),
        point: this._gl.createBuffer(),
      },
    };
  }

  static uploadData(type, buffer, data) {
    if (type === "element") {
      this._gl.bindBuffer(this._gl.ELEMENT_ARRAY_BUFFER, buffer);
      this._gl.bufferData(
        this._gl.ELEMENT_ARRAY_BUFFER,
        data,
        this._gl.STATIC_DRAW
      );
    } else {
      this._gl.bindBuffer(this._gl.ARRAY_BUFFER, buffer);
      this._gl.bufferData(this._gl.ARRAY_BUFFER, data, this._gl.STATIC_DRAW);
    }

    // Unbind the buffer
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, null);
  }

  static enableAttribute(buffer, attribute) {
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, buffer);
    this._gl.vertexAttribPointer(attribute, 3, this._gl.FLOAT, false, 0, 0);
    this._gl.enableVertexAttribArray(attribute);
  }

  static drawMain() {
    this._gl.clearColor(0.0, 0.0, 0.0, 1.0);
    this._gl.clearDepth(1.0);
    this._gl.enable(this._gl.DEPTH_TEST);
    this._gl.depthFunc(this._gl.LEQUAL);
    this._gl.clear(this._gl.COLOR_BUFFER_BIT | this._gl.DEPTH_BUFFER_BIT);

    // Put main data
    this.uploadData(
      "element",
      this.programInfo.buffer.index,
      this.data.indexList
    );
    this.uploadData(
      "any",
      this.programInfo.buffer.vertexPosition,
      this.data.vertexList
    );
    this.uploadData(
      "any",
      this.programInfo.buffer.modelPosition,
      this.data.positionList
    );

    // Enable attributes
    this.enableAttribute(
      this.programInfo.buffer.vertexPosition,
      this.programInfo.attribLocations.vertexPosition
    );
    this.enableAttribute(
      this.programInfo.buffer.modelPosition,
      this.programInfo.attribLocations.modelPosition
    );

    this._gl.useProgram(this.programInfo.mainShader);
    this.programInfo.uniformLocations.projectionMatrix =
      this._gl.getUniformLocation(
        this.programInfo.mainShader,
        "uProjectionMatrix"
      );
    this._gl.uniformMatrix4fv(
      this.programInfo.uniformLocations.projectionMatrix,
      false,
      this.data.projectionMatrix
    );

    this._gl.bindBuffer(
      this._gl.ELEMENT_ARRAY_BUFFER,
      this.programInfo.buffer.index
    );

    this._gl.drawElements(
      this._gl.TRIANGLES,
      this.data.indexList.length,
      this._gl.UNSIGNED_SHORT,
      0
    );

    // Draw points
    this.uploadData("any", this.programInfo.buffer.point, this.data.pointList);
    this.enableAttribute(
      this.programInfo.buffer.point,
      this.programInfo.attribLocations.pointPosition
    );

    this._gl.useProgram(this.programInfo.pointShader);
    this.programInfo.uniformLocations.projectionMatrix2 =
      this._gl.getUniformLocation(
        this.programInfo.pointShader,
        "uProjectionMatrix"
      );
    this._gl.uniformMatrix4fv(
      this.programInfo.uniformLocations.projectionMatrix2,
      false,
      this.data.projectionMatrix
    );

    this._gl.disable(this._gl.DEPTH_TEST);
    this._gl.drawArrays(this._gl.POINTS, 0, this.data.pointList.length / 3);
  }
}

class GoRenderWasmLayer {
  static async init(wasmModuleUrl) {
    const go = new Go();
    const b = await fetch(wasmModuleUrl);
    const bytes = await b.blob();
    const wasmModule = await WebAssembly.instantiate(
      await bytes.arrayBuffer(),
      go.importObject
    );
    go.run(wasmModule.instance);

    // Init webgl render
    await GoRender.init();

    let start;

    function step(timestamp) {
      if (start === undefined) start = timestamp;

      const elapsed = timestamp - start;

      // Calculate scene in golang
      goWasmRenderFrame();

      // Get state
      let state = goWasmRenderState();

      // Send golang data to webgl render
      GoRender.setWasmData(wasmModule.instance.exports.mem.buffer, state);

      // Draw scene
      GoRender.drawMain();

      window.requestAnimationFrame(step);
    }

    window.requestAnimationFrame(step);
  }
}
