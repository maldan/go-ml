package mr

import (
	mr_layer "github.com/maldan/go-ml/render/layer"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
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

type Color32 struct {
	R float32
	G float32
	B float32
	A float32
}

type Line struct {
	From ml_geom.Vector3[float32]
	To   ml_geom.Vector3[float32]
}

func DebugLine(from ml_geom.Vector3[float32], to ml_geom.Vector3[float32]) {
	State.Line.LineList = append(State.Line.LineList, ml_geom.Line[float32, ml_geom.Vector3[float32]]{
		from, to,
	})
}

func DebugPoint(to ml_geom.Vector3[float32]) {
	State.Point.PointList = append(State.Point.PointList, to)
}

func (r *RenderEngine) Render() {
	r.Main.Render()
	r.Point.Render()
	r.Line.Render()
}
