#version 330
#define PI 3.1415926538

// (Every uniform is available here.)
uniform mat4 u_view_projection;
uniform mat4 u_model;

uniform float u_normal_scaling;
uniform float u_height_scaling;

uniform vec3 u_cam_pos;
uniform vec3 u_light_pos;
uniform vec3 u_light_intensity;
uniform vec4 u_color;

// Feel free to add your own textures. If you need more than 4,
// you will need to modify the skeleton.
uniform sampler2D u_texture_1;
uniform sampler2D u_texture_2;
uniform sampler2D u_texture_3;
uniform sampler2D u_texture_4;
uniform vec2 u_texture_2_size;

// Environment map! Take a look at GLSL documentation to see how to
// sample from this.
uniform samplerCube u_texture_cubemap;

in vec4 v_position;
in vec4 v_normal;
in vec4 v_tangent;
in vec2 v_uv;

out vec4 out_color;


/******** Variables for Cloth Shading needed in Helper and Main Func ********/
//boolean to determine if the material should have subsurface scattering
bool subsurf_color = true;
//roughness, may want to change later based upon the system in place though this might be good
float roughness = u_height_scaling*u_normal_scaling; //Same component used for bump mapping

//luminance of base color
vec3 sheenColor = vec3(0.5294*u_color.x, 0.8078*u_color.y, 0.92164*u_color.z);

//defult
//vec3 sheenColor = vec3(0.04,0.04, 0.04);
vec3 diffuseColor = u_color.xyz;

/******** Distrubution Helper Fun ********/
float distribution(float rough, float NoH){
  float sin_theta = max((1.0-(NoH*NoH)), 0.0078125);
  return (((2.0+(1.0/rough))*pow(sin_theta,(1.0/rough)*0.5))/(2.0*PI));
}

/******** Visibility Helper Fun ********/
float visibility(float NoV, float NoL){
  float denom = 4.0*(NoL+NoV-(NoL*NoV));
  return (1.0/denom);
}

/******** Texture Helper Fun ********/
float help(vec2 uv) {
  // You may want to use this helper function...
  return texture(u_texture_3,uv).r;
}

void main() {

  //view unit vector
  vec3 v = normalize(u_cam_pos-v_position.xyz); //viewing vector

  //incident light unit vector
  vec3 r = u_light_pos - v_position.xyz;
  vec3 l = normalize(r); //light source vecor

  //half unit vector between l and v
  vec3 h = normalize(v + l); //half-angle vector

  //Pasted Code From Bump
  //there in order to hopefully create a more intersting and realistic texture for fabric
  //find b and create matrix
  vec3 b = cross(v_normal.xyz,v_tangent.xyz);
  mat3x3 tbn = mat3x3(v_tangent.xyz,b,v_normal.xyz);

  float du = (help(vec2(v_uv.x + 1.0/u_texture_2_size.x, v_uv.y)) - help(v_uv))*roughness;
  float dv = (help(vec2(v_uv.x,v_uv.y+1.0/u_texture_2_size.y)) - help(v_uv))*roughness;

  vec3 n0 = vec3(-du,-dv,1.0);
  vec3 nd = tbn*n0;

  //other helpful definitions
  float noh = dot(normalize(nd), h);
  float nov = abs(dot(normalize(nd), v));
  float nol = dot(normalize(nd), l);
  float loh = dot(l,h);

  /******** New Shading Method For Cloth ********/

  //Specular
  float D = distribution(roughness, noh);  //Charlie Distribution
  float V = visibility(nov,nol); //Denom component of specular f_r function
  vec3 F = sheenColor;
  vec3 Fr = (D*V)*F; //Specular component (f_r)

  //Diffuse from diffuse frag
  //vec3 kd = diffuseColor;
  //float r2 = length(r)*length(r);
  vec3 Fd = diffuseColor/PI; //Diffusion component (f_d)
  vec3 color = Fr + Fd; //Implementing Blinn Phong Model (without ambient lighting)
    color *= u_light_intensity * max(0.0, nol);

  out_color = vec4(color,1);
  out_color.a = 1.0;
}
