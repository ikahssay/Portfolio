#include "bbox.h"

#include "GL/glew.h"

#include <algorithm>
#include <iostream>

namespace CGL {

/**
 * Ray - bbox intersection.
 * Intersects ray with bounding box, does not store shading information.
 * \param r the ray to intersect with
 * \param t0 lower bound of intersection time
 * \param t1 upper bound of intersection time
 */
bool BBox::intersect(const Ray& r, double& t0, double& t1) const {

  // TODO (Part 2.2):
  // Implement ray - bounding box intersection test
  // If the ray intersected the bouding box within the range given by
  // t0, t1, update t0 and t1 with the new intersection times.
    //std::cout << min << endl;
    
    Vector3D p0 = this -> min;
    Vector3D p1 = this -> max;
    
    //Find tx1, ty1, and tz1 using formula: (p - o )/ d
    Vector3D t_0 = (p0 - (r.o)); //vectors subtract by elements
    double tx0 = t_0.x / r.d.x;
    double ty0 = t_0.y / r.d.y;
    double tz0 = t_0.z / r.d.z;
    
    //Find tx1, ty1, and tz1 using formula: (p - o )/ d
    Vector3D t_1 = (p1 - (r.o)); //vectors subtract by elements
    double tx1 = t_1.x / r.d.x;
    double ty1 = t_1.y / r.d.y;
    double tz1 = t_1.z / r.d.z;
    
    //Calculate t_mins
    double t_x_mins = std::min(tx0, tx1);
    double t_y_mins = std::min(ty0, ty1);
    double t_z_mins = std::min(tz0, tz1);
    
    //Calculate t_maxs
    double t_x_maxs = std::max(tx0, tx1);
    double t_y_maxs = std::max(ty0, ty1);
    double t_z_maxs = std::max(tz0, tz1);
    
    //Calculate t_entry -> t_entry = max(t_mins)
    double temp = std::max(t_x_mins, t_y_mins);
    double t_entry = std::max(temp, t_z_mins);
    
    //Calculate t_exit -> t_exit = min(t_maxs)
    temp = std::min(t_x_maxs, t_y_maxs);
    double t_exit = std::min(temp, t_z_maxs);
    
    //Check if there is a valid collision
        //-> If the exit is greater than the entry, then there is an intersection/collision
    if (t_exit >= t_entry) {
        t0 = t_entry;
        t1 = t_exit;
        return true;
    }
   
    
    return false;   
}

void BBox::draw(Color c, float alpha) const {

  glColor4f(c.r, c.g, c.b, alpha);

  // top
  glBegin(GL_LINE_STRIP);
  glVertex3d(max.x, max.y, max.z);
  glVertex3d(max.x, max.y, min.z);
  glVertex3d(min.x, max.y, min.z);
  glVertex3d(min.x, max.y, max.z);
  glVertex3d(max.x, max.y, max.z);
  glEnd();

  // bottom
  glBegin(GL_LINE_STRIP);
  glVertex3d(min.x, min.y, min.z);
  glVertex3d(min.x, min.y, max.z);
  glVertex3d(max.x, min.y, max.z);
  glVertex3d(max.x, min.y, min.z);
  glVertex3d(min.x, min.y, min.z);
  glEnd();

  // side
  glBegin(GL_LINES);
  glVertex3d(max.x, max.y, max.z);
  glVertex3d(max.x, min.y, max.z);
  glVertex3d(max.x, max.y, min.z);
  glVertex3d(max.x, min.y, min.z);
  glVertex3d(min.x, max.y, min.z);
  glVertex3d(min.x, min.y, min.z);
  glVertex3d(min.x, max.y, max.z);
  glVertex3d(min.x, min.y, max.z);
  glEnd();

}

std::ostream& operator<<(std::ostream& os, const BBox& b) {
  return os << "BBOX(" << b.min << ", " << b.max << ")";
}

} // namespace CGL
