varying highp vec4 vColor;
varying highp vec2 vUv;
varying highp vec3 vLighting;

uniform sampler2D uTexture;

void main() {
    highp vec4 texelColor = texture2D(uTexture, vUv);
    if (texelColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = vec4(texelColor.rgb * vLighting, texelColor.a);
}