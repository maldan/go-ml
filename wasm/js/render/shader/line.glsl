attribute vec3 aVertex;
attribute vec3 aColor;
uniform mat4 uProjectionMatrix;

varying lowp vec3 vColor;

// LIB

void main() {
    vColor = aColor;
    gl_Position = uProjectionMatrix * vec4(aVertex, 1.0);
}

// Fragment
varying lowp vec3 vColor;

void main() {
    gl_FragColor = vec4(vColor, 1.0);
}