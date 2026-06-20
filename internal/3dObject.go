package gogoraytracer

type Material struct {
	Color             *Color
	RefractionIndex   float64
	ReflectionIndex   float64
	Shininess         float64
	TransparencyIndex float64
}

type Ray struct {
	Origin *Position3D
	Vector *Vector3D
}

func (ray *Ray) Position(scalar float64) *Position3D {
	scaledVect := ray.Vector.Scale(scalar)
	return &Position3D{
		X: ray.Origin.X + scaledVect.X,
		Y: ray.Origin.Y + scaledVect.Y,
		Z: ray.Origin.Z + scaledVect.Z,
	}
}

func NewRay(origin *Position3D, Vector *Vector3D) (ray *Ray) {
	ray = &Ray{
		Origin: origin,
		Vector: Vector.Normalized(),
	}
	return ray
}

type IntersectionDistance = *float64

type Object3D interface {
	Intersect(ray *Ray) (intersection IntersectionDistance)
	GetMaterial() *Material
	NormalVector(intersection *Position3D) *Vector3D
}
