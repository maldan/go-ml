package ml_color

/*type ColorRGB[T constraints.Integer | constraints.Float] struct {
	R T
	G T
	B T
}

type ColorRGBA[T constraints.Integer | constraints.Float] struct {
	R T
	G T
	B T
	A T
}

func (c *ColorRGBA[T]) SetRGB(r T, g T, b T) {
	c.R = r
	c.G = g
	c.B = b
}

func (c ColorRGBA[T]) AddRGB(r T, g T, b T) ColorRGBA[T] {
	c.R += r
	c.G += g
	c.B += b
	return c
}

func (c ColorRGBA[T]) Avg(c2 ColorRGBA[T]) ColorRGBA[T] {
	c.R = T((mmath.Clamp((float32(c.R)+float32(c2.R))/2.0, 0, 255)))
	c.G = T((mmath.Clamp((float32(c.G)+float32(c2.G))/2.0, 0, 255)))
	c.B = T((mmath.Clamp((float32(c.B)+float32(c2.B))/2.0, 0, 255)))
	return c
}*/

/*func (c ColorRGBA[T]) Mix(c2 ColorRGBA[T], t float32) ColorRGBA[T] {
	if t > 1 {
		t = 1
	}
	if t < 0 {
		t = 0
	}
	c.R = T(mmath.Clamp(float32(c.R)*(1.0-t)+float32(c2.R)*t, 0, 255))
	c.G = T(mmath.Clamp(float32(c.G)*(1.0-t)+float32(c2.G)*t, 0, 255))
	c.B = T(mmath.Clamp(float32(c.B)*(1.0-t)+float32(c2.B)*t, 0, 255))
	return c
}*/

/*func (c ColorRGBA[T]) MulF32(v float32) ColorRGBA[T] {
	c.R = T(mmath.Clamp(float32(c.R)*v, 0, 255))
	c.G = T(mmath.Clamp(float32(c.G)*v, 0, 255))
	c.B = T(mmath.Clamp(float32(c.B)*v, 0, 255))
	return c
}

func (c ColorRGBA[T]) To01() ColorRGBA[float32] {
	c2 := ColorRGBA[float32]{}
	c2.R = float32(c.R) / 255.0
	c2.G = float32(c.G) / 255.0
	c2.B = float32(c.B) / 255.0
	return c2
}

func (c ColorRGBA[T]) To255() ColorRGBA[uint8] {
	c2 := ColorRGBA[uint8]{}
	c2.R = uint8(mmath.Clamp(float32(c.R)*255.0, 0, 255))
	c2.G = uint8(mmath.Clamp(float32(c.G)*255.0, 0, 255))
	c2.B = uint8(mmath.Clamp(float32(c.B)*255.0, 0, 255))
	return c2
}

func (c ColorRGBA[T]) Lerp(to ColorRGBA[T], t float32) ColorRGBA[T] {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return ColorRGBA[T]{
		A: T(mmath.Lerp(float32(c.A), float32(to.A), t)),
		R: T(mmath.Lerp(float32(c.R), float32(to.R), t)),
		G: T(mmath.Lerp(float32(c.G), float32(to.G), t)),
		B: T(mmath.Lerp(float32(c.B), float32(to.B), t)),
	}
}*/

/*func Lerp[T constraints.Integer | constraints.Float](from ColorRGBA[T], to ColorRGBA[T], t float32) ColorRGBA[T] {
	return ColorRGBA[T]{
		A: mmath.Lerp(from.A, to.A, t),
		R: mmath.Lerp(from.R, to.R, t),
		G: mmath.Lerp(from.G, to.G, t),
		B: mmath.Lerp(from.B, to.B, t),
	}
}
*/
