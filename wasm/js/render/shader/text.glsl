// LIB

// Attributes
attribute vec3 aVertex;
attribute vec2 aUv;
attribute vec4 aColor;

attribute vec3 aPosition;
attribute vec3 aRotation;

varying vec2 vUv;
varying vec4 vColor;

uniform mat4 uProjectionMatrix;

void main() {
    mat4 modelViewMatrix = translate(aPosition) * rotate(aRotation) * scale(aScale);

    // Set position
    gl_Position = uProjectionMatrix * modelViewMatrix * vec4(aVertex, 1.0);
    vUv = aUv;
    vColor = aColor;
}

// Fragment
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif

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