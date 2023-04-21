// LIB

// Attributes
attribute vec3 aVertex;
attribute vec2 aUv;
//attribute vec4 aColor;
attribute vec3 aPosition;
attribute vec3 aScale;

varying highp vec2 vUv;
//varying highp vec4 vColor;

uniform mat4 uProjectionMatrix;

void main() {
    // Set position
    gl_Position = uProjectionMatrix * translate(aPosition) * scale(aScale) * vec4(aVertex, 1.0);
    vUv = aUv;
    //vColor = aColor;
}
