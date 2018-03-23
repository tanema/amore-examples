package game

import "github.com/tanema/amore/gfx"

// Sky is the sky!
type Sky struct {
	skybox *Voxel
	sun    *Voxel
}

func newSky() *Sky {
	return &Sky{
		skybox: newVoxel(0, 0, 0, 0, 0, 0.55, worldSaturation, 0.5, false),
		sun:    newVoxel(0, 0, 0, sunSize, sunSize, 0.16, 0, 1, false),
	}
}

func (sky *Sky) update(world *World) {
	sky.skybox.width, sky.skybox.height = gfx.GetWidth()/2, gfx.GetHeight()/2
	sky.skybox.update(world, float32(world.size), float32(world.size))
	if world.sin > 0 { // if daytime, sun is larger than moon
		sky.sun.width, sky.sun.height = sunSize, sunSize
		sky.sun.y = -1 * float32(world.camera.visible) * cos(world.timeOfDay)
		sky.sun.z = abs(world.sin * float32(world.camera.visible) * 3.2)
	} else {
		sky.sun.width, sky.sun.height = moonSize, moonSize
		sky.sun.y = -1 * -float32(world.camera.visible) * cos(world.timeOfDay)
		sky.sun.z = abs(world.sin * -float32(world.camera.visible) * 3.2)
	}
}

func (sky *Sky) draw(world *World) {
	sky.skybox.draw(world.camera, sky.skybox.x, sky.skybox.y)
	sky.sun.draw(world.camera, sky.sun.x, sky.sun.y)
}
