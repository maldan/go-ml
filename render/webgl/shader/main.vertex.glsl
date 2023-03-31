attribute vec3 aVertexPosition;
attribute vec4 aModelPosition;
attribute vec3 aVertexColor;

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
    modelMatrix = translate(modelMatrix, aModelPosition);

    gl_Position = uProjectionMatrix * modelMatrix * vec4(aVertexPosition, 1.0);
    vColor = gl_Position.xyz;
}