// LIB

// Attributes
attribute vec3 aVertex;
attribute vec2 aUv;
attribute vec4 aColor;
attribute vec3 aPosition;
attribute vec3 aRotation;
attribute vec3 aScale;

varying highp vec2 vUv;
varying highp vec4 vColor;

uniform mat4 uProjectionMatrix;

void main() {
    mat4 rotationMatrix = rotate(aRotation);

    // Set position
    vec3 newScale = aScale;
    newScale.y *= -1.0;
    gl_Position = uProjectionMatrix * translate(aPosition) * rotationMatrix * scale(newScale) * vec4(aVertex, 1.0);

    vUv = aUv;
    vColor = aColor;
}

// Fragment
precision highp float;

varying vec4 vColor;
varying vec2 vUv;

uniform sampler2D uTexture;

void main() {
    highp vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = texelColor * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;

    /*highp vec4 texelColor = texture2D(uTexture, vUv);
    gl_FragColor = vec4(texelColor.rgb, texelColor.a);*/
}