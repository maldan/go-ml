attribute vec3 aVertexPosition;
varying lowp vec3 vColor;
uniform mat4 uProjectionMatrix;

// LIB

void main() {
    gl_Position = uProjectionMatrix * vec4(aVertexPosition, 1.0);
    gl_PointSize = 3.0;
    vColor = gl_Position.xyz;
}