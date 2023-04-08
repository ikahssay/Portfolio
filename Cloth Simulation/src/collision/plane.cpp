#include "iostream"
#include <nanogui/nanogui.h>

#include "../clothMesh.h"
#include "../clothSimulator.h"
#include "plane.h"

using namespace std;
using namespace CGL;

#define SURFACE_OFFSET 0.0001

void Plane::collide(PointMass &pm) {
  // TODO (Part 3): Handle collision with plane -> Adjust the point mass's its position if it is "inside" the plane, which we define to be when the point moves from one side of the plane to the other in the last time step. If the point mass crosses over, then we "bump" it back up to the side of the surface it originated from:

    //1. Compute where the point mass should have intersected the plane, if it had travelled in a straight line from its position towards the plane. Call this the tangent point.
    
    Vector3D pm_to_plane = this -> point - pm.position; //Line of pm to point in plane (the V vector in lecture 2: https://cs184.eecs.berkeley.edu/sp22/lecture/2-44/digital-drawing
    
    if (dot(pm_to_plane, this ->normal) >= 0 ){
        
        //2. Compute the correction vector needed to be applied to the point mass's last_position in order to reach a point slightly above the tangent point, on the same side of the plane that the point mass was before crossing over. We have provided a small constant SURFACE_OFFSET for this small displacement.
        Vector3D correction_difference = SURFACE_OFFSET * this -> normal;
        
        //3. Finally, let the point mass's new position be its last_position adjusted by the above correction vector, scaled down by friction (i.e. scaled by (1 - f)).
        correction_difference = (1 - this->friction)*correction_difference;
        pm.position =  pm.last_position + correction_difference;
       
    }
}

void Plane::render(GLShader &shader) {
  nanogui::Color color(0.7f, 0.7f, 0.7f, 1.0f);

  Vector3f sPoint(point.x, point.y, point.z);
  Vector3f sNormal(normal.x, normal.y, normal.z);
  Vector3f sParallel(normal.y - normal.z, normal.z - normal.x,
                     normal.x - normal.y);
  sParallel.normalize();
  Vector3f sCross = sNormal.cross(sParallel);

  MatrixXf positions(3, 4);
  MatrixXf normals(3, 4);

  positions.col(0) << sPoint + 2 * (sCross + sParallel);
  positions.col(1) << sPoint + 2 * (sCross - sParallel);
  positions.col(2) << sPoint + 2 * (-sCross + sParallel);
  positions.col(3) << sPoint + 2 * (-sCross - sParallel);

  normals.col(0) << sNormal;
  normals.col(1) << sNormal;
  normals.col(2) << sNormal;
  normals.col(3) << sNormal;

  if (shader.uniform("u_color", false) != -1) {
    shader.setUniform("u_color", color);
  }
  shader.uploadAttrib("in_position", positions);
  if (shader.attrib("in_normal", false) != -1) {
    shader.uploadAttrib("in_normal", normals);
  }

  shader.drawArray(GL_TRIANGLE_STRIP, 0, 4);
}
