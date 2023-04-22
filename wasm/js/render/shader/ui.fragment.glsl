precision highp float;

varying vec4 vColor;
varying vec2 vUv;

uniform sampler2D uTexture;

void main() {
    highp vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = vec4(texelColor.rgb, texelColor.a) * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;

    /*highp vec4 texelColor = texture2D(uTexture, vUv);
    gl_FragColor = vec4(texelColor.rgb, texelColor.a);*/
}