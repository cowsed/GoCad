#version 410

in vec4 color_from_vshader;
out vec4 frag_colour;

uniform vec3 normal_color;
uniform vec3 selected_color;

void main() {
    frag_colour = vec4(selected_color.x, selected_color.y, selected_color.z, 1);
}
