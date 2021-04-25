#version 410
in vec3 position;
in uint vert_color;

flat out uint mouse_state;

void main() {
    mouse_state=vert_color;
    gl_Position = vec4(position, 1.0);
}
