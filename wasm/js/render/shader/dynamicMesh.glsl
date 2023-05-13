// LIB

// Attributes
attribute vec3 aVertex;
attribute vec3 aNormal;
attribute vec2 aUv;
attribute vec4 aColor;

attribute vec3 aPosition;
attribute vec3 aRotation;
attribute vec3 aScale;

varying vec2 vUv;
varying vec3 vLighting;
varying vec4 vColor;

uniform mat4 uProjectionMatrix;
uniform mat4 uLight;

void main() {
    mat4 modelViewMatrix = translate(aPosition) * rotate(aRotation) * scale(aScale);

    // Set position
    gl_Position = uProjectionMatrix * modelViewMatrix * vec4(aVertex, 1.0);
    vUv = aUv;
    vColor = aColor;

    // Apply lighting effect
    vec3 directionalVector = normalize(vec3(uLight[0][0], uLight[0][1], uLight[0][2]));
    vec3 ambientLight = vec3(uLight[1][0], uLight[1][1], uLight[1][2]);
    vec3 directionalLightColor = vec3(uLight[2][0], uLight[2][1], uLight[2][2]);

    // Prepare normal matrix
    mat4 normalMatrix = identity();
    normalMatrix = inverse(modelViewMatrix);
    normalMatrix = transpose(normalMatrix);

    // Calculate light
    vec4 transformedNormal = normalMatrix * vec4(aNormal.xyz, 1.0);
    float directional = max(dot(transformedNormal.xyz, directionalVector), 0.0);
    vLighting = ambientLight + (directionalLightColor * directional);
}

// Fragment
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif

varying vec4 vColor;
varying vec2 vUv;
varying vec3 vLighting;

uniform sampler2D uTexture;

void main() {
    vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = vec4(texelColor.rgb * vLighting.rgb, texelColor.a) * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;
}