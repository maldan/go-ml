class GoRender {
  static _gl: WebGLRenderingContext;
  static _canvas: HTMLCanvasElement;

  static shaderSource: Record<string, string> = {};
  static layerList: GoRenderLayer[] = [];

  static downscale = 4;

  static async loadShaderSourceCode(url: string) {
    const x = await fetch(url);
    this.shaderSource[url] = await x.text();
  }

  static async init() {
    this._canvas = document.querySelector("#glcanvas") as HTMLCanvasElement;
    if (!this._canvas) throw new Error(`Canvas not found`);

    this._gl = this._canvas.getContext("webgl", {
      antialias: false,
    }) as WebGLRenderingContext;
    if (this._gl === null) throw new Error("WebGL is not supported");

    // Load shaders
    const shaderList = ["matrix.glsl"];
    shaderList.push(
      ...["main", "point", "line", "text", "ui"]
        .map((x) => [`${x}.vertex.glsl`, `${x}.fragment.glsl`])
        .flat()
    );
    for (let i = 0; i < shaderList.length; i++) {
      await this.loadShaderSourceCode(`./js/render/shader/${shaderList[i]}`);
    }

    // Inject library
    shaderList.slice(1).forEach((x) => {
      this.shaderSource[`./js/render/shader/${x}`] = this.shaderSource[
        `./js/render/shader/${x}`
      ].replace("// LIB", this.shaderSource["./js/render/shader/matrix.glsl"]);
    });

    // Main texture
    const texture = loadTexture(this._gl, "./texture.png");

    // Compile shaders
    this.layerList = [
      new GoRenderLayer("main", this._gl),
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

    window.addEventListener("resize", () => {
      this.onResize();
    });
    this.onResize();
  }

  static onResize() {
    this._canvas.setAttribute(
      "width",
      window.innerWidth / this.downscale + "px"
    );
    this._canvas.setAttribute(
      "height",
      window.innerHeight / this.downscale + "px"
    );
    this._gl.viewport(
      0,
      0,
      window.innerWidth / this.downscale,
      window.innerHeight / this.downscale
    );

    if (window.go)
      window.go.renderResize(
        window.innerWidth / this.downscale,
        window.innerHeight / this.downscale
      );
  }

  static draw() {
    this._gl.clearColor(0.0, 0.0, 0.0, 1.0);
    this._gl.clearDepth(1.0);
    this._gl.enable(this._gl.DEPTH_TEST);

    // this._gl.enable(this._gl.CULL_FACE);
    // this._gl.cullFace(this._gl.BACK);

    this._gl.depthFunc(this._gl.LEQUAL);
    this._gl.clear(this._gl.COLOR_BUFFER_BIT | this._gl.DEPTH_BUFFER_BIT);

    this.layerList.forEach((layer) => {
      layer.draw();
    });
  }
}

function loadTexture(gl: WebGLRenderingContext, url: string) {
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
    gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);
    gl.texImage2D(
      gl.TEXTURE_2D,
      level,
      internalFormat,
      srcFormat,
      srcType,
      image
    );
    //

    // gl.generateMipmap(gl.TEXTURE_2D);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
  };
  image.src = url;

  return texture;
}
