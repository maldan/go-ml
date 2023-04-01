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
      return x;
    });

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

    this.layerList.forEach((layer) => {
      layer.draw();
    });
  }
}

class GoRenderWasm {
  static wasmTime = [];
  static jsTime = [];

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
    go.run(wasmModule.instance);
    console.log(wasmModule.instance);

    // Init webgl render
    await GoRender.init();

    let start;

    const step = (timestamp) => {
      if (start === undefined) start = timestamp;
      const elapsed = timestamp - start;

      this.beforeFrame(elapsed);
      goWasmGameTick();

      // Calculate scene in golang
      let pp = performance.now();
      goWasmRenderFrame();
      this.wasmTime.push(performance.now() - pp);

      // Get state
      let state = goWasmRenderState();

      // Send golang data to webgl render
      GoRender.layerList.forEach((layer) => {
        layer.state = state;

        if (wasmModule.instance.exports.mem) {
          layer.setWasmData(
            wasmModule.instance.exports.mem.buffer,
            state[layer.name + "Layer"]
          );
        } else {
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

      this.afterFrame(elapsed);

      window.requestAnimationFrame(step);
    };

    window.requestAnimationFrame(step);

    // Timers
    const avg = (x) => x.reduce((a, b) => a + b) / x.length;
    setInterval(() => {
      document.getElementById("wasm").innerHTML = `${avg(this.wasmTime).toFixed(
        2
      )}`;
      document.getElementById("js").innerHTML = `${avg(this.jsTime).toFixed(
        2
      )}`;
      this.wasmTime.length = 0;
      this.jsTime.length = 0;
    }, 1000);
  }
}
