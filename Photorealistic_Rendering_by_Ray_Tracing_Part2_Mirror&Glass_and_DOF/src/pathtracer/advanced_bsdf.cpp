#include "bsdf.h"

#include <algorithm>
#include <iostream>
#include <utility>

#include "application/visual_debugger.h"

using std::max;
using std::min;
using std::swap;

namespace CGL {

// Mirror BSDF //

Vector3D MirrorBSDF::f(const Vector3D wo, const Vector3D wi) {
  return Vector3D();
}

Vector3D MirrorBSDF::sample_f(const Vector3D wo, Vector3D* wi, double* pdf) {

  // TODO Project 3-2: Part 1
  // Implement MirrorBSDF
    (*pdf) = 1.0;
    reflect(wo, wi);
    return reflectance / abs_cos_theta(*wi);;
}

void MirrorBSDF::render_debugger_node()
{
  if (ImGui::TreeNode(this, "Mirror BSDF"))
  {
    DragDouble3("Reflectance", &reflectance[0], 0.005);
    ImGui::TreePop();
  }
}

// Microfacet BSDF //

double MicrofacetBSDF::G(const Vector3D wo, const Vector3D wi) {
  return 1.0 / (1.0 + Lambda(wi) + Lambda(wo));
}

double MicrofacetBSDF::D(const Vector3D h) {
  // TODO Project 3-2: Part 2
  // Compute Beckmann normal distribution function (NDF) here.
  // You will need the roughness alpha.
  return 1.0;
}

Vector3D MicrofacetBSDF::F(const Vector3D wi) {
  // TODO Project 3-2: Part 2
  // Compute Fresnel term for reflection on dielectric-conductor interface.
  // You will need both eta and etaK, both of which are Vector3D.

  return Vector3D();
}

Vector3D MicrofacetBSDF::f(const Vector3D wo, const Vector3D wi) {
  // TODO Project 3-2: Part 2
  // Implement microfacet model here.

  return Vector3D();
}

Vector3D MicrofacetBSDF::sample_f(const Vector3D wo, Vector3D* wi, double* pdf) {
  // TODO Project 3-2: Part 2
  // *Importance* sample Beckmann normal distribution function (NDF) here.
  // Note: You should fill in the sampled direction *wi and the corresponding *pdf,
  //       and return the sampled BRDF value.

  *wi = cosineHemisphereSampler.get_sample(pdf);
  return MicrofacetBSDF::f(wo, *wi);
}

void MicrofacetBSDF::render_debugger_node()
{
  if (ImGui::TreeNode(this, "Micofacet BSDF"))
  {
    DragDouble3("eta", &eta[0], 0.005);
    DragDouble3("K", &k[0], 0.005);
    DragDouble("alpha", &alpha, 0.005);
    ImGui::TreePop();
  }
}

// Refraction BSDF //

Vector3D RefractionBSDF::f(const Vector3D wo, const Vector3D wi) {
  return Vector3D();
}

Vector3D RefractionBSDF::sample_f(const Vector3D wo, Vector3D* wi, double* pdf) {
  // TODO Project 3-2: Part 1
  // Implement RefractionBSDF
    (*pdf) = 1;
    
    //refract() initializes wi and returns false if refraction does not happen due to total internal reflection
    bool is_refracting = refract(wo, wi, ior);
    if (!is_refracting) {
        return Vector3D();
    } else {
        
        //Calculate eta (same as the refract())
        double eta = get_eta(wo, ior);
        return transmittance /abs_cos_theta(*wi) / (eta * eta); //refraction function
    }
}

void RefractionBSDF::render_debugger_node()
{
  if (ImGui::TreeNode(this, "Refraction BSDF"))
  {
    DragDouble3("Transmittance", &transmittance[0], 0.005);
    DragDouble("ior", &ior, 0.005);
    ImGui::TreePop();
  }
}

// Glass BSDF //

Vector3D GlassBSDF::f(const Vector3D wo, const Vector3D wi) {
  return Vector3D();
}

Vector3D GlassBSDF::sample_f(const Vector3D wo, Vector3D* wi, double* pdf) {

  // TODO Project 3-2: Part 1
  // Compute Fresnel coefficient and either reflect or refract based on it.

  // compute Fresnel coefficient and use it as the probability of reflection
  // - Fundamentals of Computer Graphics page 305

    //If total internal reflection does occur, BOTH reflection and refraction will occur
    bool no_total_internal_reflecting = refract(wo, wi, ior);
    
    if (!no_total_internal_reflecting) {
        (*pdf) = 1;
        reflect(wo, wi);
        return reflectance / abs_cos_theta(*wi);
        
    } else {
        
        double cosineTheta = abs_cos_theta(wo);
        double R0 = ((ior - 1)/(ior + 1)) * ((ior - 1)/(ior + 1)); //TA said to use ior instead
        double R = R0 + ( (1.0 - R0) * pow( (1.0 - cosineTheta) , 5) );

        //coinflip whether to reflect or refract, true = reflect, false = refract
        if (coin_flip(R)) {
            
            reflect(wo, wi);
            (*pdf) = R;
            return R * reflectance / abs_cos_theta(*wi);
            
        } else {
            
            (*pdf) = 1.0 - R;
            double eta = get_eta(wo, ior);
            refract_helper_function(wo, wi, eta);
            return (1.0 - R) * transmittance / abs_cos_theta(*wi) / (eta * eta);
        }
    }
}

void GlassBSDF::render_debugger_node()
{
  if (ImGui::TreeNode(this, "Refraction BSDF"))
  {
    DragDouble3("Reflectance", &reflectance[0], 0.005);
    DragDouble3("Transmittance", &transmittance[0], 0.005);
    DragDouble("ior", &ior, 0.005);
    ImGui::TreePop();
  }
}

void BSDF::reflect(const Vector3D wo, Vector3D* wi) {

  // TODO Project 3-2: Part 1
  // Implement reflection of wo about normal (0,0,1) and store result in wi.
    Vector3D normal = Vector3D(0.0, 0.0, 1.0);
    (*wi) = ( -1 * wo ) + (2 * dot(wo, normal) * normal); //reflection formula
}

bool BSDF::refract(const Vector3D wo, Vector3D* wi, double ior) {

  // TODO Project 3-2: Part 1
  // Use Snell's Law to refract wo surface and store result ray in wi.
  // Return false if refraction does not occur due to total internal reflection
  // and true otherwise. When dot(wo,n) is positive, then wo corresponds to a
  // ray entering the surface through vacuum.
    
    double eta = get_eta(wo, ior);
    return refract_helper_function(wo, wi, eta);
}

//CREATED HELPER FUNCTION
double BSDF::get_eta(const Vector3D wo, double ior){
    Vector3D normal = Vector3D(0.0, 0.0, 1.0);
    double eta= 0.0;
    
    //When wo's z is positive, we say then we are ENTERING the non-air material
    if (dot(wo,normal) >= 0) {
        eta = 1.0 / ior; //Index of refraction, N, is 1/ior when entering
        
    } else { //Otherwise, we are EXITING the non-air material
        eta = ior / 1.0; //Index of refraction, N, is ior when exiting
    }
    
    return eta;
}


bool BSDF::refract_helper_function(const Vector3D wo, Vector3D* wi, double eta) {
   
    //If there is total internal reflection, then return false and leave wi unused
    double total_internal_reflection = 1.0 - ( (eta * eta) * (1.0 - (wo.z * wo.z)) );
    if (total_internal_reflection < 0) {
        return false;
    } else {
        (*wi).x = -1.0 * eta * wo.x;
        (*wi).y = -1.0 * eta * wo.y;
        
        double temp = 1.0 - ( (eta * eta) * (1.0 - (wo.z * wo.z)) );
        temp = sqrt(temp);
        
    
        if (wo.z >= 0) {
            (*wi).z = -1.0 * temp;
        } else {
            (*wi).z = temp;
        }
         
        //(*wi).z = -(abs(wo.z)/wo.z) * temp;
        
        return true;
    }
  }

} // namespace CGL
