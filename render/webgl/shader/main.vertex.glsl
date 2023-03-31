// Attributes
attribute vec3 aVertex;
attribute vec4 aPosition;

varying lowp vec3 vColor;
uniform mat4 uProjectionMatrix;

// LIB

void main() {
    mat4 modelMatrix = mat4(
        vec4(1.0, 0.0, 0.0, 0.0),
        vec4(0.0, 1.0, 0.0, 0.0),
        vec4(0.0, 0.0, 1.0, 0.0),
        vec4(0.0, 0.0, 0.0, 1.0)
    );
    modelMatrix = translate(modelMatrix, aPosition);

    // Set position
    gl_Position = uProjectionMatrix * modelMatrix * vec4(aVertex, 1.0);
    vColor = gl_Position.xyz;
}