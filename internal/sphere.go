package gogoraytracer

import (
	"math"
)

type Sphere struct {
	Origin   *Position3D
	Radius   float64
	Material *Material
}

func NewSphere(origin *Position3D, radius float64, Material *Material) *Sphere {
	sphere := &Sphere{
		Origin:   origin,
		Radius:   radius,
		Material: Material,
	}
	return sphere
}

func (sphere *Sphere) Intersect(ray *Ray) IntersectionDistance {
	b := 2.0 * (ray.Vector.X*(ray.Origin.X-sphere.Origin.X) +
		ray.Vector.Y*(ray.Origin.Y-sphere.Origin.Y) +
		ray.Vector.Z*(ray.Origin.Z-sphere.Origin.Z))

	c := math.Pow(sphere.Origin.X, 2) +
		math.Pow(sphere.Origin.Y, 2) +
		math.Pow(sphere.Origin.Z, 2) +
		math.Pow(ray.Origin.X, 2) +
		math.Pow(ray.Origin.Y, 2) +
		math.Pow(ray.Origin.Z, 2) -
		2*ray.Origin.X*sphere.Origin.X -
		2*ray.Origin.Y*sphere.Origin.Y -
		2*ray.Origin.Z*sphere.Origin.Z -
		math.Pow(sphere.Radius, 2)

	tmp_det := math.Pow(b, 2) - 4*c

	if tmp_det < 0 {
		return nil
	}

	det := math.Sqrt(tmp_det)

	d1 := (-b + det) / 2
	d2 := (-b - det) / 2

	if d1 < 0 && d2 < 0 {
		return nil
	}

	if d1 >= 0 && d2 >= 0 {
		minDist := math.Min(d1, d2)
		return &minDist
	}
	minPositiveDist := math.Max(d1, d2)
	return &minPositiveDist
}

func (sphere *Sphere) GetMaterial() *Material {
	return sphere.Material
}

func (sphere *Sphere) NormalVector(position *Position3D) *Vector3D {
	return position.Substract(sphere.Origin).Normalized()
}
