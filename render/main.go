package mrender

import (
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mr_camera "github.com/maldan/go-ml/render/camera"
	mr_layer "github.com/maldan/go-ml/render/layer"
	mrender_mesh "github.com/maldan/go-ml/render/mesh"
	mrender_uv "github.com/maldan/go-ml/render/uv"
	ml_mouse "github.com/maldan/go-ml/util/io/mouse"
	ml_color "github.com/maldan/go-ml/util/media/color"
)

type RenderEngine struct {
	Main       mr_layer.MainLayer
	StaticMesh mr_layer.StaticMeshLayer
	Point      mr_layer.PointLayer
	Line       mr_layer.LineLayer
	Text       mr_layer.TextLayer
	UI         mr_layer.UILayer

	ScreenSize mmath_la.Vector2[float32]

	GlobalCamera mr_camera.PerspectiveCamera
}

var State RenderEngine = RenderEngine{}

func AddStaticMesh(mesh mrender_mesh.Mesh, position mmath_la.Vector3[float32]) *mrender_mesh.Mesh {
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
	copyMesh.SetPosition(position)

	State.StaticMesh.MeshList = append(State.StaticMesh.MeshList, copyMesh)
	State.StaticMesh.IsChanged = true
	return &State.StaticMesh.MeshList[len(State.StaticMesh.MeshList)-1]
}

func DeleteStaticMesh(mesh *mrender_mesh.Mesh) {

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
	State.GlobalCamera = mr_camera.PerspectiveCamera{Fov: 45, AspectRatio: 1}
	State.GlobalCamera.Scale = mmath_la.Vector3[float32]{1, 1, 1}
	State.GlobalCamera.Position.Z = 15.5

	State.UI.Camera.Scale = mmath_la.Vector3[float32]{1, 1, 1}
	State.UI.Camera.Area.Left = 0
	State.UI.Camera.Area.Right = 320
	State.UI.Camera.Area.Top = 0
	State.UI.Camera.Area.Bottom = 240

	State.UI.Camera.ApplyMatrix()
}

func DrawLine(from mmath_la.Vector3[float32], to mmath_la.Vector3[float32], color ml_color.ColorRGB[float32]) {
	State.Line.LineList = append(State.Line.LineList, mrender_mesh.Line{
		From:  from,
		To:    to,
		Color: color,
	})
}

func DrawRectangle(from mmath_la.Vector3[float32], to mmath_la.Vector3[float32], color ml_color.ColorRGB[float32]) {
	tFrom := from
	tTo := to
	tFrom.Z = to.Z

	// Top line
	tFrom.Y = from.Y
	tTo.Y = from.Y
	DrawLine(tFrom, tTo, color)

	// Bottom line
	tFrom.Y = to.Y
	tTo.Y = to.Y
	DrawLine(tFrom, tTo, color)

	tFrom = from
	tTo = to

	// Left line
	tFrom.X = from.X
	tTo.X = from.X
	DrawLine(tFrom, tTo, color)

	// To line
	tFrom.X = to.X
	tTo.X = to.X
	DrawLine(tFrom, tTo, color)
}

func DrawPoint(to mmath_la.Vector3[float32], size float32, color ml_color.ColorRGBA[float32]) {
	State.Point.PointList = append(State.Point.PointList, mr_layer.Point{
		Position: to,
		Size:     size,
		Color:    color,
	})
}

func LoadFont(name string, charMap map[uint8]mmath_geom.Rectangle[float32]) {
	State.Text.FontMap[name] = mr_layer.TextFont{
		Symbol: map[uint8]mmath_geom.Rectangle[float32]{},
	}

	for c, r := range charMap {
		State.Text.FontMap[name].Symbol[c] = mrender_uv.GetArea(r.Left, r.Top, r.Right, r.Bottom, 1024, 1024)
	}
}

func DrawText(font string, text string, size float32, pos mmath_la.Vector3[float32]) {
	State.Text.TextList = append(State.Text.TextList, mr_layer.Text{
		Font:     font,
		Content:  text,
		Size:     size,
		Position: pos,
	})
}

func DrawUI(
	uv mmath_geom.Rectangle[float32],
	pos mmath_la.Vector3[float32],
	size mmath_la.Vector2[float32],
	rotation float32,
	pivot mmath_la.Vector2[float32],
) {
	State.UI.ElementList = append(State.UI.ElementList, mr_layer.UIElement{
		UvArea:    uv,
		Position:  pos,
		Rotation:  mmath_la.Vector3[float32]{0, 0, rotation},
		Scale:     mmath_la.Vector3[float32]{size.X, size.Y, 1},
		Color:     ml_color.ColorRGBA[float32]{1, 1, 1, 1},
		IsVisible: true,
		IsActive:  true,
		Pivot:     pivot,
	})
}

func DrawButton(
	uv mmath_geom.Rectangle[float32],
	pos mmath_la.Vector3[float32],
	size mmath_la.Vector2[float32],
	onClick func()) {

	mp := ml_mouse.GetPosition()
	mp.X = ((mp.X + 1) / 2) * State.ScreenSize.X
	mp.Y = State.ScreenSize.Y - (((mp.Y + 1) / 2) * State.ScreenSize.Y)

	halfX := size.X / 2
	halfY := size.Y / 2

	rect := mmath_geom.Rectangle[float32]{Left: pos.X, Right: pos.X + size.X, Top: pos.Y, Bottom: pos.Y + size.Y}
	rect = rect.Add(-halfX, -halfY)

	if rect.IntersectPoint(mp) {
		if ml_mouse.IsMouseDown(ml_mouse.LeftButton) {
			pos.Y += 2
		}
		if ml_mouse.IsMouseClick(ml_mouse.LeftButton) {
			onClick()
		}

		State.UI.ElementList = append(State.UI.ElementList, mr_layer.UIElement{
			UvArea:    uv,
			Position:  pos,
			Scale:     mmath_la.Vector3[float32]{size.X, size.Y, 1},
			Color:     ml_color.ColorRGBA[float32]{1, 1, 1, 0.8},
			IsVisible: true,
			IsActive:  true,
		})
	} else {
		State.UI.ElementList = append(State.UI.ElementList, mr_layer.UIElement{
			UvArea:    uv,
			Position:  pos,
			Scale:     mmath_la.Vector3[float32]{size.X, size.Y, 1},
			Color:     ml_color.ColorRGBA[float32]{1, 1, 1, 1},
			IsVisible: true,
			IsActive:  true,
		})
	}
}

func (r *RenderEngine) Render() {
	r.GlobalCamera.ApplyMatrix()
	r.Main.Camera = r.GlobalCamera
	r.StaticMesh.Camera = r.GlobalCamera
	r.Point.Camera = r.GlobalCamera
	r.Line.Camera = r.GlobalCamera
	r.Text.Camera = r.GlobalCamera

	r.Main.Render()
	r.StaticMesh.Render()
	r.Point.Render()
	r.Line.Render()
	r.Text.Render()
	r.UI.Render()
}
