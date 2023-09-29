package mcgi_rasterizer

import (
	mcgi_mesh "github.com/maldan/go-ml/cgi/mesh"
	mmath "github.com/maldan/go-ml/math"
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

func drawLine(
	mesh *mcgi_mesh.Mesh,
	t mmath_geom.Triangle4D[float32],
	square float32,
	screen *screen,
	x1, y, x2 int,
	triangleId int,
) {
	x1 = mmath.Clamp(x1, 0, screen.Width)
	x2 = mmath.Clamp(x2, 0, screen.Width)
	width := mmath.Abs(x2 - x1)
	if width == 0 {
		return
	}
	if width > screen.Width-1 {
		width = screen.Width - 1
	}

	uvTriangle := mesh.UVTriangleList[triangleId]

	for i := 0; i < width; i++ {
		p := mmath_la.Vector2[float32]{X: float32(x1 + i), Y: float32(y)}

		a := t.A
		b := t.B
		c := t.C

		s1 := 0.5 * ((b.X-p.X)*(c.Y-p.Y) - (c.X-p.X)*(b.Y-p.Y))
		s2 := 0.5 * ((p.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(p.Y-a.Y))
		s3 := 0.5 * ((b.X-a.X)*(p.Y-a.Y) - (p.X-a.X)*(b.Y-a.Y))

		alpha := s1 / square
		beta := s2 / square
		gamma := s3 / square

		// p = oldP
		if p.X >= 0 && p.Y >= 0 && int(p.X) < screen.Width && int(p.Y) < screen.Height {
			// fn(p, alpha, beta, gamma)
			depthA := t.A.W
			depthB := t.B.W
			depthC := t.C.W

			wA := 1.0 / depthA
			wB := 1.0 / depthB
			wC := 1.0 / depthC

			wSum := alpha*wA + beta*wB + gamma*wC

			alphaPerspective := (alpha * wA) / wSum
			betaPerspective := (beta * wB) / wSum
			gammaPerspective := (gamma * wC) / wSum

			intX := alpha*t.A.X + beta*t.B.X + gamma*t.C.X
			intY := alpha*t.A.Y + beta*t.B.Y + gamma*t.C.Y
			intZ := alpha*t.A.Z + beta*t.B.Z + gamma*t.C.Z
			intW := alpha*t.A.W + beta*t.B.W + gamma*t.C.W

			xClip := intX / intW
			yClip := intY / intW
			zClip := intZ / intW

			if xClip <= -1 || xClip >= 1 {
				return
			}
			if yClip <= -1 || yClip >= 1 {
				return
			}
			if zClip <= -1 || zClip >= 1 {
				return
			}

			depth := mmath.Remap(zClip, -1, 1, 1, 0)

			// Check zBuffer
			if screen.GetZPixel(mmath.FloorInt(p.X), mmath.FloorInt(p.Y)) > depth {
				return
			}

			// UV
			itrUv := uvTriangle.InterpolateABG(alphaPerspective, betaPerspective, gammaPerspective)

			itrUv.X = mmath.Mod(itrUv.X, 1)
			itrUv.Y = mmath.Mod(itrUv.Y, 1)

			// Read diffuse color
			texelCurrent := mmath_la.Vector2[float32]{
				X: itrUv.X * float32(mesh.Texture.Diffuse.Width-1),
				Y: itrUv.Y * float32(mesh.Texture.Diffuse.Height-1),
			}

			// fmt.Printf("%v\n", texelNext.Sub(texelCurrent))
			texelColor := mesh.Texture.Diffuse.GetPixel(mmath.RoundInt(texelCurrent.X), mmath.RoundInt(texelCurrent.Y))
			// color := texelColor.MulF32(lightPower)
			color := texelColor //.MulScalar(lightPower)
			// color = colorX
			//color := texelColor

			screen.SetPixel(mmath.FloorInt(p.X), mmath.FloorInt(p.Y), color)
			screen.SetZPixel(mmath.FloorInt(p.X), mmath.FloorInt(p.Y), depth)
		}
	}
}

func drawTriangle(mesh *mcgi_mesh.Mesh, screen *screen, triangleId int) {
	triangle := mesh.TriangleList[triangleId]
	triangleRender := triangle.ToTriangle4D(1.0).MultiplyMatrix4x4(mesh.RenderMatrix)

	// Discard triangles
	minZ := triangleRender.MaxZ()
	minW := triangleRender.MaxW()
	clipMZ := minZ / minW
	if clipMZ <= -1 || clipMZ >= 1 {
		return
	}

	// To screen space
	triangleScreen := mmath_geom.Triangle4D[float32]{
		A: toScreenSpace2(screen.Width, screen.Height, triangleRender.A),
		B: toScreenSpace2(screen.Width, screen.Height, triangleRender.B),
		C: toScreenSpace2(screen.Width, screen.Height, triangleRender.C),
	}

	// Check orientation
	deltaA := triangleScreen.B.Sub(triangleScreen.A).ToVector2XY()
	deltaB := triangleScreen.C.Sub(triangleScreen.A).ToVector2XY()
	crossT := deltaB.Cross(deltaA)
	if crossT < 0 {
		return
	}

	// Check if triangle inside screen
	if !screen.IsVisible(triangleScreen.ToTriangle3D()) {
		return
	}

	// Sort triangle by vertices
	var top, middle, bottom = triangleRender.TopMiddleBottom()

	// First slope
	r := middle.Y - top.Y
	l := middle.X - top.X
	slope1 := l / r
	if r == 0 {
		slope1 = 0
	}

	// Second slope
	r = bottom.Y - top.Y
	l = bottom.X - top.X
	slope2 := l / r
	if r == 0 {
		slope2 = 0
	}

	// Calculate area
	area := triangleRender.ToTriangle3D().ToTriangle2D().Area()

	// Top middle triangle
	startY := mmath.Clamp(mmath.CeilInt(top.Y-0.5), 0, screen.Height)
	endY := mmath.Clamp(mmath.CeilInt(middle.Y-0.5), 0, screen.Height)
	for y := startY; y < endY; y++ {
		yy := float32(y) + 0.5
		xStart := mmath.CeilInt(((top.X) + (yy-top.Y)*slope1) - 0.5)
		xEnd := mmath.CeilInt(((top.X) + (yy-top.Y)*slope2) - 0.5)
		if xStart > xEnd {
			xStart, xEnd = xEnd, xStart
		}

		drawLine(mesh, triangleRender, area, screen, xStart, y, xEnd, triangleId)
	}

	// Middle bottom triangle
	r = bottom.Y - middle.Y
	l = bottom.X - middle.X
	slope1 = l / r
	if r == 0 {
		slope1 = 0
	}

	startY = mmath.Clamp(mmath.CeilInt(middle.Y-0.5), 0, screen.Height)
	endY = mmath.Clamp(mmath.CeilInt(bottom.Y-0.5), 0, screen.Height)
	for y := startY; y < endY; y++ {
		yy := float32(y) + 0.5
		xStart := mmath.CeilInt(((middle.X) + (yy-middle.Y)*slope1) - 0.5)
		xEnd := mmath.CeilInt(((top.X) + (yy-top.Y)*slope2) - 0.5)
		if xStart > xEnd {
			xStart, xEnd = xEnd, xStart
		}

		drawLine(mesh, triangleRender, area, screen, xStart, y, xEnd, triangleId)
	}
}

func toScreenSpace2(width int, height int, v mmath_la.Vector4[float32]) mmath_la.Vector4[float32] {
	// Проверка деления на ноль
	if v.W == 0 {
		return mmath_la.Vector4[float32]{}
	}

	// Преобразование X и Y координат
	x := (v.X/v.W + 1) * 0.5 * float32(width)
	y := (1 - v.Y/v.W) * 0.5 * float32(height)

	return mmath_la.Vector4[float32]{x, y, v.Z, v.W}
}
