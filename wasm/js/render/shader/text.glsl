// LIB

// Attributes
attribute vec3 aVertex;
attribute vec2 aUv;
attribute vec4 aColor;

attribute vec3 aPosition;

varying vec2 vUv;
varying vec4 vColor;

uniform mat4 uProjectionMatrix;

void main() {
    // Set position
    gl_Position = uProjectionMatrix * translate(aPosition) * vec4(aVertex, 1.0);
    vUv = aUv;
    vColor = aColor;
}

// Fragment
precision highp float;

varying vec4 vColor;
varying vec2 vUv;

uniform sampler2D uTexture;

void main() {
    vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = vec4(texelColor.rgb, texelColor.a) * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;
}