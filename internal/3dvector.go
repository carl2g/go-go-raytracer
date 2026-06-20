package gogoraytracer

import (
	"math"
)

type Vector3D struct {
	X, Y, Z, Distance float64
}

func NewVector3D(x, y, z float64) *Vector3D {
	vector := &Vector3D{
		X: x,
		Y: y,
		Z: z,
	}
	vector.SetDistance()
	return vector
}

func (v *Vector3D) Add(other *Vector3D) *Vector3D {
	return NewVector3D(
		v.X+other.X,
		v.Y+other.Y,
		v.Z+other.Z,
	)
}

func (v *Vector3D) Substract(other *Vector3D) *Vector3D {
	return NewVector3D(
		v.X-other.X,
		v.Y-other.Y,
		v.Z-other.Z,
	)
}

func (v *Vector3D) SetDistance() {
	v.Distance = math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

func (v *Vector3D) Normalized() *Vector3D {
	return NewVector3D(
		v.X/v.Distance,
		v.Y/v.Distance,
		v.Z/v.Distance,
	)
}

func (v *Vector3D) CrossProduct(other *Vector3D) *Vector3D {
	return NewVector3D(
		v.Y*other.Z-v.Z*other.Y,
		v.Z*other.X-v.X*other.Z,
		v.Y*other.X-v.X*other.Y,
	)
}

func (v *Vector3D) Angle(other *Vector3D) (angle float64) {
	precision := math.Pow10(9)
	x := math.Trunc(
		v.ScalarProduct(other)*precision,
	) / (precision * (other.Distance * v.Distance))
	if x >= 1 {
		return 0
	} else if x <= -1 {
		return math.Pi
	}
	return math.Acos(x)
}

func (v *Vector3D) ScalarProduct(other *Vector3D) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vector3D) Scale(scalar float64) *Vector3D {
	return NewVector3D(
		v.X*scalar,
		v.Y*scalar,
		v.Z*scalar,
	)
}
