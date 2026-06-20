package gogoraytracer

type Position3D struct {
	X, Y, Z float64
}

func (p *Position3D) Substract(other *Position3D) *Vector3D {
	return NewVector3D(
		p.X-other.X,
		p.Y-other.Y,
		p.Z-other.Z,
	)
}

func (p *Position3D) Add(other *Position3D) *Vector3D {
	return NewVector3D(
		p.X+other.X,
		p.Y+other.Y,
		p.Z+other.Z,
	)
}

func (p *Position3D) Translate(other *Vector3D) *Position3D {
	return &Position3D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}
