// LIB

// Attributes
attribute vec3 aVertex;
attribute vec3 aNormal;
attribute vec2 aUv;
attribute vec4 aColor;

attribute vec3 aPosition;
attribute vec3 aRotation;
attribute vec3 aScale;

varying highp vec2 vUv;
varying highp vec3 vLighting;
varying highp vec4 vColor;

uniform mat4 uProjectionMatrix;

void main() {
    mat4 rotationMatrix = rotate(aRotation);

    // Set position
    vec3 position = aPosition;
    gl_Position = uProjectionMatrix * translate(position) * rotationMatrix * scale(aScale) * vec4(aVertex, 1.0);
    vUv = aUv;
    vColor = aColor;

    // Apply lighting effect
    highp vec3 ambientLight = vec3(0.3, 0.3, 0.3);
    highp vec3 directionalLightColor = vec3(1, 1, 1);
    highp vec3 directionalVector = normalize(vec3(0.3, 0.4, 0.8));

    // Prepare normal matrix
    mat4 normalMatrix = identity();
    normalMatrix = inverse(rotationMatrix);

    // Calculate light
    highp vec4 transformedNormal = normalMatrix * vec4(aNormal.xyz, 1.0);
    highp float directional = max(dot(transformedNormal.xyz, directionalVector), 0.0);
    vLighting = ambientLight + (directionalLightColor * directional);
}
