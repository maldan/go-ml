precision highp float;

varying vec4 vColor;
varying vec2 vUv;
varying vec3 vLighting;

uniform sampler2D uTexture;

void main() {
    highp vec4 texelColor = texture2D(uTexture, vUv);
    vec4 finalColor = vec4(texelColor.rgb * vLighting.rgb, texelColor.a) * vColor.rgba;
    if (finalColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = finalColor;
}