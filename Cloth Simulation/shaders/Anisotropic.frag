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


/******** Variables for Anisotropic Cloth Shading needed in Helper and Main Func ********/

//Amount of anisotropy. Scalar between −1 and 1
float anisotropy = u_normal_scaling;
float roughness = u_height_scaling*u_height_scaling; //Same component used for bump mapping

//luminance of base color
//vec3 sheenColor = vec3(0.2126*u_color.x, 0.7152*u_color.y, 0.0722*u_color.z);

//defult
vec3 diffuseColor = vec3(0.24,0,0.03);
vec3 sheenColor = u_color.xyz;

/******** Distrubution Helper Fun ********/
float distribution(float NoH, const vec3 h, const vec3 t, const vec3 b, float at, float ab) {
  float ToH = dot(t, h);
  float BoH = dot(b, h);
  float a2 = at * ab;
  highp vec3 v = vec3(ab * ToH, at * BoH, a2 * NoH);
  highp float v2 = dot(v, v);
  float w2 = a2 / v2;
  return a2 * w2 * w2 * (1.0 / PI);
}

/******** Visibility Helper Fun ********/
float visibility(float at, float ab, float ToV, float BoV, float ToL, float BoL, float NoV, float NoL){
  float lL = NoL* length(vec3(at*ToV, ab*BoV, NoV));
  float lV = NoV * length(vec3(at*ToL, ab*BoL, NoL));
  float v = 0.5/(lL+lV);
  return min(v, 65504.0);
}


/******** Texture Helper Fun ********/
//Denim = texture_3.png
//Silk = texture_4.png
float bumpHelper(vec2 uv) {
  // You may want to use this helper function...
  return texture(u_texture_4,uv).r;
}

void main() {

  /******** Bump Code ********/
  //view unit vector
  vec3 v = normalize(u_cam_pos-v_position.xyz); //viewing vector

  //incident light unit vector
  vec3 r = u_light_pos - v_position.xyz;
  vec3 l = normalize(r); //light source vecor

  //half unit vector between l and v
  vec3 h = normalize(v + l); //half-angle vector

  //find b and create matrix
  vec3 b = cross(v_normal.xyz,v_tangent.xyz);
  mat3x3 tbn = mat3x3(v_tangent.xyz,b,v_normal.xyz);

  float du = (bumpHelper(vec2(v_uv.x + 1.0/u_texture_2_size.x, v_uv.y)) - bumpHelper(v_uv))*u_height_scaling*u_normal_scaling;
  float dv = (bumpHelper(vec2(v_uv.x,v_uv.y+1.0/u_texture_2_size.y)) - bumpHelper(v_uv))*u_height_scaling*u_normal_scaling;

  vec3 n0 = vec3(-du,-dv,1.0);
  vec3 nd = tbn*n0;

  /******** Anisotrophic Lighting For Cloth ********/

  // helpful definitions
  float noh = dot(v_normal.xyz, h);
  float nov = abs(dot(v_normal.xyz, v));
  float nol = dot(v_normal.xyz, l);
  float loh = dot(l,h);



  float at = max(roughness * (1.0 + anisotropy), 0.001);
  float ab = max(roughness * (1.0 - anisotropy), 0.001);
  vec3 t = v_tangent.xyz;
  vec3 b2 = cross(v_normal.xyz, t);

  float tov = dot(t,v);
  float bov = dot(b,v);
  float tol = dot(t,l);
  float bol = dot(b,l);


  //note for anna, i used h=u_height_scaling and t = v_tangent, is this correct?
  //h should be the half unit vector which we already calculated in the above code
  //b is the bitangent we calculated it to get the matrix for bump I recalculated it with the new normal
  //unsure of if it will really make a difference tho
  // t should be v_tangent tho

  //Specular
  float D = distribution(noh, h, t, b2, at, ab);  //aniso distribution
  float V = visibility(at, ab, tov, bov, tol, bol, nov, nol); //aniso visibility
  vec3 F = sheenColor;
  vec3 Fr = (D*V)*F ; //Specular component (f_r)


  //Diffuse
  vec3 Fd = diffuseColor/PI; //Diffusion component (f_d)
  vec3 color = (Fr * u_light_intensity/4.0) + (Fd * u_light_intensity); //Implementing Blinn Phong Model (without ambient lighting)
  color *= max(0.0, nol);
    
  out_color = vec4(color,1);
  out_color.a = 1.0;
}
