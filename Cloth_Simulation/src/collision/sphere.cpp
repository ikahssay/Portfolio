#include <nanogui/nanogui.h>

#include "../clothMesh.h"
#include "../misc/sphere_drawing.h"
#include "sphere.h"

using namespace nanogui;
using namespace CGL;

void Sphere::collide(PointMass &pm) {
  // TODO (Part 3): Implement Sphere::collide, which takes in a point mass and adjusts its position if it intersects with or is inside the sphere. If the point mass intersects with or is inside the sphere, then "bump" it up to the surface of the sphere:
    
    //1. Compute where the point mass should have intersected the sphere by extending the path between its 'position' and the sphere's origin to the sphere's surface. Call the surface intersection point the tangent point.
        //-> You can check to see if the distance between the point masses' position and the sphere's origin is less than the sphere's radius. This would mean the point is within the sphere.
    double pm_and_sphere_distance = (this -> origin - pm.position).norm(); //Distance = || pos_b - pos_a ||
    if (pm_and_sphere_distance <= this -> radius) {
        
        //2. Compute the correction vector needed to be applied to the point mass's last_position in order to reach the tangent point.
        Vector3D correction_difference = (this->radius - pm_and_sphere_distance) * (pm.last_position - this-> origin).unit(); //unit() vector is pointing opposite direction from sphere so cloth can move distance (this->radius - pm_and_sphere_distance) away from sphere
        
        //3. Finally, let the point mass's new position be its last_position adjusted by the above correction vector, scaled down by friction (i.e. scaled by (1 - f)).
        correction_difference = (1 - this->friction)*correction_difference;
        pm.position =  pm.last_position + correction_difference; //updates position to be top of surface
    }
    
    //4. Make sure to update Cloth::simulate to account for potential collisions.

}

void Sphere::render(GLShader &shader) {
  // We decrease the radius here so flat triangles don't behave strangely
  // and intersect with the sphere when rendered
  m_sphere_mesh.draw_sphere(shader, origin, radius * 0.92);
}
