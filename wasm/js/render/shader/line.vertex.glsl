attribute vec3 aVertex;
attribute vec3 aColor;
uniform mat4 uProjectionMatrix;

varying lowp vec3 vColor;

// LIB

void main() {
    vColor = aColor;
    gl_Position = uProjectionMatrix * vec4(aVertex, 1.0);
}