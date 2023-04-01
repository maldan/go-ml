class GoRenderLayer {
  _gl = null;

  shader = null;
  attributeList = {};
  uniformList = {};
  bufferList = {};
  dataList = {};
  name = "";
  state = {};

  constructor(name, gl) {
    this.name = name;
    this._gl = gl;
  }

  init(vertex, fragment) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex", "index", "position", "rotation"].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    ["aVertex", "aPosition", "aRotation"].forEach((x) => {
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
    this.setDataArray("position", state, float32Array);
    this.setDataArray("rotation", state, float32Array);
    this.setDataArray("index", state, shortArray);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  setDataArray(name, state, array, length = 0) {
    let offsetSize = 1;
    if (array instanceof Uint16Array) offsetSize = 2;
    if (array instanceof Float32Array) offsetSize = 4;

    this.dataList[name] = array.subarray(
      ~~state[name + "Pointer"] / offsetSize,
      ~~state[name + "Pointer"] / offsetSize +
        (length || ~~state[name + "Amount"])
    );
  }

  uploadData(type, name) {
    if (type === "element") {
      this._gl.bindBuffer(this._gl.ELEMENT_ARRAY_BUFFER, this.bufferList[name]);
      this._gl.bufferData(
        this._gl.ELEMENT_ARRAY_BUFFER,
        this.dataList[name],
        this._gl.DYNAMIC_DRAW
      );
    } else {
      this._gl.bindBuffer(this._gl.ARRAY_BUFFER, this.bufferList[name]);
      this._gl.bufferData(
        this._gl.ARRAY_BUFFER,
        this.dataList[name],
        this._gl.DYNAMIC_DRAW
      );
    }

    // Unbind the buffer
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, null);
  }

  enableAttribute(name) {
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, this.bufferList[name]);

    let attributeName = "a" + name[0].toUpperCase() + name.slice(1);

    this._gl.vertexAttribPointer(
      this.attributeList[attributeName],
      3,
      this._gl.FLOAT,
      false,
      0,
      0
    );
    this._gl.enableVertexAttribArray(this.attributeList[attributeName]);
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
    // Set program
    this._gl.useProgram(this.shader);

    // Put main data
    this.uploadData("element", "index");
    this.uploadData("any", "vertex");
    this.uploadData("any", "position");
    this.uploadData("any", "rotation");

    // Enable attributes
    this.enableAttribute("vertex");
    this.enableAttribute("position");
    this.enableAttribute("rotation");

    // Set projection
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

class GoRenderPointLayer extends GoRenderLayer {
  init(vertex, fragment) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex"].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    ["aVertex"].forEach((x) => {
      this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
    });
    ["uProjectionMatrix"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });
  }

  setWasmData(memory, state) {
    let float32Array = new Float32Array(memory);

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  draw() {
    // Set program
    this._gl.useProgram(this.shader);

    // Draw points
    this.uploadData("any", "vertex");
    this.enableAttribute("vertex");

    // Set projection
    this.setUniform("projectionMatrix");

    this._gl.disable(this._gl.DEPTH_TEST);
    this._gl.drawArrays(this._gl.POINTS, 0, this.dataList["vertex"].length / 3);
  }
}

class GoRenderLineLayer extends GoRenderLayer {
  init(vertex, fragment) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex", "color"].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    ["aVertex", "aColor"].forEach((x) => {
      this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
    });
    ["uProjectionMatrix"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });
  }

  setWasmData(memory, state) {
    let float32Array = new Float32Array(memory);

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("color", state, float32Array);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  draw() {
    // Set program
    this._gl.useProgram(this.shader);

    // Draw points
    this.uploadData("any", "vertex");
    this.uploadData("any", "color");
    this.enableAttribute("vertex");
    this.enableAttribute("color");

    // Set projection
    this.setUniform("projectionMatrix");

    this._gl.disable(this._gl.DEPTH_TEST);
    this._gl.drawArrays(this._gl.LINES, 0, this.dataList["vertex"].length / 3);
  }
}
