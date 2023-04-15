#version 330

// The camera's position in world-space
uniform vec3 u_cam_pos;

// Color
uniform vec4 u_color;

// Properties of the single point light
uniform vec3 u_light_pos;
uniform vec3 u_light_intensity;

// We also get the uniform texture we want to use.
uniform sampler2D u_texture_1;

// These are the inputs which are the outputs of the vertex shader.
in vec4 v_position;
in vec4 v_normal;

// This is where the final pixel color is output.
// Here, we are only interested in the first 3 dimensions (xyz).
// The 4th entry in this vector is for "alpha blending" which we
// do not require you to know about. For now, just set the alpha
// to 1.
out vec4 out_color;

float radius;
vec4 diffuse_result;
vec3 l;

void main() {
  // YOUR CODE HERE
   //u_light_pos and v_position are already in world space. This should also help with calculating r squared.
    radius = distance(vec4(u_light_pos,0), v_position);
    l = (u_light_pos - vec3(v_position.xyz)) / distance(vec3(v_position.xyz), u_light_pos);
    
    diffuse_result = 1.0 * (vec4(u_light_intensity,0) / (radius * radius)) * max(0, dot(normalize(vec3(v_normal.xyz)), l));
    
   //The out_color should be the u_color scaled by your calculated diffuse value for each color channel.
    out_color = diffuse_result * u_color;
    out_color.a = 1;
    
  // (Placeholder code. You will want to replace it.)
  //out_color = (vec4(1, 1, 1, 0) + v_normal) / 2;
  //out_color.a = 1;
}
