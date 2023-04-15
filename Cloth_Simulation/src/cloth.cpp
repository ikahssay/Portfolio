#include <iostream>
#include <math.h>
#include <random>
#include <vector>

#include "cloth.h"
#include "collision/plane.h"
#include "collision/sphere.h"

using namespace std;

Cloth::Cloth(double width, double height, int num_width_points,
             int num_height_points, float thickness) {
  this->width = width;
  this->height = height;
  this->num_width_points = num_width_points;
  this->num_height_points = num_height_points;
  this->thickness = thickness;

  buildGrid();
  buildClothMesh();
}

Cloth::~Cloth() {
  point_masses.clear();
  springs.clear();

  if (clothMesh) {
    delete clothMesh;
  }
}

void Cloth::buildGrid() {
  // TODO (Part 1): Build an evenly spaced grid of masses and springs.

    //There should be num_width_points by num_height_points total masses spanning width and height lengths, respectively.

    Vector3D position = Vector3D(0.0, 0.0, 0.0);

    double x = 0.0;
    double y = 0.0;
    double z = 0.0;

    for (int row = 0; row < num_height_points; row++) {
        for (int col = 0; col < num_width_points; col++) {
            double temp = (rand() * 1.0f) / RAND_MAX;
            double rand_offset = (temp * 2 - 1) / 1000; //generates random numbers from -1/1000 to 1/1000

            //If the cloth's orientation is HORIZONTAL, then set the y coordinate for all point masses to 1 while varying positions over the xz plane
            if (this->orientation == HORIZONTAL) {

                /*
                int total_masses = num_width_points * num_height_points;
                Vector3D start = Vector3D(0.0, 0.0, 0.0);
                Vector3D end = Vector3D(width, height, 0.0);
                Vector3D offset = (end - start) / (total_masses - 1);
                 */

                x = (double) (col * width) / (num_width_points - 1);
                z =  (double) (row * height) / (num_height_points - 1);
                position = Vector3D(x, 1.0, z);
            } else {
            //If the cloth's orientation is VERTICAL, then generate a small random offset between -1/1000 and 1/1000
            // for each point mass and use that as the z coordinate while varying positions on the xy plane
                //Use rand() function to produce the random offset

                /*
                int total_masses = num_width_points * num_height_points;
                Vector3D start = Vector3D(0.0, 0.0, 0.0);
                Vector3D end = Vector3D(width, height, 0.0);
                Vector3D offset = (end - start) / (total_masses - 1);
                 */

                x = (double) (col * width) / (num_width_points - 1);
                y = (double) (row * height) / (num_height_points - 1);
                position = Vector3D(x, y, rand_offset);
            }

            //If the point mass's (x,y) index is within the cloth's pinned vector (which stores the indices of the
            // pinned masses for each ROW), then set the point mass's PINNED boolean to TRUE
            // Masses corresponding to indices in pinned vector should have their pinned field set to true.
            bool isPinned = false;

            for (int i = 0; i < pinned.size(); i++) {

                //looking at col/row (of our point masses vector) in pinned
                vector<int> test = {col, row};
                if (pinned[i] == test) {
                    isPinned = true;
                }
            }

            PointMass m = PointMass(position, isPinned);
            point_masses.emplace_back(m);

        }

    }

    /*
     Next, create springs to apply the structual, shear, and bending constraints between point masses. Each of these
     springs takes in pointers to the two point masses that belong at the two ends of the spring and an enum that
     represents the type of constraint (STRUCTURAL, SHEARING, or BENDING).

     1) Structural constraints exist between a point mass and the point mass to its left as well as the point mass above it.
     2) Shearing constraints exist between a point mass and the point mass to its diagonal upper left as well as the point mass to its diagonal upper right.
     3) Bending constraints exist between a point mass and the point mass two away to its left as well as the point mass two above it.
     */
    //Initializing variables
    int index = 0;
    int above_pm = 0;
    int upper_left_pm = 0;
    int upper_right_pm = 0;
    int two_above_pm = 0;

    for (int row = 0; row < num_height_points; row++) {
        for (int col = 0; col < num_width_points; col++) {
            index =  (row * num_width_points) + col;

            //If the pm is the first of its column, there is nothing to the left of it
            if (col != 0) {
                Spring s = Spring(&point_masses[index], &point_masses[index - 1], STRUCTURAL);
                springs.emplace_back(s);
            }

            //If the pm is the first of its row, there is nothing above it
            if (row != 0) { //if index is NOT within the range of num_width_points <-> #cols in row
                above_pm = ((row-1) * num_width_points) + col;
                Spring s = Spring(&point_masses[index], &point_masses[above_pm], STRUCTURAL);
                springs.emplace_back(s);
            }

            //If the pm is the first of its column AND the first of its row, then there is nothing to the upper left of it
            if (col != 0 && row != 0) {
                upper_left_pm = ((row-1) * num_width_points) + (col - 1);
                Spring s = Spring(&point_masses[index], &point_masses[upper_left_pm], SHEARING);
                springs.emplace_back(s);
            }

            //If the pm is the LAST of its column AND the first of its row, there is nothing to the upper right of it
            if ( col != num_width_points - 1 && row != 0) {
                upper_right_pm = ((row-1) * num_width_points) + (col + 1);
                Spring s = Spring(&point_masses[index], &point_masses[upper_right_pm], SHEARING);
                springs.emplace_back(s);
            }

            //If the pm is the first of its column  OR the second of its column, there is nothing 2 to the left of it
            if (col > 1) {
                Spring s = Spring(&point_masses[index], &point_masses[index - 2], BENDING);
                springs.emplace_back(s);
            }

            //If the pm is the first of its row OR the second of its row, there is nothing 2 above it
            if (row > 1) {
                two_above_pm = ((row-2) * num_width_points) + col;
                Spring s = Spring(&point_masses[index], &point_masses[two_above_pm], BENDING);
                springs.emplace_back(s);
            }

        }

    }

}

void Cloth::simulate(double frames_per_sec, double simulation_steps, ClothParameters *cp,
                     vector<Vector3D> external_accelerations,
                     vector<CollisionObject *> *collision_objects) {
  double mass = width * height * cp->density / num_width_points / num_height_points;
  double delta_t = 1.0f / frames_per_sec / simulation_steps;
//  this->thickness = 0.01;



// TODO (Part 2): Compute total force acting on each point mass.
//  reset();
  Vector3D totalExForce = Vector3D(0.0, 0.0, 0.0);
    int i = 0;

  for (Vector3D& a : external_accelerations) {
      totalExForce += mass * a;
  }
  for (PointMass& pm : this->point_masses) {
      // Be sure to clear/reset forces at the start of each call to simulate so as to not accumulate extra,
      // nonexistant forces.
      pm.forces = Vector3D(0, 0, 0);

      // First, compute a total external force based on the external_accelerations and the mass
      // (recall that Newton's 2nd Law states that F = ma). Apply this external force to every point mass.
      pm.forces = totalExForce;
      i++;
  }

  Vector3D b_to_a = Vector3D(0.0, 0.0, 0.0);
  float displacement = 0.0;
  Vector3D spring_force = Vector3D(0.0, 0.0, 0.0);

  // Next, apply the spring correction forces.
  for (Spring& currentSpring : this->springs) {

      //For each spring, skip over the spring if that spring's constraint type is currently disabled.
      //Otherwise, compute the force applied to the two masses on its ends using Hooke's law.
      if ((currentSpring.spring_type == STRUCTURAL && cp->enable_structural_constraints) ||
          (currentSpring.spring_type == SHEARING && cp->enable_shearing_constraints) ||
          (currentSpring.spring_type == BENDING && cp->enable_bending_constraints)) {
          b_to_a = currentSpring.pm_b->position - currentSpring.pm_a->position;
          displacement = b_to_a.norm() - currentSpring.rest_length;
          spring_force =  cp->ks * b_to_a.unit() * displacement;

          //Because the bending constraint should be weaker that structural or shearing constraints,
          // you should multiply your ks by a small constant to achieve this. For example, 0.2 works well.
          if (currentSpring.spring_type == BENDING) {
              spring_force =  0.2 * cp->ks * b_to_a.unit() * displacement;
          }

      }

      //The force vector is the vector pointing from one point mass to the other with magnitude equal to ||F_s||
      // Apply this force to one point mass and an equal, but opposite force to the other.
      currentSpring.pm_a->forces += spring_force;
      currentSpring.pm_b->forces -= spring_force;

  }

  // TODO (Part 2): Use Verlet integration to compute new point mass positions
    //For each point mass, update the value in position according to the Verlet integration equations and store the previous time step's position in last_position
    for (PointMass& pm : this->point_masses) {
        if (!pm.pinned) {
            Vector3D temp_position = pm.position;
            Vector3D totalAcceleration = pm.forces / mass;

            //Note that the damping value is in units of percentage, so divide by 100 before subtracting from 1
            pm.position = pm.position + (1- (cp->damping/100))*(pm.position - pm.last_position) + totalAcceleration * (delta_t * delta_t); //Update position

            pm.last_position = temp_position; //Update last position
        }
    }

  // TODO (Part 4): Handle self-collisions.
  // As in Part 3, you will also need to update Cloth::simulate to account for potential self-collisions,
  // similarly accomplished by calling self_collide on each PointMass.

    build_spatial_map();

    for (PointMass& pm : this->point_masses) {
        self_collide(pm, simulation_steps);
    }


  // TODO (Part 3): Handle collisions with other primitives.
    //-> For every PointMass, you will want to try to collide it with every possible CollisionObject.
    for (PointMass& pm : this->point_masses) {
        for (CollisionObject *col_obj : *collision_objects){
            col_obj -> collide(pm);
        }
    }


  // TODO (Part 2): Constrain the changes to be such that the spring does not change
  // in length more than 10% per timestep [Provot 1995].

  //For each spring, apply this constraint by correcting the two point masses' positions such that the spring's length
  // is at most 10% greater than its rest_length at the end of any time step.
    for (Spring& currentSpring : this->springs) {
        //Do nothing if both are pinned
       if (!currentSpring.pm_a->pinned || !currentSpring.pm_b->pinned) {

            Vector3D b_to_a = currentSpring.pm_b->position - currentSpring.pm_a->position;
            double springLength = b_to_a.norm();

            if (springLength >= currentSpring.rest_length * 1.1 ) {
                // Maintain the same vector direction between the two point masses and only modify their distance apart
                // from each other.
                double differenceToMakeUp = springLength - ( 1.1 * currentSpring.rest_length);

                //Perform half of the correction to each point mass, unless one of them is pinned, in
                // which case apply the correction entirely to one point mass.
                if (!currentSpring.pm_a->pinned && !currentSpring.pm_b->pinned) {
                    currentSpring.pm_a->position += b_to_a.unit() * differenceToMakeUp/2;
                    currentSpring.pm_b->position -= b_to_a.unit() * differenceToMakeUp/2;
                }else if (!currentSpring.pm_a->pinned) {
                    currentSpring.pm_a->position += b_to_a.unit() * differenceToMakeUp;
                } else if (!currentSpring.pm_b->pinned) {
                    currentSpring.pm_b->position -= b_to_a.unit() * differenceToMakeUp;
                }
            }

        }
    }

}

void Cloth::build_spatial_map() {
  for (const auto &entry : map) {
    delete(entry.second);
  }
  map.clear();

  //cout << "building spacial map\n";
  // TODO (Part 4): Build a spatial map out of all of the point masses.
  // loop over all point masses and use the Cloth::hash_position method to populate the map as described above.
    for (PointMass& pm : this->point_masses) {
        float hashed_position = hash_position(pm.position);
        if (map[hashed_position] == NULL) {
            //create vector
            vector<PointMass *>* pmsVector = new vector<PointMass *>; //create new vector
            pmsVector -> push_back(&pm); //add to vector
            map[hashed_position] = pmsVector; //add vector to map
        } else {
            //vector already exists, so add to vector
            vector<PointMass *>* pmsVector = map[hashed_position]; //grab vector
            pmsVector -> push_back(&pm); //update vector
            map[hashed_position] = pmsVector; //update MAP's vector
        }

    }

}

void Cloth::self_collide(PointMass &pm, double simulation_steps) {
  // TODO (Part 4): Handle self-collision for a given point mass.
  // Cloth::self_collide takes in a point mass and looks up potential candidates for collision using the hash table.
  // For each pair between the point mass and a candidate point mass, determine whether they are within
  // 2 * thickness apart.
    float hashed_position = hash_position(pm.position);
    vector<PointMass *>* hashPmsVector = map[hashed_position];
    Vector3D sumCorrectionVectors = Vector3D(0.0, 0.0, 0.0);
    int counter = 0;
    for (PointMass *hashPm : *hashPmsVector) {
        // Make sure to not collide a point mass with itself!
        if (&pm == hashPm) {
            continue;
        } else if ((pm.position - hashPm->position).norm() < 2 * thickness){
            double pmDifference = (pm.position - hashPm->position).norm();
            // compute a correction vector that can be applied to the point mass (not the candidate one) such that
            // the pair would be 2 * thickness distance apart.
            double differenceToMakeUp = 2.0 * thickness - pmDifference;
            Vector3D correction_difference = (pm.position - hashPm->position).unit() * differenceToMakeUp;
            sumCorrectionVectors += correction_difference;
            counter++;
        }
    }

    if (counter == 0) {
        return;
    }

    // The final correction vector to the point mass's position
    // is the average of all of these pairwise correction vectors, scaled down by simulation_steps
    // saw on piazza we divide by simulation_steps
    Vector3D averageCorrectionVector = sumCorrectionVectors / counter / simulation_steps;
    pm.position =  pm.position + averageCorrectionVector;
  //update Cloth::simulate to account for potential self-collisions

}

float Cloth::hash_position(Vector3D pos) {
  // TODO (Part 4): Hash a 3D position into a unique float identifier that represents membership in some 3D box volume.

  //Effectively partition the 3D space into 3D boxes with dimensions w * h * t, where:
    float w = 3.0 * width / num_width_points;
    float h = 3.0 * height / num_height_points;
    float t = max(w, h);

  // Then, take the position and truncate its coordinates to the closest 3D box (hint: think modulo).
    float closestBoxX = pos.x - fmod(pos.x, w);
    float closestBoxY = pos.y - fmod(pos.y, h);
    float closestBoxZ = pos.z - fmod(pos.z, t);

  // Using these new coordinates, compute a unique number that corresponds to those 3D coordinates and return it.
  // This will be used as the unique key in our hash table. You may find the fmod function useful
  return (closestBoxX * 113.0 * 113.0) + (closestBoxY * 113.0) + closestBoxZ; //make each variable slightly different in case you have a position with switched values (i.e. you dont want position (2, 5, 7) to be in the same (hashed) bin as position (5, 2, 7).
}

///////////////////////////////////////////////////////
/// YOU DO NOT NEED TO REFER TO ANY CODE BELOW THIS ///
///////////////////////////////////////////////////////

void Cloth::reset() {
  PointMass *pm = &point_masses[0];
  for (int i = 0; i < point_masses.size(); i++) {
    pm->position = pm->start_position;
    pm->last_position = pm->start_position;
    pm++;
  }
}

void Cloth::buildClothMesh() {
  if (point_masses.size() == 0) return;

  ClothMesh *clothMesh = new ClothMesh();
  vector<Triangle *> triangles;

  // Create vector of triangles
  for (int y = 0; y < num_height_points - 1; y++) {
    for (int x = 0; x < num_width_points - 1; x++) {
      PointMass *pm = &point_masses[y * num_width_points + x];
      // Get neighboring point masses:
      /*                      *
       * pm_A -------- pm_B   *
       *             /        *
       *  |         /   |     *
       *  |        /    |     *
       *  |       /     |     *
       *  |      /      |     *
       *  |     /       |     *
       *  |    /        |     *
       *      /               *
       * pm_C -------- pm_D   *
       *                      *
       */

      float u_min = x;
      u_min /= num_width_points - 1;
      float u_max = x + 1;
      u_max /= num_width_points - 1;
      float v_min = y;
      v_min /= num_height_points - 1;
      float v_max = y + 1;
      v_max /= num_height_points - 1;

      PointMass *pm_A = pm                       ;
      PointMass *pm_B = pm                    + 1;
      PointMass *pm_C = pm + num_width_points    ;
      PointMass *pm_D = pm + num_width_points + 1;

      Vector3D uv_A = Vector3D(u_min, v_min, 0);
      Vector3D uv_B = Vector3D(u_max, v_min, 0);
      Vector3D uv_C = Vector3D(u_min, v_max, 0);
      Vector3D uv_D = Vector3D(u_max, v_max, 0);


      // Both triangles defined by vertices in counter-clockwise orientation
      triangles.push_back(new Triangle(pm_A, pm_C, pm_B,
                                       uv_A, uv_C, uv_B));
      triangles.push_back(new Triangle(pm_B, pm_C, pm_D,
                                       uv_B, uv_C, uv_D));
    }
  }

  // For each triangle in row-order, create 3 edges and 3 internal halfedges
  for (int i = 0; i < triangles.size(); i++) {
    Triangle *t = triangles[i];

    // Allocate new halfedges on heap
    Halfedge *h1 = new Halfedge();
    Halfedge *h2 = new Halfedge();
    Halfedge *h3 = new Halfedge();

    // Allocate new edges on heap
    Edge *e1 = new Edge();
    Edge *e2 = new Edge();
    Edge *e3 = new Edge();

    // Assign a halfedge pointer to the triangle
    t->halfedge = h1;

    // Assign halfedge pointers to point masses
    t->pm1->halfedge = h1;
    t->pm2->halfedge = h2;
    t->pm3->halfedge = h3;

    // Update all halfedge pointers
    h1->edge = e1;
    h1->next = h2;
    h1->pm = t->pm1;
    h1->triangle = t;

    h2->edge = e2;
    h2->next = h3;
    h2->pm = t->pm2;
    h2->triangle = t;

    h3->edge = e3;
    h3->next = h1;
    h3->pm = t->pm3;
    h3->triangle = t;
  }

  // Go back through the cloth mesh and link triangles together using halfedge
  // twin pointers

  // Convenient variables for math
  int num_height_tris = (num_height_points - 1) * 2;
  int num_width_tris = (num_width_points - 1) * 2;

  bool topLeft = true;
  for (int i = 0; i < triangles.size(); i++) {
    Triangle *t = triangles[i];

    if (topLeft) {
      // Get left triangle, if it exists
      if (i % num_width_tris != 0) { // Not a left-most triangle
        Triangle *temp = triangles[i - 1];
        t->pm1->halfedge->twin = temp->pm3->halfedge;
      } else {
        t->pm1->halfedge->twin = nullptr;
      }

      // Get triangle above, if it exists
      if (i >= num_width_tris) { // Not a top-most triangle
        Triangle *temp = triangles[i - num_width_tris + 1];
        t->pm3->halfedge->twin = temp->pm2->halfedge;
      } else {
        t->pm3->halfedge->twin = nullptr;
      }

      // Get triangle to bottom right; guaranteed to exist
      Triangle *temp = triangles[i + 1];
      t->pm2->halfedge->twin = temp->pm1->halfedge;
    } else {
      // Get right triangle, if it exists
      if (i % num_width_tris != num_width_tris - 1) { // Not a right-most triangle
        Triangle *temp = triangles[i + 1];
        t->pm3->halfedge->twin = temp->pm1->halfedge;
      } else {
        t->pm3->halfedge->twin = nullptr;
      }

      // Get triangle below, if it exists
      if (i + num_width_tris - 1 < 1.0f * num_width_tris * num_height_tris / 2.0f) { // Not a bottom-most triangle
        Triangle *temp = triangles[i + num_width_tris - 1];
        t->pm2->halfedge->twin = temp->pm3->halfedge;
      } else {
        t->pm2->halfedge->twin = nullptr;
      }

      // Get triangle to top left; guaranteed to exist
      Triangle *temp = triangles[i - 1];
      t->pm1->halfedge->twin = temp->pm2->halfedge;
    }

    topLeft = !topLeft;
  }

  clothMesh->triangles = triangles;
  this->clothMesh = clothMesh;
}
