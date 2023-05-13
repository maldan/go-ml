// LIB

// Attributes
attribute vec2 aVertex;
attribute vec2 aUv;
attribute vec4 aColor;
attribute vec2 aPosition;
attribute float aRotation;
attribute vec2 aScale;

varying vec2 vUv;
varying vec4 vColor;

uniform mat4 uProjectionMatrix;

void main() {
    // Set position
    vec3 newScale = vec3(aScale, 1.0);
    newScale.y *= -1.0;
    gl_Position = uProjectionMatrix * translate(vec3(aPosition, 0.0)) * rotate(vec3(0.0, 0.0, aRotation)) * scale(newScale) * vec4(aVertex.x, aVertex.y, 0.0, 1.0);

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
    vec4 finalColor = texelColor * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;
}