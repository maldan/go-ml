package mr

import (
	mr_layer "github.com/maldan/go-ml/render/layer"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	ml_image "github.com/maldan/go-ml/util/media/image"
)

type RenderEngine struct {
	Main  mr_layer.MainLayer
	Point mr_layer.PointLayer
	Line  mr_layer.LineLayer
}

var State RenderEngine = RenderEngine{}

func AddMesh(mesh *mr_mesh.Mesh) {
	State.Main.MeshList = append(State.Main.MeshList, mesh)
}

func Init() {
	State.Main.Init()
	State.Point.Init()
	State.Line.Init()
}

func DrawLine(from ml_geom.Vector3[float32], to ml_geom.Vector3[float32], color ml_image.ColorRGB[float32]) {
	State.Line.LineList = append(State.Line.LineList, mr_mesh.Line{
		From:  from,
		To:    to,
		Color: color,
	})
}

func DrawRectangle(from ml_geom.Vector3[float32], to ml_geom.Vector3[float32], color ml_image.ColorRGB[float32]) {
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

func DebugPoint(to ml_geom.Vector3[float32]) {
	State.Point.PointList = append(State.Point.PointList, to)
}

func (r *RenderEngine) Render() {
	r.Main.Render()
	r.Point.Render()
	r.Line.Render()
}
