#version 410
in vec3 position;
in vec4 vert_color;

out vec4 color_from_vshader;

void main() {
    color_from_vshader=vert_color;
    gl_Position = vec4(position, 1.0);
}
