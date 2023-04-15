#version 330

uniform vec3 u_cam_pos;
uniform vec3 u_light_pos;
uniform vec3 u_light_intensity;

uniform vec4 u_color;

uniform sampler2D u_texture_2;
uniform vec2 u_texture_2_size;

uniform float u_normal_scaling;
uniform float u_height_scaling;

in vec4 v_position;
in vec4 v_normal;
in vec4 v_tangent;
in vec2 v_uv;

out vec4 out_color;

//Parameters for Bump Mapping
vec3 b;
mat3 TBN;

float first_arg_du;
float dU;
float second_arg_dv;
float dV;
vec3 n_0;
vec3 n_d;

//Parameters for Blinn Phong Shading
//In Piazza, the TA (Divi) said that the spec uses: ka = 0.1, ks = 0.5, p = 100, kd = u_color, and Ia = a vector of all 1s
//Diffusion parameters
float radius;
vec4 diffuse_result;
vec3 l;
vec4 k_d = u_color;

//Ambient lighting parameters
float k_a = 0.15;
vec4 I_a = vec4(1, 1, 1, 1);

//Specular reflection parameters
vec3 v;
vec3 phongH;
vec4 specular_result;
float k_s = 0.9;
float p = 100.0;

vec4 result;
float h(vec2 uv) {
  // You may want to use this helper function...
    //One such h(u, v) you could use would be the r component of the color vector stored in the texture at coordinates (u, v).
    result = texture(u_texture_2, uv);
    return result.r;
}

void main() {
  // YOUR CODE HERE
  //TBN Matrix -> transforms the something from object space into model space
    b = cross(vec3(v_normal.xyz), vec3(v_tangent.xyz));
    TBN = mat3(vec3(v_tangent.xyz), b, vec3(v_normal.xyz));

  //Calculating local space normal, n_0
    //dU = ( h(u + (1.0/w), v ) - h(u,v)) * k_h * k_n
    //dV = ( h(u, v + (1.0/h) - h(u,v)) * k_h * k_n
    //n_0 = (-dU, -dV, 1)

    first_arg_du = v_uv.x + (1.0/u_texture_2_size.x);
    dU = ( h(vec2(first_arg_du, v_uv.y)) - h(v_uv) ) * u_height_scaling * u_normal_scaling;

    second_arg_dv = v_uv.y + (1.0/u_texture_2_size.y);
    dV = ( h(vec2(v_uv.x, second_arg_dv)) - h(v_uv) ) * u_height_scaling * u_normal_scaling;

    n_0 = vec3((-1.0 * dU), (-1.0 * dV), 1);

    //Calculating displaced model space normal, n_d
    n_d = TBN * n_0;

    //phong stuff copied over
    radius = distance(vec4(u_light_pos,0), v_position);
    l = (u_light_pos - vec3(v_position.xyz)) / distance(vec3(v_position.xyz), u_light_pos);
    diffuse_result = k_d * (vec4(u_light_intensity,0) / (radius * radius)) * max(0, dot(normalize(vec3(n_d.xyz)), l));

    v = (u_cam_pos - v_position.xyz) / distance(u_cam_pos, vec3(v_position.xyz));
    phongH  = (v + l) / length(v + l);
    specular_result = k_s * (vec4(u_light_intensity,0) / (radius * radius)) * pow( max(0, dot(normalize(vec3(n_d.xyz)), phongH)), p);

    //out_color = Ambient lighting + diffusion lighting + specular lighting
    out_color = (k_a * I_a) + diffuse_result + specular_result;
    out_color.a = 1;
}
