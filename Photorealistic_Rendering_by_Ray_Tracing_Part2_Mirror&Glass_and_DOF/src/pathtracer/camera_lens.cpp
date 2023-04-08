#include "camera.h"

#include <iostream>
#include <sstream>
#include <fstream>

#include "CGL/misc.h"
#include "CGL/vector2D.h"
#include "CGL/vector3D.h"

using std::cout;
using std::endl;
using std::max;
using std::min;
using std::ifstream;
using std::ofstream;

namespace CGL {

using Collada::CameraInfo;

Ray Camera::generate_ray_for_thin_lens(double x, double y, double rndR, double rndTheta) const {

  // TODO Project 3-2: Part 4
  // compute position and direction of ray from the input sensor sample coordinate.
  // Note: use rndR and rndTheta to uniformly sample a unit disk.

  //From piazza to run: set flags for radius = 0.05/0.1/0.25/1 and focal length = 4.2/4.4/4.6/4.8

  // pFilm is the point on the screen we're rendering
  // pLens is the point on the screen we uniformly sampled
  // pFocus is the point of the object in the world (or in the Plane of Focus)

  //1. Shoot ray like normal "Look up your code from Project 3-1, Part 1 to figure out the generated ray
  // direction (red segment in the figure)."
  Ray normalRay = generate_ray(x, y);

  //2. Uniformly sample the disk representing the thin lens (find pLens)
  Vector3D pLens = Vector3D(lensRadius*sqrt(rndR)*cos(rndTheta), lensRadius*sqrt(rndR)*sin(rndTheta), 0);

  //3.Calculate pFocus by intersecting the plane of focus with the red segment
  //saw something onn piazza saying:  I don't think you actually need to compute the intersection point with the
  // ray-plane equation - I just scaled by direction vector from 3.1 by the focalDistance
  // from piazza: Remember to subtract pLens from pFocus for your blue ray because the origin is no longer at (0, 0, 0). 
  Vector3D pFocus = normalRay.d * focalDistance - pLens;

  //4. Calculate the ray that originates from pLens, and set its direction towards pFocus (blue segment in the figure).
  Ray newRay = Ray(pLens, pFocus,fClip, 0); //(origin, direction, max_t, depth)
  newRay.min_t = nClip;
    
  //5. Normalize the direction of the ray, perform the camera-to-world conversion for both its origin and direction,
  // add pos to the ray's origin, and set the near and far clips to be the same as in Project 3-1, Part 1.
  newRay.d.normalize();
  newRay.d = c2w * newRay.d;
  newRay.o = c2w * newRay.o;
  newRay.o += pos;

  return newRay;
}


} // namespace CGL
