package gogoraytracer

type ObjectManager struct {
	Objects      []Object3D
	LightSources []*LightSource
}

func NewObjectManager() (manager *ObjectManager) {
	manager = &ObjectManager{}
	return manager
}

func (manager *ObjectManager) AddObject(object Object3D) {
	manager.Objects = append(manager.Objects, object)
}

func (manager *ObjectManager) removeObject(index int) {
	manager.Objects = append(manager.Objects[:index], manager.Objects[index+1:]...)
}

func (manager *ObjectManager) AddLightSource(lightSource *LightSource) {
	manager.LightSources = append(manager.LightSources, lightSource)
}

func (manager *ObjectManager) GetObjects() []Object3D {
	return manager.Objects
}

func (manager *ObjectManager) GetLightSources() []*LightSource {
	return manager.LightSources
}
