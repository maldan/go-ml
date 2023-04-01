varying lowp vec4 vColor;
varying highp vec2 vUv;
varying highp vec3 vLighting;

uniform sampler2D uTexture;

void main() {
    highp vec4 texelColor = texture2D(uTexture, vUv);
    gl_FragColor = vec4(texelColor.rgb * vLighting, texelColor.a);
}