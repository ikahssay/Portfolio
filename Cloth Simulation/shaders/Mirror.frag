#version 330


uniform vec3 u_cam_pos;

uniform samplerCube u_texture_cubemap;

in vec4 v_position;
in vec4 v_normal;
in vec4 v_tangent;

out vec4 out_color;

//Parameters for Mirror.frag
vec3 w_out;
vec3 w_in;

void main() {
//1. Using the camera's position u_cam_pos and the fragment's position v_position, compute the outgoing eye-ray, w_o.
  w_out = (u_cam_pos - v_position.xyz) / distance(u_cam_pos, vec3(v_position.xyz));
  w_out = normalize(w_out);
    
//2. Then reflect w_o across the surface normal given in v_normal to get w_i.
    w_in = (-1 * w_out) + ( 2 * dot(w_out, vec3(v_normal.xyz)) * vec3(v_normal.xyz) ); //reflection equation

//3. Finally, sample the environment map u_texture_cubemap for the incoming direction w_i
   out_color = texture(u_texture_cubemap, w_in);
    
  // YOUR CODE HERE
  //out_color = (vec4(1, 1, 1, 0) + v_normal) / 2;
  out_color.a = 1;
}
