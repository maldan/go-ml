class GoRender {
  static _gl: WebGLRenderingContext;
  static _canvas: HTMLCanvasElement;

  static shaderSource: Record<string, string> = {};
  static layerList: GoRenderLayer[] = [];
  static renderTexture: WebGLTexture;
  static framebuffer: WebGLFramebuffer;
  static depthTexture: WebGLTexture;

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
    alert(extensions);*/

    /*if (!this._gl.getExtension("WEBGL_depth_texture")) {
      alert("WEBGL_depth_texture is not supported");
      throw new Error("WEBGL_depth_texture is not supported");
    }*/

    /*if (!this._gl.getExtension("OES_texture_half_float")) {
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
      "postprocessing",
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
      new GoRenderStaticMeshLayer("staticMesh", this._gl),
      new GoRenderDynamicMeshLayer("dynamicMesh", this._gl),
      new GoRenderPointLayer("point", this._gl),
      new GoRenderLineLayer("line", this._gl),
      // new GoRenderPostProcessingLayer("postprocessing", this._gl),
      new GoRenderUILayer("ui", this._gl),

      /*
        new GoRenderTextLayer("text", this._gl),
      */
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

    // Frame buffer
    /*this.renderTexture = createTexture(this._gl, width, height);
    this.depthTexture = createDepthTexture(this._gl, width, height);

    this.framebuffer = this._gl.createFramebuffer() as WebGLFramebuffer;
    this._gl.bindFramebuffer(this._gl.FRAMEBUFFER, this.framebuffer);
    this._gl.framebufferTexture2D(
      this._gl.FRAMEBUFFER,
      this._gl.COLOR_ATTACHMENT0,
      this._gl.TEXTURE_2D,
      this.renderTexture,
      0
    );
    this._gl.framebufferTexture2D(
      this._gl.FRAMEBUFFER,
      this._gl.DEPTH_ATTACHMENT,
      this._gl.TEXTURE_2D,
      this.depthTexture,
      0
    );*/

    // render to the canvas
    this._gl.bindFramebuffer(this._gl.FRAMEBUFFER, null);
  }

  static draw() {
    this._gl.clearColor(0.0, 0.0, 0.0, 1.0);
    this._gl.clearDepth(1.0);
    this._gl.clear(this._gl.COLOR_BUFFER_BIT | this._gl.DEPTH_BUFFER_BIT);

    this._gl.enable(this._gl.CULL_FACE);
    this._gl.cullFace(this._gl.BACK);

    this._gl.enable(this._gl.BLEND);
    this._gl.blendFunc(this._gl.SRC_ALPHA, this._gl.ONE_MINUS_SRC_ALPHA);

    this._gl.enable(this._gl.DEPTH_TEST);
    this._gl.depthFunc(this._gl.LEQUAL);

    this.layerList.forEach((layer) => {
      layer.draw();
    });

    // Render to texture
    /*this._gl.bindFramebuffer(this._gl.FRAMEBUFFER, this.framebuffer);
    this._gl.framebufferTexture2D(
      this._gl.FRAMEBUFFER,
      this._gl.DEPTH_ATTACHMENT,
      this._gl.TEXTURE_2D,
      this.depthTexture,
      0
    );*/

    /*// Render layers
    this.layerList
      .filter((x) => x.isRenderToTexture)
      .forEach((layer) => {
        layer.draw();
      });

    // Render to canvas
    this._gl.bindFramebuffer(this._gl.FRAMEBUFFER, null);

    this._gl.disable(this._gl.BLEND);
    this._gl.disable(this._gl.DEPTH_TEST);

    // Render postprocessing
    this.layerList
      .filter((x) => x.name == "postprocessing")
      .forEach((layer) => {
        layer.texture = this.renderTexture;
        layer.draw();
      });

    // Render other layers
    this.layerList
      .filter((x) => !x.isRenderToTexture && x.name != "postprocessing")
      .forEach((layer) => {
        layer.draw();
      });*/
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

    gl.texParameteri(
      gl.TEXTURE_2D,
      gl.TEXTURE_MIN_FILTER,
      gl.LINEAR_MIPMAP_LINEAR
    );
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
  };
  image.src = url;

  return texture;
}

function createTexture(
  gl: WebGLRenderingContext,
  width: number,
  height: number
): WebGLTexture {
  const texture = gl.createTexture();
  if (!texture) throw new Error(`Can't create texture`);
  gl.bindTexture(gl.TEXTURE_2D, texture);

  const level = 0;
  const internalFormat = gl.RGBA;
  const border = 0;
  const srcFormat = gl.RGBA;
  const srcType = gl.UNSIGNED_BYTE;

  gl.texImage2D(
    gl.TEXTURE_2D,
    level,
    internalFormat,
    width,
    height,
    border,
    srcFormat,
    srcType,
    null
  );

  // set the filtering so we don't need mips
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

  return texture;
}

function createDepthTexture(
  gl: WebGLRenderingContext,
  width: number,
  height: number
): WebGLTexture {
  const texture = gl.createTexture();
  if (!texture) throw new Error(`Can't create texture`);
  gl.bindTexture(gl.TEXTURE_2D, texture);

  gl.texImage2D(
    gl.TEXTURE_2D,
    0,
    gl.DEPTH_COMPONENT,
    gl.drawingBufferWidth,
    gl.drawingBufferHeight,
    0,
    gl.DEPTH_COMPONENT,
    gl.UNSIGNED_SHORT,
    null
  );

  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

  return texture;
}
