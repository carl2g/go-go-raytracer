package gogoraytracer

type Plane struct {
	Origin     *Position3D
	NormalVect *Vector3D
	Material   *Material
}

func NewPlane(origin *Position3D, normalVect *Vector3D, Material *Material) *Plane {
	plane := &Plane{
		Origin:     origin,
		NormalVect: normalVect.Normalized(),
		Material:   Material,
	}
	return plane
}

func (plane *Plane) Intersect(ray *Ray) IntersectionDistance {
	denum := plane.NormalVect.ScalarProduct(ray.Vector)

	if denum == 0 {
		return nil
	}

	num := plane.NormalVect.ScalarProduct(
		plane.Origin.Substract(ray.Origin),
	)
	d := num / denum

	if d <= 0 {
		return nil
	}

	return &d
}

func (plane *Plane) GetMaterial() *Material {
	return plane.Material
}
func (plane *Plane) NormalVector(intersection *Position3D) *Vector3D {
	return plane.NormalVect
}
