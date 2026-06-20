package gogoraytracer

import (
	"math"
	"sync"
)

type Drawer interface {
	Draw()
	SetPixel(pixel *Pixel)
	getHeight() int
	getWidth() int
}

func RunEngine(drawer Drawer, pixelsChan chan PixelPosition) {
	var wg sync.WaitGroup

	pixelsToCompute(&wg, drawer, pixelsChan)

	camera := NewCamera(
		NewRay(
			&Position3D{X: -60.0, Y: -10.0, Z: 20.0},
			NewVector3D(3., 3.0, -1.),
		),
		drawer.getHeight(),
		drawer.getWidth(),
	)

	objectManager := presetObjectManager()
	for range 100 {
		wg.Go(func() {
			for pixelPosition := range pixelsChan {
				color := intersectionColor(
					NewRay(
						camera.MainRay.Origin,
						camera.PixelToVector[pixelPosition.Y][pixelPosition.X],
					),
					objectManager,
					0,
				)

				drawer.SetPixel(
					&Pixel{
						Color: color,
						PixelPosition: &PixelPosition{
							X: pixelPosition.X,
							Y: pixelPosition.Y,
						},
					},
				)
			}
		})

	}

	wg.Wait()
	drawer.Draw()
}

func objectIntersection(ray *Ray, objectManager *ObjectManager) (intersectedObject Object3D, closestIntersection IntersectionDistance) {
	for _, object := range objectManager.GetObjects() {
		intersectionDistance := object.Intersect(ray)
		if intersectionDistance == nil || *intersectionDistance < 0.00001 {
			continue
		}
		if closestIntersection == nil || *intersectionDistance < *closestIntersection {
			closestIntersection = intersectionDistance
			intersectedObject = object
		}
	}
	if closestIntersection == nil {
		return nil, nil
	}
	return intersectedObject, closestIntersection
}

func intersectionColor(ray *Ray, objectManager *ObjectManager, depth int) (color *Color) {
	var closestIntersection IntersectionDistance = nil

	intersectedObject, closestIntersection := objectIntersection(ray, objectManager)

	if intersectedObject == nil || closestIntersection == nil {
		return blackColor
	}

	objectMaterial := intersectedObject.GetMaterial()
	objColor := objectMaterial.Color
	scalar := *closestIntersection
	transparency := objectMaterial.TransparencyIndex
	reflection := objectMaterial.ReflectionIndex
	intersectionPosition := ray.Position(scalar)

	lightCoef := lightCoef(ray, intersectedObject, intersectionPosition, objectManager)
	reflectionColor, refractionColor := calculateReflectionRefaction(ray, intersectionPosition, intersectedObject, objectManager, depth)
	specularColor := specularColor(ray, intersectionPosition, intersectedObject, objectManager)

	surface := objColor.Scale(lightCoef)

	return surface.Scale(math.Max(1-transparency-reflection, 0)).
		Add(reflectionColor).
		Add(refractionColor).
		Add(specularColor)
}

func calculateReflectionRefaction(ray *Ray, intersectionPosition *Position3D, intersectedObject Object3D, objectManager *ObjectManager, depth int) (*Color, *Color) {
	objectMaterial := intersectedObject.GetMaterial()
	refractionColor := blackColor
	reflectionColor := blackColor

	if objectMaterial.ReflectionIndex == 0 && objectMaterial.RefractionIndex == 1 {
		return reflectionColor, refractionColor
	}

	if depth > 3 {
		return objectMaterial.Color, objectMaterial.Color
	}

	if objectMaterial.TransparencyIndex > 0 {
		refractionColor = calculateRefraction(ray, intersectionPosition, intersectedObject, objectManager, depth)
	}

	if objectMaterial.ReflectionIndex > 0 {
		reflectionColor = calculateReflection(ray, intersectionPosition, intersectedObject, objectManager, depth)
	}

	intersectionNormal := intersectedObject.NormalVector(intersectionPosition)
	cosAngleNomrmalRay := -ray.Vector.ScalarProduct(intersectionNormal)

	r := objectMaterial.ReflectionIndex

	if objectMaterial.TransparencyIndex > 0 {
		n1, n2 := 1.0, objectMaterial.RefractionIndex
		if cosAngleNomrmalRay < 0 {
			n1, n2 = n2, n1
		}

		cosθ := math.Abs(cosAngleNomrmalRay)
		r0 := math.Pow((n1-n2)/(n1+n2), 2)
		r = r0 + (1.0-r0)*math.Pow(1.0-cosθ, 5)
	}

	return reflectionColor.Scale(r), refractionColor.Scale(1.0 - r)
}

func calculateRefraction(ray *Ray, intersectionPosition *Position3D, intersectedObject Object3D, objectManager *ObjectManager, depth int) *Color {
	refractionIndexPrev := 1.0
	refractionIndexCurent := intersectedObject.GetMaterial().RefractionIndex

	intersectionNormal := intersectedObject.NormalVector(intersectionPosition)
	viewVector := ray.Vector.Normalized()
	cosAngleNomrmalRay := viewVector.ScalarProduct(intersectionNormal)

	if cosAngleNomrmalRay > 0 {
		refractionIndexPrev, refractionIndexCurent = refractionIndexCurent, refractionIndexPrev
		intersectionNormal = intersectionNormal.Scale(-1)
	} else {
		cosAngleNomrmalRay = -cosAngleNomrmalRay
	}

	sinAngleNomrmalRay := math.Sqrt(1 - math.Pow(cosAngleNomrmalRay, 2))

	t := refractionIndexPrev / refractionIndexCurent * sinAngleNomrmalRay

	g := math.Sqrt(1 - math.Pow(t, 2))

	parallelVecttoNormal := viewVector.Add(
		intersectionNormal.Scale(cosAngleNomrmalRay),
	).Scale(t / sinAngleNomrmalRay)

	perpendicularlVecttoNormal := intersectionNormal.Scale(-g)

	newRayVector := parallelVecttoNormal.Add(perpendicularlVecttoNormal).Normalized()
	refractionRay := &Ray{
		Vector: newRayVector,
		Origin: &Position3D{
			X: intersectionPosition.X + viewVector.X*0.000001,
			Y: intersectionPosition.Y + viewVector.Y*0.000001,
			Z: intersectionPosition.Z + viewVector.Z*0.000001,
		},
	}

	return intersectionColor(refractionRay, objectManager, depth+1)
}

func calculateReflection(ray *Ray, intersectionPosition *Position3D, intersectedObject Object3D, objectManager *ObjectManager, depth int) *Color {
	intersectionNormal := intersectedObject.NormalVector(intersectionPosition)
	angleIntersection := intersectionNormal.ScalarProduct(ray.Vector.Normalized())

	projectionVector := intersectionNormal.Scale(angleIntersection)
	reflexionVector := ray.Vector.Substract(projectionVector.Scale(2))

	reflexionRay := &Ray{
		Origin: &Position3D{
			X: intersectionPosition.X - ray.Vector.X*0.000001,
			Y: intersectionPosition.Y - ray.Vector.Y*0.000001,
			Z: intersectionPosition.Z - ray.Vector.Z*0.000001,
		},
		Vector: reflexionVector.Normalized(),
	}

	intersectedColor := intersectionColor(reflexionRay, objectManager, depth+1)
	return intersectedColor
}

func specularColor(ray *Ray, intersectionPosition *Position3D, intersectedObject Object3D, objectManager *ObjectManager) (specularColor *Color) {
	specularColor = blackColor.Copy()

	for _, lightSource := range objectManager.GetLightSources() {
		vectorToLightSource := lightSource.Origin.Substract(intersectionPosition)
		rayToLight := &Ray{
			Origin: intersectionPosition,
			Vector: vectorToLightSource.Normalized(),
		}

		lightIntersectedObject, intersectionDistance := objectIntersection(rayToLight, objectManager)

		if lightIntersectedObject != nil {
			scalar := *intersectionDistance
			lightRayIntersectionPosition := rayToLight.Position(scalar).Substract(intersectionPosition)

			if lightRayIntersectionPosition.Distance < vectorToLightSource.Distance {
				continue
			}
		}

		objectNormalVector := intersectedObject.NormalVector(intersectionPosition)
		halfWayVector := rayToLight.Vector.Substract(ray.Vector).Normalized()
		dotProd := math.Max(halfWayVector.ScalarProduct(objectNormalVector), 0)

		specularCoef := math.Pow(dotProd, intersectedObject.GetMaterial().Shininess)

		specularColor = specularColor.Add(
			whiteColor.Scale(lightSource.Intensity * specularCoef),
		)
	}

	return specularColor
}

func lightCoef(ray *Ray, intersectedObject Object3D, intersectionPosition *Position3D, objectManager *ObjectManager) (lightCoef float64) {
	lightCoef = 0.0
	const ambiantLight = 0.1

	for _, lightSource := range objectManager.GetLightSources() {
		vectorToLightSource := lightSource.Origin.Substract(intersectionPosition)
		rayToLight := &Ray{
			Origin: &Position3D{
				X: intersectionPosition.X - ray.Vector.X*0.00001,
				Y: intersectionPosition.Y - ray.Vector.Y*0.00001,
				Z: intersectionPosition.Z - ray.Vector.Z*0.00001,
			},
			Vector: vectorToLightSource.Normalized(),
		}

		normVect := intersectedObject.NormalVector(intersectionPosition)
		dotProd := normVect.ScalarProduct(rayToLight.Vector)
		tmpLightCoef := 1.0

		lightIntersectedObject, intersectionDistance := objectIntersection(rayToLight, objectManager)

		if lightIntersectedObject != nil {
			transparency := lightIntersectedObject.GetMaterial().TransparencyIndex
			scalar := *intersectionDistance
			lightRayIntersectionPosition := rayToLight.Position(scalar).Substract(intersectionPosition)

			if lightRayIntersectionPosition.Distance < vectorToLightSource.Distance {
				if transparency == 0 {
					continue
				} else {
					tmpLightCoef = tmpLightCoef * transparency
				}
			}
		}

		lightCoef += tmpLightCoef * lightSource.Intensity * math.Max(dotProd, 0)
	}

	return math.Min(lightCoef+ambiantLight, 1)
}

func presetObjectManager() (objectManger *ObjectManager) {
	objectManager := NewObjectManager()

	objectManager.AddLightSource(
		NewLightSource(
			&Position3D{
				-20, 50, 40,
			},
			0.8,
		),
	)
	objectManager.AddObject(
		NewSphere(
			&Position3D{
				X: 0, Y: 50, Z: 15,
			},
			10,
			&Material{
				Color: &Color{
					R: 255,
					G: 0,
					B: 0,
				},
				RefractionIndex:   1.,
				ReflectionIndex:   0.2,
				Shininess:         256,
				TransparencyIndex: 0.,
			},
		),
	)

	objectManager.AddObject(
		NewSphere(
			&Position3D{
				X: -10, Y: 50, Z: 20,
			},
			10,
			&Material{
				Color: &Color{
					R: 125,
					G: 0,
					B: 125,
				},
				RefractionIndex:   1.5,
				ReflectionIndex:   0,
				Shininess:         256,
				TransparencyIndex: 1.,
			},
		),
	)

	objectManager.AddObject(
		NewPlane(
			&Position3D{
				X: 0, Y: 0, Z: 0,
			},
			NewVector3D(
				0.0, 0.0, 1.0,
			),
			&Material{
				Color: &Color{
					R: 0,
					G: 255,
					B: 0,
				},
				RefractionIndex:   1.33,
				ReflectionIndex:   0.3,
				Shininess:         128,
				TransparencyIndex: 0.5,
			},
		),
	)

	objectManager.AddObject(
		NewPlane(
			&Position3D{
				X: 30, Y: 0, Z: 0,
			},
			NewVector3D(
				-1.0, 0.0, -0.4,
			),
			&Material{
				Color: &Color{
					R: 255,
					G: 125,
					B: 135,
				},
				RefractionIndex:   1.0,
				ReflectionIndex:   0.0,
				Shininess:         128,
				TransparencyIndex: 0.0,
			},
		),
	)
	return objectManager
}

func pixelsToCompute(wg *sync.WaitGroup, drawer Drawer, pixelsChan chan PixelPosition) {
	wg.Go(func() {
		for height := 0; height < drawer.getHeight(); height++ {
			for width := 0; width < drawer.getWidth(); width++ {
				pixelsChan <- PixelPosition{X: int32(width), Y: int32(height)}
			}
		}

		defer close(pixelsChan)
	})
}
