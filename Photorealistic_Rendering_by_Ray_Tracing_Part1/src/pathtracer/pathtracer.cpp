#include "pathtracer.h"

#include "scene/light.h"
#include "scene/sphere.h"
#include "scene/triangle.h"


using namespace CGL::SceneObjects;

namespace CGL {

PathTracer::PathTracer() {
  gridSampler = new UniformGridSampler2D();
  hemisphereSampler = new UniformHemisphereSampler3D();

  tm_gamma = 2.2f;
  tm_level = 1.0f;
  tm_key = 0.18;
  tm_wht = 5.0f;
}

PathTracer::~PathTracer() {
  delete gridSampler;
  delete hemisphereSampler;
}

void PathTracer::set_frame_size(size_t width, size_t height) {
  sampleBuffer.resize(width, height);
  sampleCountBuffer.resize(width * height);
}

void PathTracer::clear() {
  bvh = NULL;
  scene = NULL;
  camera = NULL;
  sampleBuffer.clear();
  sampleCountBuffer.clear();
  sampleBuffer.resize(0, 0);
  sampleCountBuffer.resize(0, 0);
}

void PathTracer::write_to_framebuffer(ImageBuffer &framebuffer, size_t x0,
                                      size_t y0, size_t x1, size_t y1) {
  sampleBuffer.toColor(framebuffer, x0, y0, x1, y1);
}

Vector3D
PathTracer::estimate_direct_lighting_hemisphere(const Ray &r,
                                                const Intersection &isect) {
  // Estimate the lighting from this intersection coming directly from a light.
  // For this function, sample uniformly in a hemisphere.

  // Note: When comparing Cornel Box (CBxxx.dae) results to importance sampling, you may find the "glow" around the light source is gone.
  // This is totally fine: the area lights in importance sampling has directionality, however in hemisphere sampling we don't model this behaviour.

  // make a coordinate system for a hit point
  // with N aligned with the Z direction.
  Matrix3x3 o2w;
  make_coord_space(o2w, isect.n);
  Matrix3x3 w2o = o2w.T();

  // w_out points towards the source of the ray (e.g.,
  // toward the camera if this is a primary ray)
  const Vector3D hit_p = r.o + r.d * isect.t;
  const Vector3D w_out = w2o * (-r.d);

  // This is the same number of total samples as
  // estimate_direct_lighting_importance (outside of delta lights). We keep the
  // same number of samples for clarity of comparison.
  int num_samples = scene->lights.size() * ns_area_light;
  Vector3D L_out;

  // TODO (Part 3): Write your sampling loop here
// <<<<<<< Lucy's
    Vector3D sumIncomingLight = Vector3D(0, 0, 0);

    for (int i = 0; i < num_samples; i++) {
        Vector3D hemisphereSample = hemisphereSampler->get_sample(); //object space
        //convert sample to world space
        hemisphereSample = o2w * hemisphereSample;
        //generate ray with constructor, setting origin as hit_p (which is in world space)
        Ray generatedRay(hit_p, hemisphereSample, r.max_t, 0);
        generatedRay.min_t = EPS_F;

        Intersection newIntersection;
        if (bvh->intersect(generatedRay, &newIntersection)) {
            //if it intersects, use the equation
            sumIncomingLight += isect.bsdf->f(w_out, hemisphereSample) * newIntersection.bsdf->get_emission() * 2 * PI *
                                dot(isect.n.unit(), hemisphereSample.unit());
// ======= IMAN's
//     //1. Estimate how much light arrived at an intersection point
//         //Use Monte Carlo Estimator to integrate over all light arriving in a hemisphere around the point of interest, hit_p (i.e. Implement Irradiance)
//
//         Vector3D sampled_direction_from_hemisphere;
//         Vector3D sum_of_samples = Vector3D(0.0, 0.0, 0.0);
//         double pdf = 1/ (2 * PI);
//         Vector3D bsdf_ratio = Vector3D(0.0, 0.0, 0.0);
//
//         for (int i = 0; i < num_samples; i++) {
//         //A. Uniformly sample directions from the hemisphere
//             //i. Sample new w_i (object coordinate)
//             sampled_direction_from_hemisphere = hemisphereSampler->get_sample(); //Returns an object coordinate
//             sampled_direction_from_hemisphere = o2w * sampled_direction_from_hemisphere;
//
//
//         //B. Check if a NEWLY GENERATED ray going FROM hit_p in the SAMPLED DIRECTION intrersects with a light source, and use that in Monte Carlo Estimator
//             //--> generate a ray using the world version of w_j
//             //Assumed camera targetPosition is hit_p (i.e. assume ray's origin is hit_p) TODO: RIGHT?
//             //const Ray generatedRay = camera->generate_ray( (float) sampled_direction_from_hemisphere.x/sampleBuffer.w, (float) sampled_direction_from_hemisphere.y/sampleBuffer.h);
//             Ray generatedRay = Ray(hit_p, sampled_direction_from_hemisphere, 0); //(origin, direction, depth)
//             generatedRay.min_t = EPS_F;
//
//
//             //iii. Use the intersect test to see what w_i originated from before hitting hit_p (i.e. see if rays match up)
//             Intersection newisect;
//             if ( bvh -> intersect(generatedRay, &newisect)) {
//             //If it intersects at point `new_isect`, use the BSDF's emission at `new_isect`, the BSDF's ratio f at `isect` and the angle of the sampled ray from the surface normal to compute the contribution to add to L_out using the inside of the summation in the spec:
//                 //---> Recall that the cosine of the angle between two unit vectors is equal to their dot product.
//                 //L_out = est_radiance_global_illumination(generatedRay);
//                 L_out = newisect.bsdf -> get_emission();
//                 bsdf_ratio = isect.bsdf -> f(w_out, sampled_direction_from_hemisphere);
//                 double cos_angle = dot(newisect.n.unit(), sampled_direction_from_hemisphere.unit()) ;
//
//                 sum_of_samples += (L_out * bsdf_ratio * cos_angle * pdf);
//             }

// >>>>>>> 9bd82462ad7e2e849259ca2f41550dbc413177d2
        }
    }
    return sumIncomingLight/(double) num_samples;
}

Vector3D
PathTracer::estimate_direct_lighting_importance(const Ray &r,
                                                const Intersection &isect) {
  // Estimate the lighting from this intersection coming directly from a light.
  // To implement importance sampling, sample only from lights, not uniformly in
  // a hemisphere.

  // make a coordinate system for a hit point
  // with N aligned with the Z direction.
    Matrix3x3 o2w;
    make_coord_space(o2w, isect.n);
    Matrix3x3 w2o = o2w.T();

  // w_out points towards the source of the ray (e.g.,
  // toward the camera if this is a primary ray)

  const Vector3D hit_p = r.o + r.d * isect.t;
  const Vector3D w_out = w2o * (-r.d);
  Vector3D L_out;

  double numSamples = 0; //numSamples = counter
  Vector3D bsdf_ratio = Vector3D(0, 0, 0);
  Vector3D sumIncomingLight = Vector3D(0, 0, 0);

    //loop through all lights
    for (SceneLight* light : scene->lights) {
        //if light is a point light -> sample once, else sample ns_area_light number of times
        int numLoops = light->is_delta_light() ? 1 : ns_area_light;

        for (int i = 0; i < numLoops; i++) {
            Vector3D sampledDirection;
            double dist_to_light;
            double pdf;

        //1. Sample directions between the light source and the hit_p
            /*--> sample_L RETURNS the emmitted radiance (i.e. L-out), and writes:
                    - A world space unit wi vector giving the sampled direction between p and the light source
                    - A ditToLight double giving the distance btwn p and the light source in the wi direction
                    - A pdf double giving the value of the probability density function evaluated at the wi direction

             */
            Vector3D L_out = light->sample_L(hit_p, &sampledDirection, &dist_to_light, &pdf); //L_out = emitted radiance

            //same as task 3, with max_t a little different

            Ray generatedRay = Ray(hit_p, sampledDirection, dist_to_light - EPS_F, 0);
            generatedRay.min_t = EPS_F;
            Intersection newIntersection;

            if (!bvh->intersect(generatedRay, &newIntersection)) {
             //plugging emittedRadiance and pdf in
                Vector3D bsdf_ratio = isect.bsdf->f(w_out, sampledDirection);
                double cos_angle = dot(isect.n.unit(), (sampledDirection).unit());

                sumIncomingLight += (L_out * bsdf_ratio * cos_angle * (1.0 / pdf));
            }

            numSamples++;
        }
    }
    Vector3D avgIncomingLight = sumIncomingLight/(double)numSamples;
    return avgIncomingLight;
}

Vector3D PathTracer::zero_bounce_radiance(const Ray &r,
                                          const Intersection &isect) {
  // TODO: Part 3, Task 2

  // Returns the light that results from no bounces of light
   return isect.bsdf -> get_emission();
}

Vector3D PathTracer::one_bounce_radiance(const Ray &r,
                                         const Intersection &isect) {
  // TODO: Part 3, Task 3
  // Returns either the direct illumination by hemisphere or importance sampling
  // depending on `direct_hemisphere_sample`
//  return estimate_direct_lighting_hemisphere(r, isect);
    return direct_hemisphere_sample ? estimate_direct_lighting_hemisphere(r, isect) : estimate_direct_lighting_importance(r, isect);
}

Vector3D PathTracer::at_least_one_bounce_radiance(const Ray &r,
                                                  const Intersection &isect) {
  Matrix3x3 o2w;
  make_coord_space(o2w, isect.n);
  Matrix3x3 w2o = o2w.T();

  Vector3D hit_p = r.o + r.d * isect.t;
  Vector3D w_out = w2o * (-r.d);

  Vector3D L_out(0, 0, 0);

  // TODO: Part 4, Task 2
  // Returns the one bounce radiance + radiance from extra bounces at this point.
  // Should be called recursively to simulate extra bounces.

  // end recursion when depth = 1, or on coin flip probability as long as indirect has been run
  // once (the r.depth < max_ray_depth makes sure depth has been decremented once)
  if (r.depth == 0) {
      return L_out;
  }
  if (r.depth == 1 || (r.depth < max_ray_depth && coin_flip(0.3))) {
      return one_bounce_radiance(r, isect);
        //for write up: only indirect lighting
//        return Vector3D(0, 0, 0);
  } else {
      L_out = one_bounce_radiance(r, isect);
      Vector3D w_in;
      double pdf;
      Vector3D fValue = isect.bsdf->sample_f(w_out, &w_in, &pdf);
      //convert to world frame
      w_in = o2w * w_in;
      //decrement depth
      Ray newRay = Ray(hit_p, w_in, r.max_t, r.depth - 1);
      newRay.min_t = EPS_F;
      Intersection newP;
      if (bvh->intersect(newRay, &newP)) {
          //from piazza: only use the cosine term if it is positive (doesn't make a big difference in code)
//          w_in = -1.0 * w_in;
//          if (dot(isect.n.unit(), (w_in).unit()) > 0) {
              L_out += at_least_one_bounce_radiance(newRay, newP) * fValue *
                       (1.0 / pdf) * (1.0/0.7) *
                       dot(isect.n.unit(), (w_in).unit());
//          }

      }
  }
  return L_out;
}

Vector3D PathTracer::est_radiance_global_illumination(const Ray &r) {
  Intersection isect;
  Vector3D L_out;

  // You will extend this in assignment 3-2.
  // If no intersection occurs, we simply return black.
  // This changes if you implement hemispherical lighting for extra credit.

  // The following line of code returns a debug color depending
  // on whether ray intersection with triangles or spheres has
  // been implemented.
  //TODO: UNCOMMENT THIS LINE WHEN RENDERING FOR WRITEUP!!
  // REMOVE THIS LINE when you are ready to begin Part 3.
    if (!bvh->intersect(r, &isect)) {
        return envLight ? envLight->sample_dir(r) : L_out;
    } else {
        return zero_bounce_radiance(r, isect) + at_least_one_bounce_radiance(r, isect);
//        return at_least_one_bounce_radiance(r, isect);
    }


  //L_out = (isect.t == INF_D) ? debug_shading(r.d) : normal_shading(isect.n);
  //return L_out;

  // TODO (Part 3): Return the direct illumination.


  // TODO (Part 4): Accumulate the "direct" and "indirect"
  // parts of global illumination into L_out rather than just direct


}

void PathTracer::raytrace_pixel(size_t x, size_t y) {
  // TODO (Part 1.2):
  // Make a loop that generates num_samples camera rays and traces them
  // through the scene. Return the average Vector3D.
  // You should call est_radiance_global_illumination in this function.

  int num_samples = ns_aa;          // total samples to evaluate
  Vector2D origin = Vector2D(x, y); // bottom left corner of the pixel
  Vector3D sumRadiance = Vector3D(0, 0, 0);

    //part 1:
    //  for (int i = 0; i < num_samples; i++) {
    //      Vector2D randomSample = gridSampler->get_sample();
    //      Ray randomRay = camera->generate_ray((float) ((origin.x + randomSample.x) / (float) sampleBuffer.w), (float) ((origin.y + randomSample.y) / (float) sampleBuffer.h));
    //      //part 4, initialize your camera rays' depths as max_ray_depth in raytrace_pixel
    //      randomRay.depth = max_ray_depth;
    //      Vector3D estimatedRadiance = est_radiance_global_illumination(randomRay);
    //      sumRadiance += estimatedRadiance;
    //  }
    //  Vector3D avgRadiance = sumRadiance / num_samples;
    //  sampleBuffer.update_pixel(avgRadiance, x, y);

  // TODO (Part 5):
  // Modify your implementation to include adaptive sampling.
  // Use the command line parameters "samplesPerBatch" and "maxTolerance"
//    cout << "num samples: " << num_samples << " samples per batch: " << samplesPerBatch << "\n";
    float s1 = 0; //sum of sample's illuminances
    float s2 = 0; //sum of sample's illuminances squared
    int num_samples_to_converge = num_samples;
    for (int i = 0; i < num_samples; i++) {
        Vector2D randomSample = gridSampler->get_sample();
        Ray randomRay = camera->generate_ray((float) ((origin.x + randomSample.x) / (float) sampleBuffer.w), (float) ((origin.y + randomSample.y) / (float) sampleBuffer.h));
        //part 4, initialize your camera rays' depths as max_ray_depth in raytrace_pixel
        randomRay.depth = max_ray_depth;
        Vector3D estimatedRadiance = est_radiance_global_illumination(randomRay);
        sumRadiance += estimatedRadiance;
        s1 += estimatedRadiance.illum();
        s2 += estimatedRadiance.illum()*estimatedRadiance.illum();

        //check convergence every samples per batch
        if (i > 0 && i % samplesPerBatch == 0) {
            float mu = s1/(float)i; //mean
            float variance = 1.0/((float)i - 1.0) * (s2 - (s1*s1)/(float)i); //sigma squared
            float sigma = sqrt(variance);

            float convergence_I = 1.96 * (sigma / sqrt(i));
//            cout << "mu: " << mu << endl;
//            cout << "sigma: " << sigma << endl;
//            cout << "I: " << convergence_I << endl;
//            cout << "threshold: " << maxTolerance * mu << "\n\n";

            //check if converged
            if (convergence_I <= maxTolerance * mu) {
                //set the number of samples we used to converge
                num_samples_to_converge = i;
                break;
            }
        }
    }

    Vector3D avgRadiance = sumRadiance / num_samples_to_converge;
    sampleBuffer.update_pixel(avgRadiance, x, y);
    sampleCountBuffer[x + y * sampleBuffer.w] = num_samples_to_converge;
}

void PathTracer::autofocus(Vector2D loc) {
  Ray r = camera->generate_ray(loc.x / sampleBuffer.w, loc.y / sampleBuffer.h);
  Intersection isect;

  bvh->intersect(r, &isect);

  camera->focalDistance = isect.t;
}

} // namespace CGL
