#include "sphere.h"

#include <cmath>

#include "pathtracer/bsdf.h"
#include "util/sphere_drawing.h"

namespace CGL {
namespace SceneObjects {

bool Sphere::test(const Ray &r, double &t1, double &t2) const {

  // TODO (Part 1.4):
  // Implement ray - sphere intersection test.
  // Return true if there are intersections and writing the
  // smaller of the two intersection times in t1 and the larger in t2.


  return true;

}

bool Sphere::has_intersection(const Ray &r) const {

  // TODO (Part 1.4):
  // Implement ray - sphere intersection.
  // Note that you might want to use the the Sphere::test helper here.

    //Get the coordinates of: (o + td -c)^2 - R^2 = 0
    double a =  dot(r.d , r.d);
    double b = dot((2.0 * (r.o - (this ->o))), r.d);
    double temp = dot((r.o - (this ->o)), (r.o - (this ->o)));
    double c = temp - (this->r2);
    
    //Solving the quadratic formula. REMEMBER: its b PLUS OR MINUS!
        //-> Need to caluclate t for both values!
    temp = (b*b) - (4.0 * a * c);
    temp = sqrt(temp); // square root of (b^2 - 4ac)
    double numerator = -b + temp;
    double denominator = 2.0 * a;
    
    double larger_t = numerator / denominator;
    numerator = -b - temp;
    double smaller_t = numerator /denominator;
    
    if (smaller_t >= r.min_t && smaller_t <= r.max_t) {
        //Update ray's max_t to be the CLOSEST intersection point!
        r.max_t = smaller_t;
        return true;
    } else if (larger_t >= r.min_t && larger_t <= r.max_t) {  //ELSE check if the OTHER intersection point (thats larger than smaller_t) is within range
        r.max_t = larger_t;
        return true;
    }
    
  return false;
}

bool Sphere::intersect(const Ray &r, Intersection *i) const {

  // TODO (Part 1.4):
  // Implement ray - sphere intersection.
  // Note again that you might want to use the the Sphere::test helper here.
  // When an intersection takes place, the Intersection data should be updated
  // correspondingly.

    //Get the coordinates of: (o + td -c)^2 - R^2 = 0
    double a =  dot(r.d , r.d);
    double b = dot((2.0 * (r.o - (this ->o))), r.d);
    double temp = dot((r.o - (this ->o)), (r.o - (this ->o)));
    double c = temp - (this->r2);
    
    //Solving the quadratic formula. REMEMBER: its b PLUS OR MINUS!
        //-> Need to caluclate t for both values!
    temp = (b*b) - (4.0 * a * c);
    temp = sqrt(temp);
    double numerator = -b + temp;
    double denominator = 2.0 * a;
    
    double larger_t = numerator / denominator;
    numerator = -b - temp;
    double smaller_t = numerator /denominator;
    
    
    if (smaller_t >= r.min_t && smaller_t <= r.max_t) {
        //Update ray's max_t to be the CLOSEST intersection point!
        r.max_t = smaller_t;
        
        //Populate Intersection struct
        /*
         - t is the t-value of the input ray where the intersection occurs.
         - n is the surface normal at the intersection. You should use barycentric coordinates to interpolate the three vertex normals of the triangle, n1, n2, and n3.
         - primitive points to the primitive that was intersected (use the this pointer).
         - bsdf points to the surface material (BSDF) at the hit point (use the get_bsdf() method). BSDF stands for Bidirectional Scattering Distribution Function, a generalization of BRDF that accounts for both reflection and transmission.
         */
        
        //Calculate the normal vector! normal vector = point of intersection - center of spher
        Vector3D point_of_intersect = r.at_time(smaller_t);
        Vector3D normal_vector = point_of_intersect - this -> o;
        normal_vector.normalize();
        
        i -> t = smaller_t;
        i -> n = normal_vector;
        i -> primitive = this;
        i -> bsdf = this -> get_bsdf();
        
        return true;
    } else if (larger_t >= r.min_t && larger_t <= r.max_t) {  //ELSE check if the OTHER intersection point (thats larger than smaller_t) is within range
        r.max_t = larger_t;
        
        //Populate Intersection struct
        /*
         - t is the t-value of the input ray where the intersection occurs.
         - n is the surface normal at the intersection. You should use barycentric coordinates to interpolate the three vertex normals of the triangle, n1, n2, and n3.
         - primitive points to the primitive that was intersected (use the this pointer).
         - bsdf points to the surface material (BSDF) at the hit point (use the get_bsdf() method). BSDF stands for Bidirectional Scattering Distribution Function, a generalization of BRDF that accounts for both reflection and transmission.
         */
        
        //Calculate the normal vector! normal vector = point of intersection - center of spher
        Vector3D point_of_intersect = r.at_time(larger_t);
        Vector3D normal_vector = point_of_intersect - this -> o;
        normal_vector.normalize();
        
        i -> t = larger_t;
        i -> n = normal_vector;
        i -> primitive = this;
        i -> bsdf = this -> get_bsdf();
        
        return true;
    }
    
  return false;
}

void Sphere::draw(const Color &c, float alpha) const {
  Misc::draw_sphere_opengl(o, r, c);
}

void Sphere::drawOutline(const Color &c, float alpha) const {
  // Misc::draw_sphere_opengl(o, r, c);
}

} // namespace SceneObjects
} // namespace CGL
