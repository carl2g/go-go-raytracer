package gogoraytracer

// type Plane struct {
// 	Origin     *Position3D
// 	NormalVect *Vector3D
// 	Material   *Material
// }

// func NewPlane(origin *Position3D, normalVect *Vector3D, Material *Material) (plane *Plane) {
// 	plane = &Plane{
// 		Origin:     origin,
// 		NormalVect: normalVect,
// 		Material:   Material,
// 	}
// 	return plane
// }

// func (plane *Plane) Intersect(ray *Ray) (intersectionDistance IntersectionDistance) {
// 	denum := plane.NormalVect.ScalarProduct(ray.Vector)

// 	if denum == 0 {
// 		return
// 	}

// 	num := plane.NormalVect.ScalarProduct(
// 		plane.Origin.Substract(ray.Origin),
// 	)
// 	d := num / denum

// 	if d <= 0 {
// 		return
// 	}

// 	return &d
// }

// func (plane *Plane) GetMaterial() *Material {
// 	return plane.Material
// }
