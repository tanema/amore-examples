package game

type gameObject interface {
	update(dt float32)
	tag() string
	destroy()
	push(strength float32)
	damage(intensity float32)
	draw(bool)
	updateOrder() int
}
