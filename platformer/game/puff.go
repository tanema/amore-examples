package game

type Puff struct {
	*Entity
	lived    float32
	lifeTime float32
}

func newPuff(gameMap *Map, x, y, vx, vy, minSize, maxSize float32) *Puff {
	puff := &Puff{}
	puff.Entity = newEntity(gameMap,
		puff, "puff",
		x, y,
		randRange(minSize, maxSize),
		randRange(minSize, maxSize),
	)
	puff.lifeTime = 0.1 + randMax(1)
	puff.vx, puff.vy = vx, vy
	return puff
}

func (puff *Puff) expand(dt float32) {
	cx, cy := puff.GetCenter()
	percent := puff.lived / puff.lifeTime
	if percent < 0.2 {
		puff.w = puff.w + (200+percent)*dt
		puff.h = puff.h + (200+percent)*dt
	} else {
		puff.w = puff.w + (20+percent)*dt
	}
	puff.l = cx - puff.w/2
	puff.t = cy - puff.h/2
}

func (puff *Puff) update(dt float32) {
	puff.lived = puff.lived + dt

	if puff.lived >= puff.lifeTime {
		puff.destroy()
	} else {
		puff.expand(dt)
		next_l, next_t := puff.l+puff.vx*dt, puff.t+puff.vy*dt
		puff.body.Update(next_l, next_t)
		puff.l, puff.t = next_l, next_t
	}
}

func (puff *Puff) draw(debug bool) {
	percent := min(1, (puff.lived/puff.lifeTime)*1.8)
	r, g, b := 255-floor(155*percent), 255-floor(155*percent), float32(100)
	l, t, w, h := puff.Extents()
	drawFilledRectangle(l, t, w, h, r, g, b)
}
