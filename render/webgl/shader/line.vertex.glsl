attribute vec3 aVertex;
uniform mat4 uProjectionMatrix;

// LIB

void main() {
    gl_Position = uProjectionMatrix * vec4(aVertex, 1.0);
}