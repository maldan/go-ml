varying highp vec4 vColor;

void main() {
    if (vColor.a <= 0.0) {
        discard;
    }

    gl_FragColor = vColor;
}