// Vertex shader program
const vsSource = `
    attribute vec3 aVertexPosition;
    attribute vec3 aVertexColor;
    
    varying lowp vec3 vColor;
      
    void main() {
      gl_Position = vec4(aVertexPosition, 1.0);
      vColor = gl_Position.xyz;
    }
`;

const fsSource = `
    varying lowp vec3 vColor;
    
    void main() {
      gl_FragColor = vec4(vColor, 1.0);
    }
`;

function initShaderProgram(gl, vsSource, fsSource) {
  const vertexShader = loadShader(gl, gl.VERTEX_SHADER, vsSource);
  const fragmentShader = loadShader(gl, gl.FRAGMENT_SHADER, fsSource);

  const shaderProgram = gl.createProgram();
  gl.attachShader(shaderProgram, vertexShader);
  gl.attachShader(shaderProgram, fragmentShader);
  gl.linkProgram(shaderProgram);

  if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS)) {
    console.log(
      `Unable to initialize the shader program: ${gl.getProgramInfoLog(
        shaderProgram
      )}`
    );
    return null;
  }

  return shaderProgram;
}

function loadShader(gl, type, source) {
  const shader = gl.createShader(type);
  gl.shaderSource(shader, source);
  gl.compileShader(shader);
  if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
    console.log(
      `An error occurred compiling the shaders: ${gl.getShaderInfoLog(shader)}`
    );
    gl.deleteShader(shader);
    return null;
  }

  return shader;
}

function initColorBuffer(gl) {
  var colors = [
    [1.0, 1.0, 1.0, 1.0], // Front face: white
    [1.0, 0.0, 0.0, 1.0], // Back face: red
    [0.0, 1.0, 0.0, 1.0], // Top face: green
    [0.0, 0.0, 1.0, 1.0], // Bottom face: blue
    [1.0, 1.0, 0.0, 1.0], // Right face: yellow
    [1.0, 0.0, 1.0, 1.0], // Left face: purple

    [1.0, 1.0, 1.0, 1.0], // Front face: white
    [1.0, 0.0, 0.0, 1.0], // Back face: red
    [0.0, 1.0, 0.0, 1.0], // Top face: green
    [0.0, 0.0, 1.0, 1.0], // Bottom face: blue
    [1.0, 1.0, 0.0, 1.0], // Right face: yellow
    [1.0, 0.0, 1.0, 1.0], // Left face: purple
  ];

  var generatedColors = [];

  for (j = 0; j < 6; j++) {
    var c = colors[j];

    for (var i = 0; i < 4; i++) {
      generatedColors = generatedColors.concat(c);
    }
  }

  var cubeVerticesColorBuffer = gl.createBuffer();
  gl.bindBuffer(gl.ARRAY_BUFFER, cubeVerticesColorBuffer);
  gl.bufferData(
    gl.ARRAY_BUFFER,
    new Float32Array(generatedColors),
    gl.STATIC_DRAW
  );
  return cubeVerticesColorBuffer;
}

function initPositionBuffer(gl) {
  const buffer = gl.createBuffer();
  gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
  gl.bufferData(gl.ARRAY_BUFFER, GLOBAL_VERTEX_ARRAY, gl.DYNAMIC_DRAW);
  return buffer;
}

function initIndexBuffer(gl) {
  const buffer = gl.createBuffer();
  gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer);
  gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, GLOBAL_INDEX_ARRAY, gl.DYNAMIC_DRAW);
  return buffer;
}

function drawScene(gl, programInfo, buffer, indexBuffer, colorBuffer) {
  gl.clearColor(0.0, 0.0, 0.0, 1.0);
  gl.clearDepth(1.0);
  gl.enable(gl.DEPTH_TEST);
  gl.depthFunc(gl.LEQUAL);
  gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

  // Put data
  gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);
  gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, GLOBAL_INDEX_ARRAY, gl.DYNAMIC_DRAW);
  gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
  gl.bufferData(gl.ARRAY_BUFFER, GLOBAL_VERTEX_ARRAY, gl.DYNAMIC_DRAW);

  // Enable attributes
  gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
  gl.vertexAttribPointer(
    programInfo.attribLocations.vertexPosition,
    3,
    gl.FLOAT,
    false,
    0,
    0
  );
  gl.enableVertexAttribArray(programInfo.attribLocations.vertexPosition);

  gl.bindBuffer(gl.ARRAY_BUFFER, colorBuffer);
  gl.vertexAttribPointer(
    programInfo.attribLocations.colorPosition,
    3,
    gl.FLOAT,
    false,
    0,
    0
  );
  gl.enableVertexAttribArray(programInfo.attribLocations.colorPosition);

  /*vertexColorAttribute = gl.getAttribLocation(shaderProgram, "aVertexColor");
  gl.enableVertexAttribArray(vertexColorAttribute);*/

  gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);
  gl.useProgram(programInfo.program);
  gl.drawElements(
    gl.TRIANGLES,
    GLOBAL_INDEX_ARRAY.length,
    gl.UNSIGNED_SHORT,
    0
  );
  //gl.drawElements(gl.TRIANGLE_STRIP, 0, GLOBAL_INDEX_ARRAY.length / 6);
}

function main() {
  const canvas = document.querySelector("#glcanvas");
  const gl = canvas.getContext("webgl");

  // Only continue if WebGL is available and working
  if (gl === null) {
    console.log(
      "Unable to initialize WebGL. Your browser or machine may not support it."
    );
    return;
  }

  const shaderProgram = initShaderProgram(gl, vsSource, fsSource);
  const programInfo = {
    program: shaderProgram,
    attribLocations: {
      vertexPosition: gl.getAttribLocation(shaderProgram, "aVertexPosition"),
      colorPosition: gl.getAttribLocation(shaderProgram, "aVertexColor"),
    },
  };

  const buffers = initPositionBuffer(gl);
  const buffers2 = initIndexBuffer(gl);
  const buffers3 = initColorBuffer(gl);

  this.realDraw = () => {
    drawScene(gl, programInfo, buffers, buffers2, buffers3);
  };
}

function realDraw() {}
