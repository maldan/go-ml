mat4 translate(mat4 mx, vec4 v) {
    mx[3][0] = mx[0][0] *v.x + mx[1][0] *v.y + mx[2][0] *v.z + mx[3][0];
    mx[3][1] = mx[0][1] *v.x + mx[1][1] *v.y + mx[2][1] *v.z + mx[3][1];
    mx[3][2] = mx[0][2] *v.x + mx[1][2] *v.y + mx[2][2] *v.z + mx[3][2];
    mx[3][3] = mx[0][3] *v.x + mx[1][3] *v.y + mx[2][3] *v.z + mx[3][3];

    return mx;
}