// LIB

attribute vec4 aVertex;
attribute vec4 aColor;

uniform mat4 uProjectionMatrix;

varying highp vec4 vColor;

void main() {
    gl_Position = uProjectionMatrix * vec4(aVertex.xyz, 1.0);
    gl_PointSize = aVertex.w;

    vColor = aColor;
}