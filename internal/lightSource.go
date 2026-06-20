package gogoraytracer

type LightSource struct {
	Origin    *Position3D
	Intensity float64
}

func NewLightSource(origin *Position3D, intensity float64) *LightSource {
	return &LightSource{
		Origin:    origin,
		Intensity: intensity,
	}
}
