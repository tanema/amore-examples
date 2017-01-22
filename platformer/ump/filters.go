package ump

func touchFilter(world *World, col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	return col.Touch.X, col.Touch.Y, []*Collision{}
}

func crossFilter(world *World, col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	return goalX, goalY, world.Project(body, goalX, goalY)
}

func slideFilter(world *World, col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	sx, sy := col.Touch.X, col.Touch.Y
	if col.Move.X != 0 || col.Move.Y != 0 {
		if col.Normal.X == 0 {
			sx = goalX
		} else {
			sy = goalY
		}
	}
	col.Data = Point{X: sx, Y: sy}
	body.x, body.y = col.Touch.X, col.Touch.Y
	return sx, sy, world.Project(body, sx, sy)
}

func bounceFilter(world *World, col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	tx, ty := col.Touch.X, col.Touch.Y
	bx, by := tx, ty
	if col.Move.X != 0 || col.Move.Y != 0 {
		bnx, bny := goalX-tx, goalY-ty
		if col.Normal.X == 0 {
			bny = -bny
		} else {
			bnx = -bnx
		}
		bx, by = tx+bnx, ty+bny
	}
	col.Data = Point{X: bx, Y: by}
	body.x, body.y = col.Touch.X, col.Touch.Y
	return bx, by, world.Project(body, bx, by)
}
