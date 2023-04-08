#include "triangle.h"

#include "CGL/CGL.h"
#include "GL/glew.h"

namespace CGL {
namespace SceneObjects {

Triangle::Triangle(const Mesh *mesh, size_t v1, size_t v2, size_t v3) {
  p1 = mesh->positions[v1];
  p2 = mesh->positions[v2];
  p3 = mesh->positions[v3];
  n1 = mesh->normals[v1];
  n2 = mesh->normals[v2];
  n3 = mesh->normals[v3];
  bbox = BBox(p1);
  bbox.expand(p2);
  bbox.expand(p3);

  bsdf = mesh->get_bsdf();
}

BBox Triangle::get_bbox() const { return bbox; }

bool Triangle::has_intersection(const Ray &r) const {
  // Part 1, Task 3: implement ray-triangle intersection
  // The difference between this function and the next function is that the next
  // function records the "intersection" while this function only tests whether
  // there is a intersection.

    //Get interpolated normalized vector on triangle
    //User Molar Trumbore Algorithm to get barycentric coords b1 and b2 (which is beta & gamma)
        //-> also outputs time t of intersection
    //Get interpolated normalized vector on triangle
    //User Molar Trumbore Algorithm to get barycentric coords b1 and b2 (which is beta & gamma)
        //-> also outputs time t of intersection
    Vector3D e1 = (this -> p2) - (this ->p1);
    Vector3D e2 = (this -> p3) - (this ->p1);
    Vector3D s = r.o - (this -> p1);
    Vector3D s1 = cross(r.d, e2);
    Vector3D s2 = cross(s, e1);
    
    double x = dot(s2 , e2);
    double y = dot(s1 , s);
    double z = dot(s2 , r.d);
    Vector3D transform = Vector3D(x, y, z);
    double divide = 1/ ( dot(s1 , e1) );
    Vector3D moller_trumbore_vector = divide * transform;
    
    double t = moller_trumbore_vector.x;
    double b1 = moller_trumbore_vector.y; // beta
    double b2 = moller_trumbore_vector.z; //gamma
    
    double alpha = 1 - b1 - b2;
    Vector3D normalized_normal_vector = (n1 * alpha) + (n2 * b1) + (n3 * b2);
    
    //Check if barycentric coords are valid!
    if (alpha > 0 && alpha < 1 && b1 > 0 && b1 < 1 && b2 > 0 && b2 < 1) {
        if (t >= r.min_t && t <= r.max_t) {
            //Update ray's max_t
            r.max_t = t;
            return true;
        }
        
    }
    
        return false;
}

bool Triangle::intersect(const Ray &r, Intersection *isect) const {
  // Part 1, Task 3:
  // implement ray-triangle intersection. When an intersection takes
  // place, the Intersection data should be updated accordingly

    //Get interpolated normalized vector on triangle
    //User Molar Trumbore Algorithm to get barycentric coords b1 and b2 (which is beta & gamma)
        //-> also outputs time t of intersection
    Vector3D e1 = (this -> p2) - (this ->p1);
    Vector3D e2 = (this -> p3) - (this ->p1);
    Vector3D s = r.o - (this -> p1);
    Vector3D s1 = cross(r.d, e2);
    Vector3D s2 = cross(s, e1);
    
    double x = dot(s2 , e2);
    double y = dot(s1 , s);
    double z = dot(s2 , r.d);
    Vector3D transform = Vector3D(x, y, z);
    double divide = 1/ ( dot(s1 , e1) );
    Vector3D moller_trumbore_vector = divide * transform;
    
    double t = moller_trumbore_vector.x;
    double b1 = moller_trumbore_vector.y; // beta
    double b2 = moller_trumbore_vector.z; //gamma
    
    double alpha = 1 - b1 - b2;
    Vector3D normalized_normal_vector = (n1 * alpha) + (n2 * b1) + (n3 * b2);
    
    //Check if barycentric coords are valid!
    if (alpha > 0 && alpha < 1 && b1 > 0 && b1 < 1 && b2 > 0 && b2 < 1) {
        if (t >= r.min_t && t <= r.max_t) {
            //Update ray's max_t
            r.max_t = t;
            
            //Populate Intersection struct
            /*
             - t is the t-value of the input ray where the intersection occurs.
             - n is the surface normal at the intersection. You should use barycentric coordinates to interpolate the three vertex normals of the triangle, n1, n2, and n3.
             - primitive points to the primitive that was intersected (use the this pointer).
             - bsdf points to the surface material (BSDF) at the hit point (use the get_bsdf() method). BSDF stands for Bidirectional Scattering Distribution Function, a generalization of BRDF that accounts for both reflection and transmission.
             */
            isect -> t = t;
            isect -> n = normalized_normal_vector;
            isect -> primitive = this;
            isect -> bsdf = get_bsdf();
            
            return true;
        }
        
    }
    

    return false;
    
}

void Triangle::draw(const Color &c, float alpha) const {
  glColor4f(c.r, c.g, c.b, alpha);
  glBegin(GL_TRIANGLES);
  glVertex3d(p1.x, p1.y, p1.z);
  glVertex3d(p2.x, p2.y, p2.z);
  glVertex3d(p3.x, p3.y, p3.z);
  glEnd();
}

void Triangle::drawOutline(const Color &c, float alpha) const {
  glColor4f(c.r, c.g, c.b, alpha);
  glBegin(GL_LINE_LOOP);
  glVertex3d(p1.x, p1.y, p1.z);
  glVertex3d(p2.x, p2.y, p2.z);
  glVertex3d(p3.x, p3.y, p3.z);
  glEnd();
}

} // namespace SceneObjects
} // namespace CGL
