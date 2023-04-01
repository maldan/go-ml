// LIB

// Attributes
attribute vec4 aVertex;
attribute vec4 aPosition;
attribute vec4 aRotation;

varying lowp vec4 vColor;
uniform mat4 uProjectionMatrix;

void main() {
    mat4 modelMatrix = mat4(
        1.0, 0.0, 0.0, 0.0,
        0.0, 1.0, 0.0, 0.0,
        0.0, 0.0, 1.0, 0.0,
        0.0, 0.0, 0.0, 1.0
    );

    modelMatrix = translate2(modelMatrix, aPosition);

    /*modelMatrix = rotateZ(modelMatrix, aRotation.z);
    modelMatrix = rotateY(modelMatrix, aRotation.y);
    modelMatrix = rotateX(modelMatrix, aRotation.x);*/
    // modelMatrix = rotateZ(modelMatrix, aRotation.z);

    // Set position
    gl_Position = uProjectionMatrix * modelMatrix * (aVertex * rotateX(aRotation.x) * rotateY(aRotation.y) * rotateZ(aRotation.z));
    // gl_Position = uProjectionMatrix * modelMatrix * (aVertex);

    vColor = gl_Position;
}
