// LIB

// Attributes
attribute vec2 aVertex;
attribute vec2 aUv;

varying highp vec2 vUv;

void main() {
    // Set position
    gl_Position = vec4(aVertex.x, aVertex.y, 0.0, 1.0);
    vUv = aUv;
}

// Fragment
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif

varying vec2 vUv;
uniform sampler2D uTexture;

void main() {
    vec4 texelColor = texture2D(uTexture, vUv);
    gl_FragColor = vec4(texelColor.rgb, 1.0);
}