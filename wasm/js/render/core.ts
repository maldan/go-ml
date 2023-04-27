class GoRender {
  static _gl: WebGLRenderingContext;
  static _canvas: HTMLCanvasElement;

  static shaderSource: Record<string, string> = {};
  static layerList: GoRenderLayer[] = [];

  static setResolution(w: number, h: number) {
    this._canvas.setAttribute("width", `${w}`);
    this._canvas.setAttribute("height", `${h}`);
    this.onResize();
  }

  static scale(v: number) {
    const width = Number(this._canvas.getAttribute("width"));
    const height = Number(this._canvas.getAttribute("height"));

    this._canvas.style.width = width * v + "px";
    this._canvas.style.height = height * v + "px";
  }

  static async init() {
    this._canvas = document.querySelector("#glcanvas") as HTMLCanvasElement;
    if (!this._canvas) throw new Error(`Canvas not found`);

    this._gl = this._canvas.getContext("webgl", {
      antialias: false,
      premultipliedAlpha: false,
      // alpha: false,
    }) as WebGLRenderingContext;
    if (this._gl === null) throw new Error("WebGL is not supported");

    /*const extensions = this._gl.getSupportedExtensions();
    alert(extensions);

    if (!this._gl.getExtension("OES_texture_half_float")) {
      alert("OES_texture_half_float is not supported");
      throw new Error("OES_texture_half_float is not supported");
    }

    if (!this._gl.getExtension("OES_element_index_uint")) {
      alert("OES_element_index_uint is not supported");
      throw new Error("OES_element_index_uint is not supported");
    }*/

    // Load shaders
    const shaderList = [
      "matrix",
      "staticMesh",
      "dynamicMesh",
      "line",
      "point",
      "text",
      "ui",
    ];

    for (let i = 0; i < shaderList.length; i++) {
      const url = `./js/render/shader/${shaderList[i]}.glsl`;
      const x = await fetch(url);
      this.shaderSource[url] = await x.text();
    }

    // Inject library
    shaderList.slice(1).forEach((x) => {
      this.shaderSource[`./js/render/shader/${x}.glsl`] = this.shaderSource[
        `./js/render/shader/${x}.glsl`
      ].replace("// LIB", this.shaderSource["./js/render/shader/matrix.glsl"]);
    });

    // Separate shaders
    shaderList.slice(1).forEach((x) => {
      const shader = this.shaderSource[`./js/render/shader/${x}.glsl`];
      const tuple = shader.split("// Fragment");

      this.shaderSource[`./js/render/shader/${x}.vertex.glsl`] = tuple[0];
      this.shaderSource[`./js/render/shader/${x}.fragment.glsl`] = tuple[1];
    });

    // Main texture
    const texture = loadTexture(this._gl, "./texture.png");

    // Compile shaders
    this.layerList = [
      new GoRenderDynamicMeshLayer("dynamicMesh", this._gl),
      new GoRenderStaticMeshLayer("staticMesh", this._gl),
      new GoRenderPointLayer("point", this._gl),
      new GoRenderLineLayer("line", this._gl),
      new GoRenderTextLayer("text", this._gl),
      new GoRenderUILayer("ui", this._gl),
    ].map((x) => {
      x.init(
        this.shaderSource[`./js/render/shader/${x.name}.vertex.glsl`],
        this.shaderSource[`./js/render/shader/${x.name}.fragment.glsl`]
      );
      x.texture = texture;
      return x;
    });

    // @ts-ignore
    window.go.canvas = this._canvas;

    window.addEventListener("resize", () => {
      this.onResize();
    });
    this.onResize();
  }

  static onResize() {
    const width = Number(this._canvas.getAttribute("width"));
    const height = Number(this._canvas.getAttribute("height"));

    this._gl.viewport(0, 0, width, height);
    if ((window as any).go?.renderResize)
      (window as any).go.renderResize(width, height);
  }

  static draw() {
    // this._gl.colorMask(false, false, false, true);
    this._gl.clearColor(0.0, 0.0, 0.0, 1.0);
    this._gl.clearDepth(1.0);
    this._gl.clear(this._gl.COLOR_BUFFER_BIT | this._gl.DEPTH_BUFFER_BIT);

    // this._gl.enable(this._gl.CULL_FACE);
    // this._gl.cullFace(this._gl.FRONT);

    this._gl.enable(this._gl.BLEND);
    this._gl.blendFunc(this._gl.SRC_ALPHA, this._gl.ONE_MINUS_SRC_ALPHA);

    this._gl.enable(this._gl.DEPTH_TEST);
    this._gl.depthFunc(this._gl.LEQUAL);

    this.layerList.forEach((layer) => {
      layer.draw();
    });
  }
}

function loadTexture(gl: WebGLRenderingContext, url: string): WebGLTexture {
  const texture = gl.createTexture();
  if (!texture) throw new Error(`Can't create texture`);
  gl.bindTexture(gl.TEXTURE_2D, texture);
  gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);
  // gl.pixelStorei(gl.UNPACK_PREMULTIPLY_ALPHA_WEBGL, true);

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
    //

    gl.generateMipmap(gl.TEXTURE_2D);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
  };
  image.src = url;

  return texture;
}
