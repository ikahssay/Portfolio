#include "bvh.h"

#include "CGL/CGL.h"
#include "triangle.h"

#include <iostream>
#include <stack>

using namespace std;

namespace CGL {
namespace SceneObjects {

BVHAccel::BVHAccel(const std::vector<Primitive *> &_primitives,
                   size_t max_leaf_size) {

  primitives = std::vector<Primitive *>(_primitives);
  root = construct_bvh(primitives.begin(), primitives.end(), max_leaf_size);
}

BVHAccel::~BVHAccel() {
  if (root)
    delete root;
  primitives.clear();
}

BBox BVHAccel::get_bbox() const { return root->bb; }

void BVHAccel::draw(BVHNode *node, const Color &c, float alpha) const {
  if (node->isLeaf()) {
    for (auto p = node->start; p != node->end; p++) {
      (*p)->draw(c, alpha);
    }
  } else {
    draw(node->l, c, alpha);
    draw(node->r, c, alpha);
  }
}

void BVHAccel::drawOutline(BVHNode *node, const Color &c, float alpha) const {
  if (node->isLeaf()) {
    for (auto p = node->start; p != node->end; p++) {
      (*p)->drawOutline(c, alpha);
    }
  } else {
    drawOutline(node->l, c, alpha);
    drawOutline(node->r, c, alpha);
  }
}

BVHNode *BVHAccel::construct_bvh(std::vector<Primitive *>::iterator start,
                                 std::vector<Primitive *>::iterator end,
                                 size_t max_leaf_size) {

  // TODO (Part 2.1):
  // Construct a BVH from the given vector of primitives and maximum leaf
  // size configuration. The starter code build a BVH aggregate with a
  // single leaf node (which is also the root) that encloses all the
  // primitives.
    

//First, compute the bounding box of a list of primitives and initialize a new BVH Node with that bounding box.
  BBox bbox;
    

/*TODO: If you partition along the middle of the longest axis, you get a segfault error where there is a node with no primitives
    TODO: So instead, get the MIDDLE PRIMITIVE along the longest axis!
 */
//Finding the middle/averaged position of all primitives, and use that for partitioning (NOT ALONG THE MIDDLE OF THE AXIS!)
  int counter = 0;
  Vector3D average_positions = 0;
  for (auto p = start; p != end; p++)
  {
    average_positions = average_positions + (*p)->get_bbox().centroid();
    counter++;
      
      
    BBox bb = (*p)->get_bbox();
    bbox.expand(bb);
  }

  average_positions = average_positions / counter;

  BVHNode *node = new BVHNode(bbox);
   
  
 
//If no more than max_leaf_size primatives in list, the node we just created is a leaf node and we should update its start and end iterators appropriately.
    //Return this leaf node to end the recursion
  if (end - start <= max_leaf_size) {
      node -> l = NULL;
      node -> r = NULL;
      
      //TA recommended to initialize node's start and end iterators in the leaf's if statemement
      node->start = start;
      node->end = end;
      
      return node;
  } else {
    //Otherwise, we need to divide the primitives into a "left" and "right" collection
      //Remember: the primitives and their bounding boxes exist in 3D, so we want to split along the "best" axis -> the midpoint of the longest axis (like done in Discussion)
      //Compute the split point along the "best" axis and use it to divide all primitives along a "left" and "right" collection BASED ON THE CENTROID OF THEIR BOUNDING BOXES
      //Set the current node's left and right children by recursively callling construct_bvh(...)
      
      Vector3D axis_lengths = bbox.extent;
      double temp = max(axis_lengths.x, axis_lengths.y);
      double longest_axis = max(temp, axis_lengths.z);
      
      if (longest_axis == axis_lengths.x) {
          auto bound = std::partition(start, end, [&average_positions](Primitive * p){return p -> get_bbox().centroid().x < average_positions.x; });
          
          node->l = construct_bvh(start, bound, max_leaf_size);
          node->r = construct_bvh(bound, end, max_leaf_size);
          
      } else if (longest_axis == axis_lengths.y) { //Split along y axis, it's the longest...
          double midpoint = bbox.centroid().y;
          auto bound = std::partition(start, end, [&average_positions](Primitive * p){return p -> get_bbox().centroid().y < average_positions.y; });
          
          node -> l = construct_bvh(start, bound, max_leaf_size);
          node -> r = construct_bvh(bound, end, max_leaf_size);
          
      } else if (longest_axis == axis_lengths.z) { //Split along z axis, it's the longest...
          double midpoint = bbox.centroid().z;
          auto bound = std::partition(start, end, [&average_positions](Primitive * p){return p -> get_bbox().centroid().z < average_positions.z; });
          
          node -> l = construct_bvh(start, bound, max_leaf_size);
          node -> r = construct_bvh(bound, end, max_leaf_size);
      }

      return node;
  }

}

bool BVHAccel::has_intersection(const Ray &ray, BVHNode *node) const {
  // TODO (Part 2.3):
  // Fill in the intersect function.
  // Take note that this function has a short-circuit that the
  // Intersection version cannot, since it returns as soon as it finds
  // a hit, it doesn't actually have to find the closest hit.


/* Skeleton code
  for (auto p : primitives) {
    total_isects++;
    if (p->has_intersection(ray))
      return true;
  }
  return false;
  */
    //Create local variables of min_t and max_t and pass that in as an argument -> don't wanna change the ray!!
    double local_min_t = ray.min_t;
    double local_max_t = ray.max_t;
    if ((node -> bb).intersect(ray, local_min_t, local_max_t) ) {
        
        //If node is a leaf node, test intersection with all objects/primitives
        if (node -> isLeaf()) {
            for (auto p = node-> start; p < node -> end; p++) {
              total_isects++;
              if ((*p)->has_intersection(ray))
                return true;
            }
            
            return false;
        }
        
        //If ray hits the node's box, and its NOT a leaf, check the left and right branches recursively and look for intersections within those branches.
        if (has_intersection(ray, node -> r)) { //If it has an intersection on right, don't need to check the left
            return true;
        } else if (has_intersection(ray, node -> l)){
            return true;
        }
        
    } else {
        //If ray misses the node's bbox, return false
        return false;
    }
    
    
     
    return false; //If all else fails, return false
}

//TODO: REMEMBER TO CHECK BOTH LEFT AND RIGHT BRANCHES TO CHECK THE CLOSEST!
bool BVHAccel::intersect(const Ray &ray, Intersection *i, BVHNode *node) const {
  // TODO (Part 2.3):
  // Fill in the intersect function.

/* Skeleton code
  bool hit = false;
  for (auto p : primitives) {
    total_isects++;
    hit = p->intersect(ray, i) || hit;
  }
  return hit;
*/
    
    //Create local variables of min_t and max_t and pass that in as an argument -> don't wanna change the ray!!
    double local_min_t = ray.min_t;
    double local_max_t = ray.max_t;
    if ((node -> bb).intersect(ray, local_min_t, local_max_t) ) {
        
        //If node is a leaf node, test intersection with all objects
        if (node -> isLeaf()) {
            
            bool hit = false;
            for (auto p = node-> start; p < node -> end; p++) {
              total_isects++; //WHATS THE PURPOSE OF THIS?
                
            //All primitives update i and r.max_t correctly in their own intersection functions, so you don't need to worry about updating them in your BVH intersection functions
              hit = (*p)->intersect(ray, i) | hit;
                   
            }
            
            return hit;
        }
        
        //If ray hits the node's box, and its NOT a leaf, check the left and right branches recursively and look for intersections within those branches.
        bool left_hit = intersect(ray, i, node -> l);
        bool right_hit = intersect(ray, i, node -> r);
        
        return left_hit | right_hit;
        
    } else {
        //If ray misses the node's bbox, return false
        return false;
    }
     
    return false; //If all else fails, return false


}

} // namespace SceneObjects
} // namespace CGL
