attribute vec3 aVertex;
attribute vec4 aColor;

uniform mat4 uPerspectiveCamera;
uniform mat4 uUiCamera;

varying vec4 vColor;

void main() {
    vColor = aColor;

    // Set position
    if (aVertex.z >= 1000000.0) {
        gl_Position = uUiCamera * vec4(aVertex.xy, 0.0, 1.0);
    } else {
        gl_Position = uPerspectiveCamera * vec4(aVertex, 1.0);
    }
}

// Fragment
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif

varying vec4 vColor;

void main() {
    if (vColor.a <= 0.0) {
        discard;
    }
    gl_FragColor = vColor;
}