#version 330

uniform vec4 u_color;
uniform vec3 u_cam_pos;
uniform vec3 u_light_pos;
uniform vec3 u_light_intensity;

in vec4 v_position;
in vec4 v_normal;
in vec2 v_uv;

out vec4 out_color;


//Parameters for Blinn Phong Shading
//In Piazza, the TA (Divi) said that the spec uses: ka = 0.1, ks = 0.5, p = 100, kd = u_color, and Ia = a vector of all 1s
//Diffusion parameters
float radius;
vec4 diffuse_result;
vec3 l;
vec4 k_d = u_color;

//Ambient lighting parameters
float k_a = 0.15;
vec4 I_a = normalize(vec4(1, 1, 1, 1));

//Specular reflection parameters
vec3 v;
vec3 h;
vec4 specular_result;
float k_s = 0.9;
float p = 100.0;


/* FUNCTIONS:
 - length() takes in one vector and returns the magnitude of that.

 - distance() takes in two and returns the distance between those. Basically, distance(a, b) is the same as length(a - b)
 
 */
void main() {
  // YOUR CODE HERE
    //u_light_pos and v_position are already in world space. This should also help with calculating r squared.
    radius = distance(vec4(u_light_pos,0), v_position);
    l = (u_light_pos - vec3(v_position.xyz)) / distance(vec3(v_position.xyz), u_light_pos);
    diffuse_result = k_d * (vec4(u_light_intensity,0) / (radius * radius)) * max(0, dot(normalize(vec3(v_normal.xyz)), l));
    
    v = (u_cam_pos - v_position.xyz) / distance(u_cam_pos, vec3(v_position.xyz));
    h  = (v + l) / length(v + l);
    specular_result = k_s * (vec4(u_light_intensity,0) / (radius * radius)) * pow( max(0, dot(normalize(vec3(v_normal.xyz)), h)), p);
    
   //out_color = Ambient lighting + diffusion lighting + specular lighting
    out_color = (k_a * I_a) + diffuse_result + specular_result;
    out_color.a = 1;
}
