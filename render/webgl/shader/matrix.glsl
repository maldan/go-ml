precision mediump float;

mat4 translate2(mat4 mx, vec4 v) {
    mx[3][0] = mx[0][0] *v.x + mx[1][0] *v.y + mx[2][0] *v.z + mx[3][0];
    mx[3][1] = mx[0][1] *v.x + mx[1][1] *v.y + mx[2][1] *v.z + mx[3][1];
    mx[3][2] = mx[0][2] *v.x + mx[1][2] *v.y + mx[2][2] *v.z + mx[3][2];
    mx[3][3] = mx[0][3] *v.x + mx[1][3] *v.y + mx[2][3] *v.z + mx[3][3];

    return mx;
}

mat4 translate(vec4 v) {
    return mat4(
        1.0,		0,			0,			0,
        0.0,		1.0,		0,			0,
        0.0,		0,			1.0,		0,
        v.x,		v.y,		v.z,		1.0
    );
}

mat4 rotateX(float angle) {
    /*float s = sin(rad);
    float c = cos(rad);

    mx[1][0] = mx[1][0] * c + mx[2][0] * s;
    mx[1][1] = mx[1][1] * c + mx[2][1] * s;
    mx[1][2] = mx[1][2] * c + mx[2][2] * s;
    mx[1][3] = mx[1][3] * c + mx[2][3] * s;
    mx[2][0] = mx[2][0] * c - mx[1][0] * s;
    mx[2][1] = mx[2][1] * c - mx[1][1] * s;
    mx[2][2] = mx[2][2] * c - mx[1][2] * s;
    mx[2][3] = mx[2][3] * c - mx[1][3] * s;

    return mx;*/
    return mat4(	1.0,		0,			0,			0,
    0, 	cos(angle),	-sin(angle),		0,
    0, 	sin(angle),	 cos(angle),		0,
    0, 			0,			  0, 		1);
}

mat4 rotateY(float angle) {
    /*float s = sin(rad);
    float c = cos(rad);

    mx[0][0] = mx[0][0] * c - mx[2][0] * s;
    mx[0][1] = mx[0][1] * c - mx[2][1] * s;
    mx[0][2] = mx[0][2] * c - mx[2][2] * s;
    mx[0][3] = mx[0][3] * c - mx[2][3] * s;
    mx[2][0] = mx[0][0] * s + mx[2][0] * c;
    mx[2][1] = mx[0][1] * s + mx[2][1] * c;
    mx[2][2] = mx[0][2] * s + mx[2][2] * c;
    mx[2][3] = mx[0][3] * s + mx[2][3] * c;

    return mx;*/
    return mat4(
        cos(angle),		0,		sin(angle),	0,
        0,		1.0,			 0,	0,
        -sin(angle),	0,		cos(angle),	0,
        0, 		0,				0,	1
    );
}

mat4 rotateZ(float angle) {
    /*float s = sin(rad);
    float c = cos(rad);

    mx[0][0] = mx[0][0] * c + mx[1][0] * s;
    mx[0][1] = mx[0][1] * c + mx[1][1] * s;
    mx[0][2] = mx[0][2] * c + mx[1][2] * s;
    mx[0][3] = mx[0][3] * c + mx[1][3] * s;
    mx[1][0] = mx[1][0] * c - mx[0][0] * s;
    mx[1][1] = mx[1][1] * c - mx[0][1] * s;
    mx[1][2] = mx[1][2] * c - mx[0][2] * s;
    mx[1][3] = mx[1][3] * c - mx[0][3] * s;

    return mx;*/
    return mat4(
        cos(angle),		-sin(angle),	    0,	0,
        sin(angle),		cos(angle),		0,	0,
        0,				0,		        1,	0,
        0,				0,		        0,	1
    );
}

vec4 applyMatrixToVec4(mat4 mx, vec4 v) {
    float x = v.x;
    float y = v.y;
    float z = v.z;
    float w = mx[0][3] * x + mx[1][3] * y + mx[2][3] * z + mx[3][3];

    if (w == 0.0) {
        w = 1.0;
    }

    v.x = (mx[0][0] * x + mx[1][0] * y + mx[2][0] * z + mx[3][0]) / w;
    v.y = (mx[0][1] * x + mx[1][1] * y + mx[2][1] * z + mx[3][1]) / w;
    v.z = (mx[0][2] * x + mx[1][2] * y + mx[2][2] * z + mx[3][2]) / w;

    return v;
}

