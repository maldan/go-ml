// LIB

// Attributes
attribute vec3 aVertex;
attribute vec3 aNormal;
attribute vec2 aUv;
attribute vec4 aColor;

varying highp vec2 vUv;
varying highp vec3 vLighting;
varying highp vec4 vColor;
varying highp vec4 vPosition;

uniform mat4 uProjectionMatrix;
uniform mat4 uLight;

void main() {
    // Set position
    gl_Position = uProjectionMatrix * vec4(aVertex, 1.0);
    vUv = aUv;
    vColor = aColor;
    vPosition = gl_Position;

    // Apply lighting effect
    highp vec3 directionalVector = normalize(vec3(uLight[0][0], uLight[0][1], uLight[0][2]));
    highp vec3 ambientLight = vec3(uLight[1][0], uLight[1][1], uLight[1][2]);
    highp vec3 directionalLightColor = vec3(uLight[2][0], uLight[2][1], uLight[2][2]);

    // Prepare normal matrix
    mat4 normalMatrix = identity();
    normalMatrix = inverse(normalMatrix);
    normalMatrix = transpose(normalMatrix);

    // Calculate light
    highp vec4 transformedNormal = normalMatrix * vec4(aNormal.xyz, 1.0);
    highp float directional = max(dot(transformedNormal.xyz, directionalVector), 0.0);
    vLighting = ambientLight + (directionalLightColor * directional);
}

// Fragment
precision highp float;

varying vec4 vColor;
varying vec2 vUv;
varying vec3 vLighting;
varying highp vec4 vPosition;

uniform sampler2D uTexture;

void main() {
    highp vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = vec4(texelColor.rgb * vLighting.rgb, texelColor.a) * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;
}