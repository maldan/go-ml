attribute vec3 aVertex;
attribute vec4 aColor;
uniform mat4 uProjectionMatrix;

varying lowp vec4 vColor;

// LIB

void main() {
    vColor = aColor;
    gl_Position = uProjectionMatrix * vec4(aVertex, 1.0);
}

// Fragment
varying lowp vec4 vColor;

void main() {
    if (vColor.a <= 0.0) {
        discard;
    }
    gl_FragColor = vColor;
}