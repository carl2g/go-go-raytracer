package gogoraytracer

import (
	"math"
)

type Camera struct {
	MainRay       *Ray
	PixelToVector [][]*Vector3D
	widthAngle    float64
	heightAngle   float64
}

func NewCamera(ray *Ray, height, width int) (camera *Camera) {
	camera = &Camera{
		MainRay:     ray,
		widthAngle:  math.Pi / 2.0,
		heightAngle: (math.Pi / 2.0) * (float64(height) / float64(width)),
	}
	camera.PixelToVector = make([][]*Vector3D, height)
	for i := range camera.PixelToVector {
		camera.PixelToVector[i] = make([]*Vector3D, width)
	}
	camera.setPixelsToVectorsEquiDistancePixels(height, width)
	// camera.setPixelsToVectorsEquiAngles(height, width)
	return camera
}

func (camera *Camera) setPixelsToVectorsEquiDistancePixels(height, width int) {
	forward := camera.MainRay.Vector.Normalized()

	centerAngleX := camera.widthAngle / 2.0
	halfWidth := forward.Y * math.Tan(centerAngleX)
	halfHeight := halfWidth * float64(height) / float64(width)
	widhtStep := 2 * halfWidth / float64(width)
	heightStep := 2 * halfHeight / float64(height)
	halfWidhtStep := widhtStep / 2
	halfHeightStep := heightStep / 2

	worldUp := NewVector3D(0, 0, 1)
	right := forward.CrossProduct(worldUp).Normalized()
	up := right.CrossProduct(forward).Normalized()

	for y := range height {
		zAxis := halfHeight - (heightStep*float64(y) + halfHeightStep)
		for x := range width {
			xAxis := (widhtStep*float64(x) + halfWidhtStep) - halfWidth
			camera.PixelToVector[y][x] = forward.
				Add(right.Scale(xAxis)).
				Substract(up.Scale(zAxis)).
				Normalized()
		}
	}
}

func (camera *Camera) setPixelsToVectorsEquiAngles(height, width int) {
	centerVector := camera.MainRay.Vector.Normalized()
	centerAngleX := camera.widthAngle / 2.0
	centerAngleY := camera.heightAngle / 2.0
	widthAngleStep := camera.widthAngle / float64(width)
	heightAngleStep := camera.heightAngle / float64(height)

	for y := range height {
		zPixel := heightAngleStep*float64(y) + heightAngleStep/2
		zAxis := centerVector.Y * math.Tan(centerAngleY-zPixel)
		for x := range width {
			xPixel := widthAngleStep*float64(x) + widthAngleStep/2
			xAxis := centerVector.Y * math.Tan(xPixel-centerAngleX)
			camera.PixelToVector[y][x] = NewVector3D(
				centerVector.X+xAxis,
				centerVector.Y,
				centerVector.Z+zAxis,
			).Normalized()
		}
	}
}
