package mrender

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mr_camera "github.com/maldan/go-ml/render/camera"
	mr_layer "github.com/maldan/go-ml/render/layer"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	ml_color "github.com/maldan/go-ml/util/media/color"
)

type RenderEngine struct {
	Main  mr_layer.MainLayer
	Point mr_layer.PointLayer
	Line  mr_layer.LineLayer

	GlobalCamera mr_camera.PerspectiveCamera
}

var State RenderEngine = RenderEngine{}

func AllocateMesh(mesh *mr_mesh.Mesh) *mr_mesh.Mesh {
	mesh.Id = len(State.Main.AllocatedMesh)
	State.Main.AllocatedMesh = append(State.Main.AllocatedMesh, mesh)
	return mesh
}

func InstanceMesh(mesh *mr_mesh.Mesh) *mr_mesh.MeshInstance {
	instance := mr_mesh.MeshInstance{
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
	State.Point.Init()
	State.Line.Init()
}

func DrawLine(from mmath_la.Vector3[float32], to mmath_la.Vector3[float32], color ml_color.ColorRGB[float32]) {
	State.Line.LineList = append(State.Line.LineList, mr_mesh.Line{
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

func DrawPoint(to mmath_la.Vector3[float32]) {
	State.Point.PointList = append(State.Point.PointList, to)
}

func (r *RenderEngine) Render() {
	r.GlobalCamera.ApplyMatrix()
	r.Main.Camera = r.GlobalCamera
	r.Point.Camera = r.GlobalCamera
	r.Line.Camera = r.GlobalCamera

	r.Main.Render()
	r.Point.Render()
	r.Line.Render()
}
