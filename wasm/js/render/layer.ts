class GoRenderLayer {
  public _gl: WebGLRenderingContext;
  public shader: WebGLProgram;

  public attributeList: Record<string, any> = {};
  public uniformList: Record<string, any> = {};
  public bufferList: Record<string, any> = {};
  public dataList: Record<string, any> = {};
  public name = "";
  public state: Record<string, any> = {};
  public texture: WebGLTexture;

  constructor(name: string, gl: WebGLRenderingContext) {
    this.name = name;
    this._gl = gl;
  }

  init(vertex: string, fragment: string) {}

  loadShader(type: number, source: string) {
    const shader = this._gl.createShader(type);
    if (!shader) throw new Error(`Can't create shader`);

    this._gl.shaderSource(shader, source);
    this._gl.compileShader(shader);

    if (!this._gl.getShaderParameter(shader, this._gl.COMPILE_STATUS)) {
      const info = this._gl.getShaderInfoLog(shader);
      this._gl.deleteShader(shader);
      throw new Error(
        `An error occurred compiling the shaders: ${info}\n${source}`
      );
    }

    return shader;
  }

  compileShader(vertex: string, fragment: string): WebGLProgram {
    const vertexShader = this.loadShader(this._gl.VERTEX_SHADER, vertex);
    if (!vertexShader) throw new Error(`Vertex shader is null`);
    const fragmentShader = this.loadShader(this._gl.FRAGMENT_SHADER, fragment);
    if (!fragmentShader) throw new Error(`Fragment shader is null`);

    const shaderProgram = this._gl.createProgram();
    if (!shaderProgram) throw new Error(`Can't create shader program`);

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

  setWasmData(memory: ArrayBuffer, state: any) {}

  setDataArray(
    name: string,
    state: any,
    array: Uint16Array | Float32Array,
    length = 0
  ) {
    let offsetSize = 1;
    if (array instanceof Uint16Array) offsetSize = 2;
    if (array instanceof Float32Array) offsetSize = 4;

    this.dataList[name] = array.subarray(
      ~~state[name + "Pointer"] / offsetSize,
      ~~state[name + "Pointer"] / offsetSize +
        (length || ~~state[name + "Amount"])
    );
  }

  uploadData(type: string, name: string) {
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

  enableAttribute(name: string, size = 3) {
    this._gl.bindBuffer(this._gl.ARRAY_BUFFER, this.bufferList[name]);

    let attributeName = "a" + name[0].toUpperCase() + name.slice(1);

    this._gl.vertexAttribPointer(
      this.attributeList[attributeName],
      size,
      this._gl.FLOAT,
      false,
      0,
      0
    );
    this._gl.enableVertexAttribArray(this.attributeList[attributeName]);
  }

  setUniform(name: string) {
    let uniformName = "u" + name[0].toUpperCase() + name.slice(1);
    this._gl.uniformMatrix4fv(
      this.uniformList[uniformName],
      false,
      this.dataList[name]
    );
  }

  setTexture() {
    this._gl.activeTexture(this._gl.TEXTURE0);
    this._gl.bindTexture(this._gl.TEXTURE_2D, this.texture);
    this._gl.uniform1i(this.uniformList["uTexture"], 0);
  }

  draw() {}
}

class GoRenderPointLayer extends GoRenderLayer {
  init(vertex: string, fragment: string) {
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

  setWasmData(memory: ArrayBuffer, state: any) {
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

    this.enableAttribute("vertex", 4);
    this.enableAttribute("color", 4);

    // Set projection
    this.setUniform("projectionMatrix");

    // this._gl.disable(this._gl.DEPTH_TEST);
    this._gl.enable(this._gl.DEPTH_TEST);
    this._gl.drawArrays(this._gl.POINTS, 0, this.dataList["vertex"].length / 4);
  }
}

class GoRenderLineLayer extends GoRenderLayer {
  init(vertex: string, fragment: string) {
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

  setWasmData(memory: ArrayBuffer, state: any) {
    let float32Array = new Float32Array(memory);

    // Get pointers
    const mv = (window as any).go.memoryView;
    ["vertex", "color"].forEach((x) => {
      state[x + "Pointer"] = mv.getUint32(
        (window as any).go.pointer[`renderLineLayer_${x}`],
        true
      );
      state[x + "Amount"] = mv.getUint32(
        (window as any).go.pointer[`renderLineLayer_${x}`] + 8,
        true
      );
    });

    // Get camera matrix
    state["projectionMatrixPointer"] = (window as any).go.pointer[
      `renderCamera_matrix`
    ];

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
    this.enableAttribute("color", 4);

    // Set projection
    this.setUniform("projectionMatrix");

    this._gl.enable(this._gl.DEPTH_TEST);
    // this._gl.disable(this._gl.DEPTH_TEST);
    this._gl.drawArrays(this._gl.LINES, 0, this.dataList["vertex"].length / 3);
  }
}

class GoRenderTextLayer extends GoRenderLayer {
  init(vertex: string, fragment: string) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex", "index", "position", "uv", "color"].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    ["aVertex", "aPosition", "aUv", "aColor"].forEach((x) => {
      this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
    });
    ["uProjectionMatrix", "uTexture"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });
  }

  setWasmData(memory: ArrayBuffer, state: any) {
    let shortArray = new Uint16Array(memory);
    let float32Array = new Float32Array(memory);

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("uv", state, float32Array);
    this.setDataArray("position", state, float32Array);
    this.setDataArray("color", state, float32Array);

    this.setDataArray("index", state, shortArray);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  draw() {
    // Set program
    this._gl.useProgram(this.shader);

    // Put main data
    this.uploadData("element", "index");
    this.uploadData("any", "vertex");
    this.uploadData("any", "uv");
    this.uploadData("any", "position");
    this.uploadData("any", "color");

    // Enable attributes
    this.enableAttribute("vertex");
    this.enableAttribute("uv", 2);
    this.enableAttribute("position");
    this.enableAttribute("color");

    // Set projection
    this.setUniform("projectionMatrix");

    // Set texture
    this.setTexture();

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

class GoRenderUILayer extends GoRenderLayer {
  init(vertex: string, fragment: string) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex", "index", "position", "rotation", "scale", "uv", "color"].forEach(
      (x) => {
        this.bufferList[x] = this._gl.createBuffer();
      }
    );
    ["aVertex", "aPosition", "aRotation", "aScale", "aUv", "aColor"].forEach(
      (x) => {
        this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
      }
    );
    ["uProjectionMatrix", "uTexture"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });
  }

  setWasmData(memory: ArrayBuffer, state: any) {
    let shortArray = new Uint16Array(memory);
    let float32Array = new Float32Array(memory);

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("uv", state, float32Array);
    this.setDataArray("position", state, float32Array);
    this.setDataArray("rotation", state, float32Array);
    this.setDataArray("scale", state, float32Array);
    this.setDataArray("color", state, float32Array);

    this.setDataArray("index", state, shortArray);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  draw() {
    // Set program
    this._gl.useProgram(this.shader);

    // Put main data
    this.uploadData("element", "index");
    this.uploadData("any", "vertex");
    this.uploadData("any", "uv");
    this.uploadData("any", "position");
    this.uploadData("any", "rotation");
    this.uploadData("any", "scale");
    this.uploadData("any", "color");

    // Enable attributes
    this.enableAttribute("vertex");
    this.enableAttribute("uv", 2);

    this.enableAttribute("position");
    this.enableAttribute("rotation");
    this.enableAttribute("scale");
    this.enableAttribute("color", 4);

    // Set projection
    this.setUniform("projectionMatrix");

    // Set texture
    this.setTexture();

    this._gl.bindBuffer(
      this._gl.ELEMENT_ARRAY_BUFFER,
      this.bufferList["index"]
    );

    this._gl.disable(this._gl.DEPTH_TEST);
    this._gl.drawElements(
      this._gl.TRIANGLES,
      this.dataList["index"].length,
      this._gl.UNSIGNED_SHORT,
      0
    );
  }
}

class GoRenderDynamicMeshLayer extends GoRenderLayer {
  init(vertex: string, fragment: string) {
    this.shader = this.compileShader(vertex, fragment);

    [
      "vertex",
      "index",
      "position",
      "rotation",
      "uv",
      "normal",
      "scale",
      "color",
    ].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    [
      "aVertex",
      "aPosition",
      "aRotation",
      "aScale",
      "aUv",
      "aNormal",
      "aColor",
    ].forEach((x) => {
      this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
    });
    ["uProjectionMatrix", "uTexture", "uLight"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });

    this.dataList["light"] = new Float32Array([
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ]);
  }

  setWasmData(memory: ArrayBuffer, state: any) {
    let byteArray = new Uint8Array(memory);
    let shortArray = new Uint16Array(memory);
    let float32Array = new Float32Array(memory);

    // Get pointers
    const mv = (window as any).go.memoryView;
    [
      "vertex",
      "uv",
      "normal",
      "position",
      "rotation",
      "scale",
      "color",
      "index",
    ].forEach((x) => {
      state[x + "Pointer"] = mv.getUint32(
        (window as any).go.pointer[`renderDynamicMeshLayer_${x}`],
        true
      );
      state[x + "Amount"] = mv.getUint32(
        (window as any).go.pointer[`renderDynamicMeshLayer_${x}`] + 8,
        true
      );
    });

    // Get camera matrix
    state["projectionMatrixPointer"] = (window as any).go.pointer[
      `renderCamera_matrix`
    ];

    // Set light
    const w = window as any;
    const lightPtr = w.go.pointer[`renderLight`];
    for (let i = 0; i < 3; i++) {
      this.dataList["light"][i] = w.go.memory.readF32(lightPtr + i * 4);
      this.dataList["light"][i + 4] = w.go.memory.readF32(
        lightPtr + 12 + i * 4
      );
      this.dataList["light"][i + 8] = w.go.memory.readF32(
        lightPtr + 24 + i * 4
      );
    }

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("normal", state, float32Array);
    this.setDataArray("uv", state, float32Array);
    this.setDataArray("position", state, float32Array);
    this.setDataArray("rotation", state, float32Array);
    this.setDataArray("scale", state, float32Array);
    this.setDataArray("color", state, float32Array);
    this.setDataArray("index", state, shortArray);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  draw() {
    // Set program
    this._gl.useProgram(this.shader);

    // Put main data
    this.uploadData("element", "index");
    this.uploadData("any", "vertex");
    this.uploadData("any", "normal");
    this.uploadData("any", "uv");
    this.uploadData("any", "position");
    this.uploadData("any", "rotation");
    this.uploadData("any", "scale");
    this.uploadData("any", "color");

    // Enable attributes
    this.enableAttribute("vertex");
    this.enableAttribute("normal");
    this.enableAttribute("uv", 2);
    this.enableAttribute("position");
    this.enableAttribute("rotation");
    this.enableAttribute("scale");
    this.enableAttribute("color", 4);

    // Set projection
    this.setUniform("projectionMatrix");
    this.setUniform("light");

    // Set texture
    this.setTexture();

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

class GoRenderStaticMeshLayer extends GoRenderLayer {
  init(vertex: string, fragment: string) {
    this.shader = this.compileShader(vertex, fragment);

    ["vertex", "index", "uv", "normal", "color"].forEach((x) => {
      this.bufferList[x] = this._gl.createBuffer();
    });
    ["aVertex", "aUv", "aNormal", "aColor"].forEach((x) => {
      this.attributeList[x] = this._gl.getAttribLocation(this.shader, x);
    });
    ["uProjectionMatrix", "uTexture"].forEach((x) => {
      this.uniformList[x] = this._gl.getUniformLocation(this.shader, x);
    });
  }

  setWasmData(memory: ArrayBuffer, state: any) {
    let byteArray = new Uint8Array(memory);
    let shortArray = new Uint16Array(memory);
    let float32Array = new Float32Array(memory);

    // Get pointers
    const mv = (window as any).go.memoryView;
    ["vertex", "uv", "normal", "color", "index"].forEach((x) => {
      state[x + "Pointer"] = mv.getUint32(
        (window as any).go.pointer[`renderStaticMeshLayer_${x}`],
        true
      );
      state[x + "Amount"] = mv.getUint32(
        (window as any).go.pointer[`renderStaticMeshLayer_${x}`] + 8,
        true
      );
    });

    // Get camera matrix
    state["projectionMatrixPointer"] = (window as any).go.pointer[
      `renderCamera_matrix`
    ];

    this.setDataArray("vertex", state, float32Array);
    this.setDataArray("uv", state, float32Array);
    this.setDataArray("normal", state, float32Array);
    this.setDataArray("color", state, float32Array);
    this.setDataArray("index", state, shortArray);
    this.setDataArray("projectionMatrix", state, float32Array, 16);
  }

  draw() {
    // Set program
    this._gl.useProgram(this.shader);

    // Put main data
    this.uploadData("element", "index");
    this.uploadData("any", "vertex");
    this.uploadData("any", "normal");
    this.uploadData("any", "uv");
    this.uploadData("any", "color");

    // Enable attributes
    this.enableAttribute("vertex");
    this.enableAttribute("normal");
    this.enableAttribute("uv", 2);
    this.enableAttribute("color", 4);

    // Set projection
    this.setUniform("projectionMatrix");

    // Set texture
    this.setTexture();

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
