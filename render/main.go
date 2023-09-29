package mrender

import (
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mr_layer "github.com/maldan/go-ml/render/layer"
	mrender_mesh "github.com/maldan/go-ml/render/mesh"
	ml_mouse "github.com/maldan/go-ml/util/io/mouse"
	ml_color "github.com/maldan/go-ml/util/media/color"
)

type RenderState struct {
	IsRender bool
	Delta    float32
}

type RenderLight struct {
	Direction mmath_la.Vector3[float32]
	Ambient   mmath_la.Vector3[float32]
	Color     ml_color.ColorRGB[float32]
}

type RenderEngine struct {
	Main       mr_layer.MainLayer
	StaticMesh mr_layer.StaticMeshLayer
	Point      mr_layer.PointLayer
	Line       mr_layer.LineLayer
	Text       mr_layer.TextLayer
	UI         mr_layer.UILayer

	ScreenSize mmath_la.Vector2[float32]

	Camera      mrender_camera.PerspectiveCamera
	UICamera    mrender_camera.OrthographicCamera
	Light       RenderLight
	Scissors    mmath_geom.Rectangle[float32]
	UseScissors bool
	State       RenderState
}

var State RenderEngine = RenderEngine{}

//func AddStaticMesh(mesh mrender_mesh.Mesh, position mmath_la.Vector3[float32], rotation mmath_la.Vector3[float32]) *mrender_mesh.Mesh {
func AddStaticMesh(mesh mrender_mesh.Mesh, position mmath_la.Vector3[float32], rotation mmath_la.Vector3[float32]) *mrender_mesh.Mesh {
	if mesh.Normal == nil {
		for i := 0; i < len(mesh.Vertices); i++ {
			mesh.Normal = append(mesh.Normal, mmath_la.Vector3[float32]{0, 0, 0})
		}
	}
	if mesh.UV == nil {
		for i := 0; i < len(mesh.Vertices); i++ {
			mesh.UV = append(mesh.UV, mmath_la.Vector2[float32]{0, 0})
		}
	}
	if mesh.Color == nil {
		for i := 0; i < len(mesh.Vertices); i++ {
			mesh.Color = append(mesh.Color, ml_color.ColorRGBA[float32]{1, 1, 1, 1})
		}
	}

	// Create copy of mesh
	copyMesh := mrender_mesh.Mesh{}
	copyMesh.Vertices = make([]mmath_la.Vector3[float32], len(mesh.Vertices))
	copyMesh.UV = make([]mmath_la.Vector2[float32], len(mesh.UV))
	copyMesh.Normal = make([]mmath_la.Vector3[float32], len(mesh.Normal))
	copyMesh.Color = make([]ml_color.ColorRGBA[float32], len(mesh.Color))
	copyMesh.Indices = make([]uint16, len(mesh.Indices))

	// Copy vertices
	copy(copyMesh.Vertices, mesh.Vertices)
	copy(copyMesh.Normal, mesh.Normal)
	copy(copyMesh.UV, mesh.UV)
	copy(copyMesh.Color, mesh.Color)
	copy(copyMesh.Indices, mesh.Indices)

	// Apply matrix
	copyMesh.RotateY(rotation.Y)
	copyMesh.SetPosition(position)

	State.StaticMesh.MeshList = append(State.StaticMesh.MeshList, copyMesh)
	State.StaticMesh.IsChanged = true
	return &State.StaticMesh.MeshList[len(State.StaticMesh.MeshList)-1]
}

func ClearStaticMesh() {
	State.StaticMesh.MeshList = State.StaticMesh.MeshList[:0]
	State.StaticMesh.IsChanged = true
}

func ClearDynamicMesh() {
	// State.Main.AllocatedMesh = State.Main.AllocatedMesh[:0]
	for i := 0; i < len(State.Main.MeshInstanceList); i++ {
		State.Main.MeshInstanceList[i].IsActive = false
	}
}

func AllocateMesh(mesh mrender_mesh.Mesh) *mrender_mesh.Mesh {
	mesh.Id = len(State.Main.AllocatedMesh)

	if mesh.Normal == nil {
		for i := 0; i < len(mesh.Vertices); i++ {
			mesh.Normal = append(mesh.Normal, mmath_la.Vector3[float32]{0, 0, 0})
		}
	}
	if mesh.UV == nil {
		for i := 0; i < len(mesh.Vertices); i++ {
			mesh.UV = append(mesh.UV, mmath_la.Vector2[float32]{0, 0})
		}
	}

	// mesh.Color = ml_color.ColorRGBA[float32]{1, 1, 1, 1}
	State.Main.AllocatedMesh = append(State.Main.AllocatedMesh, mesh)
	return &State.Main.AllocatedMesh[len(State.Main.AllocatedMesh)-1]
}

func InstanceMesh(mesh *mrender_mesh.Mesh) *mrender_mesh.MeshInstance {
	instance := mrender_mesh.MeshInstance{
		Id:        mesh.Id,
		Scale:     mmath_la.Vector3[float32]{1, 1, 1},
		Color:     ml_color.ColorRGBA[float32]{1, 1, 1, 1},
		IsVisible: true,
		IsActive:  true,
	}

	// Find free instance cell
	for i := 0; i < len(State.Main.MeshInstanceList); i++ {
		if State.Main.MeshInstanceList[i].IsActive {
			continue
		}

		State.Main.MeshInstanceList[i] = instance
		return &State.Main.MeshInstanceList[i]
	}

	panic("mesh instance overflow")
}

func Init() {
	State.Main.Init()
	State.StaticMesh.Init()
	State.Point.Init()
	State.Line.Init()
	State.Text.Init()
	State.UI.Init()

	// Init global camera
	State.Camera = mrender_camera.PerspectiveCamera{Fov: 45, AspectRatio: 1}
	State.Camera.Scale = mmath_la.Vector3[float32]{1, 1, 1}
	State.Camera.Position.Z = 15.5
	State.Camera.Near = 0.01
	State.Camera.Far = 1000

	// Default light
	State.Light.Direction = mmath_la.Vector3[float32]{0.3, 0.4, 0.8}
	State.Light.Ambient = mmath_la.Vector3[float32]{0.3, 0.3, 0.3}
	State.Light.Color = ml_color.ColorRGB[float32]{1, 1, 1}

	State.UICamera.Scale = mmath_la.Vector3[float32]{1, 1, 1}
	State.UICamera.Area.MinX = 0
	State.UICamera.Area.MaxX = 320
	State.UICamera.Area.MinY = 0
	State.UICamera.Area.MaxY = 240
	State.UICamera.ApplyMatrix()
}

func DrawLine(from mmath_la.Vector3[float32], to mmath_la.Vector3[float32], color ml_color.ColorRGBA[float32], isUi bool) {
	if isUi && State.UseScissors {
		isOutFX := false
		isOutFY := false
		isOutTX := false
		isOutTY := false

		if from.X < State.Scissors.MinX {
			from.X = State.Scissors.MinX
			isOutFX = true
		}
		if from.X > State.Scissors.MaxX {
			from.X = State.Scissors.MaxX
			isOutFX = true
		}
		if from.Y < State.Scissors.MinY {
			from.Y = State.Scissors.MinY
			isOutFY = true
		}
		if from.Y > State.Scissors.MaxY {
			from.Y = State.Scissors.MaxY
			isOutFY = true
		}

		if to.X < State.Scissors.MinX {
			to.X = State.Scissors.MinX
			isOutTX = true
		}
		if to.X > State.Scissors.MaxX {
			to.X = State.Scissors.MaxX
			isOutTX = true
		}
		if to.Y < State.Scissors.MinY {
			to.Y = State.Scissors.MinY
			isOutTY = true
		}
		if to.Y > State.Scissors.MaxY {
			to.Y = State.Scissors.MaxY
			isOutTY = true
		}

		if isOutFX && isOutTX && isOutFY && isOutTY {
			return
		}
	}

	State.Line.LineList = append(State.Line.LineList, mrender_mesh.Line{
		From:  from,
		To:    to,
		Color: color,
		IsUi:  isUi,
	})
}

func EnableScissors(r mmath_geom.Rectangle[float32]) {
	State.Scissors = r
	State.UseScissors = true
}

func DisableScissors() {
	State.UseScissors = false
}

func DrawRectangle2(r mmath_geom.Rectangle[float32], color ml_color.ColorRGBA[float32], isUi bool) {
	DrawRectangle(
		mmath_la.Vector3[float32]{r.MinX, r.MinY, 0},
		mmath_la.Vector3[float32]{r.MaxX, r.MaxY, 0},
		color,
		isUi,
	)
}

func DrawRectangle(from mmath_la.Vector3[float32], to mmath_la.Vector3[float32], color ml_color.ColorRGBA[float32], isUi bool) {
	tFrom := from
	tTo := to
	tFrom.Z = to.Z

	// Top line
	tFrom.Y = from.Y
	tTo.Y = from.Y
	DrawLine(tFrom, tTo, color, isUi)

	// Bottom line
	tFrom.Y = to.Y
	tTo.Y = to.Y
	DrawLine(tFrom, tTo, color, isUi)

	tFrom = from
	tTo = to

	// Left line
	tFrom.X = from.X
	tTo.X = from.X
	DrawLine(tFrom, tTo, color, isUi)

	// To line
	tFrom.X = to.X
	tTo.X = to.X
	DrawLine(tFrom, tTo, color, isUi)
}

func DrawCuboid(cuboid mmath_geom.Cuboid[float32], color ml_color.ColorRGBA[float32], isUi bool) {
	from := mmath_la.Vector3[float32]{
		cuboid.MinX, cuboid.MinY, cuboid.MinZ,
	}
	to := mmath_la.Vector3[float32]{
		cuboid.MaxX, cuboid.MaxY, cuboid.MaxZ,
	}

	f1 := from
	t1 := to
	f1.Z = t1.Z
	DrawRectangle(f1, t1, color, isUi)

	f2 := from
	t2 := to
	t2.Z = f2.Z
	DrawRectangle(f2, t2, color, isUi)

	f := from
	t := to
	t.X = f.X
	t.Y = f.Y
	DrawLine(f, t, color, isUi)

	f = from
	t = to
	f.X = to.X
	t.Y = from.Y
	DrawLine(f, t, color, isUi)

	f = from
	t = to
	t.X = f.X
	f.Y = t.Y
	DrawLine(f, t, color, isUi)

	f = from
	t = to
	f.X = to.X
	f.Y = t.Y
	DrawLine(f, t, color, isUi)
}

func DrawPoint(to mmath_la.Vector3[float32], size float32, color ml_color.ColorRGBA[float32]) {
	State.Point.PointList = append(State.Point.PointList, mr_layer.Point{
		Position: to,
		Size:     size,
		Color:    color,
	})
}

func LoadFont(name string, font mr_layer.TextFont) {
	/*State.Text.FontMap[name] = mr_layer.TextFont{
		Symbol: map[uint8]mmath_geom.Rectangle[float32]{},
	}

	for c, r := range charMap {
		State.Text.FontMap[name].Symbol[c] = mrender_uv.GetArea(r.MinX, r.MinY, r.MaxX, r.MaxY, 1024, 1024)
	}*/
	State.Text.FontMap[name] = font
	State.UI.FontMap[name] = font
}

/*func LoadUIFont(name string, font mr_layer.UITextFont) {
	State.UI.FontMap[name] = font
}
*/

/*func DrawUIText(fontName string, text string, size float32, pos mmath_la.Vector2[float32]) {
	font, ok := State.UI.FontMap[fontName]
	if !ok {
		fmt.Printf("Font %v not found\n", fontName)
	}

	offsetX := float32(0)
	for _, r := range text {
		c := fmt.Sprintf("%c", r)

		c2 := font.Symbol[c]

		rect := mrender_uv.GetArea(c2.X, c2.Y, font.Size.X, font.Size.Y, 1024, 1024)

		DrawUI(
			rect,
			pos.AddXY(offsetX, 0),
			font.Size.Scale(size),
			0,
			mmath_la.Vector2[float32]{},
		)

		if font.SymbolSize[c].X == 0 {
			offsetX += font.Size.X * size
		} else {
			offsetX += font.SymbolSize[c].X * size
		}
	}
}*/

/*func DrawText(font string, text string, size float32, pos mmath_la.Vector3[float32]) {
	State.Text.TextList = append(State.Text.TextList, mr_layer.Text{
		Font:     font,
		Content:  text,
		Size:     size,
		Position: pos,
	})
}*/

func DrawUI(
	uv mmath_geom.Rectangle[float32],
	pos mmath_la.Vector2[float32],
	size mmath_la.Vector2[float32],
	rotation float32,
	pivot mmath_la.Vector2[float32],
) {
	State.UI.ElementList = append(State.UI.ElementList, mr_layer.UIElement{
		UvArea:   uv,
		Position: pos,
		Rotation: rotation,
		Scale:    size,
		Color:    ml_color.ColorRGBA[float32]{1, 1, 1, 1},
		Pivot:    pivot,
	})
}

/*func DrawUI2(x struct{ Position mmath_la.Vector2[float32] }) {

}*/

func DrawButton(
	uv mmath_geom.Rectangle[float32],
	pos mmath_la.Vector2[float32],
	size mmath_la.Vector2[float32],
	onClick func()) {

	mp := ml_mouse.GetPosition()
	mp.X = ((mp.X + 1) / 2) * State.ScreenSize.X
	mp.Y = State.ScreenSize.Y - (((mp.Y + 1) / 2) * State.ScreenSize.Y)

	halfX := size.X / 2
	halfY := size.Y / 2

	rect := mmath_geom.Rectangle[float32]{MinX: pos.X, MaxX: pos.X + size.X, MinY: pos.Y, MaxY: pos.Y + size.Y}
	rect = rect.AddXY(-halfX, -halfY)

	if rect.IntersectPoint(mp) {
		if ml_mouse.IsMouseDown(ml_mouse.LeftButton) {
			pos.Y += 2
		}
		if ml_mouse.IsMouseClick(ml_mouse.LeftButton) {
			onClick()
		}

		State.UI.ElementList = append(State.UI.ElementList, mr_layer.UIElement{
			UvArea:   uv,
			Position: pos,
			Scale:    size,
			Color:    ml_color.ColorRGBA[float32]{1, 1, 1, 0.8},
		})
	} else {
		State.UI.ElementList = append(State.UI.ElementList, mr_layer.UIElement{
			UvArea:   uv,
			Position: pos,
			Scale:    size,
			Color:    ml_color.ColorRGBA[float32]{1, 1, 1, 1},
		})
	}
}

func DrawMeshNormals(instance *mrender_mesh.MeshInstance) {
	mesh := State.Main.AllocatedMesh[instance.Id]

	for i := 0; i < len(mesh.Indices); i += 3 {
		a := mesh.Vertices[mesh.Indices[i]]
		an := mesh.Normal[mesh.Indices[i]]
		DrawLine(
			instance.Position.Add(a),
			instance.Position.Add(a).Add(an),
			ml_color.ColorRGBA[float32]{1, 0, 0, 1},
			false,
		)

		a = mesh.Vertices[mesh.Indices[i+1]]
		an = mesh.Normal[mesh.Indices[i+1]]
		DrawLine(
			instance.Position.Add(a),
			instance.Position.Add(a).Add(an),
			ml_color.ColorRGBA[float32]{0, 1, 0, 1},
			false,
		)

		a = mesh.Vertices[mesh.Indices[i+2]]
		an = mesh.Normal[mesh.Indices[i+2]]
		DrawLine(
			instance.Position.Add(a),
			instance.Position.Add(a).Add(an),
			ml_color.ColorRGBA[float32]{0, 0, 1, 1},
			false,
		)
	}
}

func (r *RenderEngine) Render() {
	r.Camera.ApplyMatrix()
	r.UICamera.ApplyMatrix()

	r.Main.Render()
	r.StaticMesh.Render()
	r.Point.Render()
	r.Line.Render()
	r.Text.Render()
	r.UI.Render()

	/*

		r.Text.Render()
	*/
}
