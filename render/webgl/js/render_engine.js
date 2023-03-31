class GoRenderLayer {
  _gl = null;

  shader = null;
  attributeList = {};
  uniformList = {};
  bufferList = {};
  dataList = {};
  name = "";

  constructor(gl) {
    this._gl = gl;
  }

  init(vertex, fragment) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex", "index", "position"].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    ["aVertex", "aPosition"].forEach((x) => {
      this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
    });
    ["uProjectionMatrix"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });
  }

  loadShader(type, source) {
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

  compileShader(vertex, fragment) {
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

  setWasmData(memory, state) {
    let byteArray = new Uint8Array(memory);
    let shortArray = new Uint16Array(memory);
    let float32Array = new Float32Array(memory);

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("index", state, shortArray);
    this.setDataArray("projectionMatrix", state, float32Array, 16);

    // this.setDataArray("main", "index", state, shortArray);

    /*GoRender.data.mainLayer.vertexList = float32Array.subarray(
      state.mainLayer.vertexPointer / 4,
      state.mainLayer.vertexPointer / 4 + state.mainLayer.vertexAmount
    );
    GoRender.data.mainLayer.indexList = float32Array.subarray(
      state.mainLayer.indexPointer / 4,
      state.mainLayer.indexPointer / 4 + state.mainLayer.indexAmount
    );*/

    /*GoRender.data.positionList = float32Array.subarray(
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
    );*/
  }

  setDataArray(name, state, array, length = 0) {
    let offsetSize = 1;
    if (array instanceof Uint16Array) offsetSize = 2;
    if (array instanceof Float32Array) offsetSize = 4;

    this.dataList[name] = array.subarray(
      state[name + "Pointer"] / offsetSize,
      state[name + "Pointer"] / offsetSize + (length ?? state[name + "Amount"])
    );
  }

  uploadData(type, name) {
    if (type === "element") {
      this._gl.bindBuffer(this._gl.ELEMENT_ARRAY_BUFFER, this.bufferList[name]);
      this._gl.bufferData(
        this._gl.ELEMENT_ARRAY_BUFFER,
        this.dataList[name],
        this._gl.STATIC_DRAW
      );
    } else {
      this._gl.bindBuffer(this._gl.ARRAY_BUFFER, this.bufferList[name]);
      this._gl.bufferData(
        this._gl.ARRAY_BUFFER,
        this.dataList[name],
        this._gl.STATIC_DRAW
      );
    }

    // Unbind the buffer
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, null);
  }

  enableAttribute(name) {
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, this.bufferList[name]);
    let attr = this._gl.getAttribLocation(
      this.shader,
      this.attributeList["a" + name[0].toUpperCase() + name.slice(1)]
    );
    this._gl.vertexAttribPointer(attr, 3, this._gl.FLOAT, false, 0, 0);
    this._gl.enableVertexAttribArray(attr);
  }

  setUniform(name) {
    let uniformName = "u" + name[0].toUpperCase() + name.slice(1);
    this._gl.uniformMatrix4fv(
      this.uniformList[uniformName],
      false,
      this.dataList[name]
    );
  }

  draw() {
    // Put main data
    this.uploadData("element", "index");
    this.uploadData("any", "vertex");
    this.uploadData("any", "position");

    // Enable attributes
    this.enableAttribute("vertex");
    this.enableAttribute("position");

    this._gl.useProgram(this.shader);

    this.setUniform("projectionMatrix");

    this._gl.bindBuffer(
      this._gl.ELEMENT_ARRAY_BUFFER,
      this.bufferList["index"]
    );

    this._gl.drawElements(
      this._gl.TRIANGLES,
      this.dataList["index"].length,
      this._gl.UNSIGNED_SHORT,
      0
    );
  }
}

class GoRender {
  static _gl = null;

  /*static data = {
    mainLayer: {
      shader: undefined,
      attribute: {
        vertex: 0,
        position: 0,
        rotation: 0,
        scale: 0,
      },
      uniform: {
        projectionMatrix: 0,
      },
      buffer: {
        vertex: 0,
        index: 0,
        position: 0,
        rotation: 0,
        scale: 0,
      },

      vertexList: [],
      indexList: [],
      positionList: [],
      rotationList: [],
      scaleList: [],
      projectionMatrix: [],
    },
    pointLayer: {
      shader: undefined,

      pointList: [],
      projectionMatrix: [],
    },
  };*/

  static shaderSource = {};
  static layerList = [];

  /*  static loadShader(type, source) {
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
  }*/

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
    this.layerList = [new GoRenderLayer(this._gl)];
    this.layerList[0].name = "main";
    this.layerList[0].init(
      this.shaderSource["./shader/main.vertex.glsl"],
      this.shaderSource["./shader/main.fragment.glsl"]
    );

    /*const mainShader = this.compileShader(
      this.shaderSource["./shader/main.vertex.glsl"],
      this.shaderSource["./shader/main.fragment.glsl"]
    );
    const pointShader = this.compileShader(
      this.shaderSource["./shader/point.vertex.glsl"],
      this.shaderSource["./shader/point.fragment.glsl"]
    );*/

    /*this.programInfo = {
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
    };*/
  }

  /*
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

  static enableAttribute(shader, buffer, attributeName) {
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, buffer);

    let attr = this._gl.getAttribLocation(shader, attributeName);
    this._gl.vertexAttribPointer(attr, 3, this._gl.FLOAT, false, 0, 0);
    this._gl.enableVertexAttribArray(attr);
  }

  static setUniform(shader, uniformName, data) {
    let uniformLocation = this._gl.getUniformLocation(shader, uniformName);
    this._gl.uniformMatrix4fv(uniformLocation, false, data);
  }
*/

  static drawMain() {
    // Put main data
    this.uploadData(
      "element",
      this.data.mainLayer.buffer.index,
      this.data.mainLayer.indexList
    );
    this.uploadData(
      "any",
      this.data.mainLayer.buffer.vertex,
      this.data.mainLayer.vertexList
    );
    this.uploadData(
      "any",
      this.data.mainLayer.buffer.position,
      this.data.mainLayer.positionList
    );

    // Enable attributes
    this.enableAttribute(
      this.data.mainLayer.shader,
      this.data.mainLayer.buffer.vertex,
      "aVertexPosition"
    );
    this.enableAttribute(
      this.data.mainLayer.shader,
      this.data.mainLayer.buffer.position,
      "aModelPosition"
    );

    this._gl.useProgram(this.data.mainLayer.shader);

    this.setUniform(
      this.data.mainLayer.shader,
      "uProjectionMatrix",
      this.data.mainLayer.uniform.projectionMatrix
    );

    this._gl.bindBuffer(
      this._gl.ELEMENT_ARRAY_BUFFER,
      this.data.mainLayer.buffer.index
    );

    this._gl.drawElements(
      this._gl.TRIANGLES,
      this.data.indexList.length,
      this._gl.UNSIGNED_SHORT,
      0
    );
  }

  static drawPoints() {
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

  static draw() {
    this._gl.clearColor(0.0, 0.0, 0.0, 1.0);
    this._gl.clearDepth(1.0);
    this._gl.enable(this._gl.DEPTH_TEST);
    this._gl.depthFunc(this._gl.LEQUAL);
    this._gl.clear(this._gl.COLOR_BUFFER_BIT | this._gl.DEPTH_BUFFER_BIT);

    this.layerList[0].draw();

    // this.drawMain();
    // this.drawPoints();
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
      GoRender.layerList.forEach((x) => {
        x.setWasmData(
          wasmModule.instance.exports.mem.buffer,
          state[x.name + "Layer"]
        );
      });

      // Draw scene
      GoRender.draw();

      window.requestAnimationFrame(step);
    }

    window.requestAnimationFrame(step);
  }
}
