#version 410
in vec3 position;
in uint vert_color;

flat out uint mouse_state;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projMatrix;
uniform mat4 MVP;



void main() {
    mouse_state=vert_color;
    gl_Position = MVP * vec4(position, 1);
}
