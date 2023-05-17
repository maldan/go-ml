// LIB

// Attributes
attribute vec3 aVertex;
// attribute vec2 aUv;
attribute vec4 aColor;

attribute vec3 aPosition;
attribute vec3 aRotation;

// varying vec2 vUv;
varying vec4 vColor;

uniform mat4 uPerspectiveCamera;
uniform mat4 uUiCamera;

void main() {
    mat4 modelViewMatrix = translate(aPosition) * rotate(aRotation);

    // Set position
    if (aVertex.z > 0.0) {
        gl_Position = uUiCamera * modelViewMatrix * vec4(aVertex.x, aVertex.y*-1.0, 0.0, 1.0);
    } else {
        gl_Position = uPerspectiveCamera * modelViewMatrix * vec4(aVertex.xy, 0.0, 1.0);
    }

    vColor = aColor;
}

// Fragment
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif

varying vec4 vColor;
// varying vec2 vUv;

// uniform sampler2D uTexture;

void main() {
    /*vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = vec4(texelColor.rgb, texelColor.a) * vColor.rgba;
    if (finalColor.a > 0.2) {
        finalColor.rgb = vColor.rgb;
    }
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;*/
    if (vColor.a <= 0.0) {
        discard;
    }
    gl_FragColor = vColor;
}