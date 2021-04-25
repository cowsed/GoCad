#version 410
flat in uint mouse_state;
out vec4 frag_colour;

uniform vec3 normal_color; //0
uniform vec3 selected_color; //1
uniform vec3 hovered_color; //2

void main() {
    vec3 cols[3];
    cols[0]=normal_color;
    cols[1]=selected_color;
    cols[2]=hovered_color;
    //Combine colors 
    vec3 col = cols[mouse_state];
     

    frag_colour = vec4(col.x, col.y, col.z, 1);
}
