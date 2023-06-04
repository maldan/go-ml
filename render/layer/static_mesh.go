package mrender_layer

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mrender_mesh "github.com/maldan/go-ml/render/mesh"
	ml_color "github.com/maldan/go-ml/util/media/color"
	"reflect"
	"unsafe"
)

type StaticMeshLayer struct {
	VertexList []float32
	UvList     []float32
	NormalList []float32
	ColorList  []float32
	IndexList  []uint16

	/*VertexAmount int
	IndexAmount  int
	UvAmount     int
	ColorAmount  int*/

	IsChanged bool

	MeshList []mrender_mesh.Mesh
	//Camera mr_camera.PerspectiveCamera

	state map[string]any
}

type StatisMeshArgs struct {
	Mesh     mrender_mesh.Mesh
	Position mmath_la.Vector3[float32]
	Rotation mmath_la.Vector3[float32]
	UvOffset mmath_la.Vector2[float32]
}

func (l *StaticMeshLayer) Add(args StatisMeshArgs) *mrender_mesh.Mesh {
	if args.Mesh.Normal == nil {
		for i := 0; i < len(args.Mesh.Vertices); i++ {
			args.Mesh.Normal = append(args.Mesh.Normal, mmath_la.Vector3[float32]{0, 0, 0})
		}
	}
	if args.Mesh.UV == nil {
		for i := 0; i < len(args.Mesh.Vertices); i++ {
			args.Mesh.UV = append(args.Mesh.UV, mmath_la.Vector2[float32]{0, 0})
		}
	}
	if args.Mesh.Color == nil {
		for i := 0; i < len(args.Mesh.Vertices); i++ {
			args.Mesh.Color = append(args.Mesh.Color, ml_color.ColorRGBA[float32]{1, 1, 1, 1})
		}
	}

	// Create copy of mesh
	copyMesh := mrender_mesh.Mesh{}
	copyMesh.Vertices = make([]mmath_la.Vector3[float32], len(args.Mesh.Vertices))
	copyMesh.UV = make([]mmath_la.Vector2[float32], len(args.Mesh.UV))
	copyMesh.Normal = make([]mmath_la.Vector3[float32], len(args.Mesh.Normal))
	copyMesh.Color = make([]ml_color.ColorRGBA[float32], len(args.Mesh.Color))
	copyMesh.Indices = make([]uint16, len(args.Mesh.Indices))

	// Copy vertices
	copy(copyMesh.Vertices, args.Mesh.Vertices)
	copy(copyMesh.Normal, args.Mesh.Normal)
	copy(copyMesh.UV, args.Mesh.UV)
	copy(copyMesh.Color, args.Mesh.Color)
	copy(copyMesh.Indices, args.Mesh.Indices)

	// Apply matrix
	copyMesh.RotateY(args.Rotation.Y)
	copyMesh.SetPosition(args.Position)
	copyMesh.OffsetUv(args.UvOffset)

	l.MeshList = append(l.MeshList, copyMesh)
	l.IsChanged = true
	return &l.MeshList[len(l.MeshList)-1]
}

func (l *StaticMeshLayer) Init() {
	l.VertexList = make([]float32, 0, 1024)
	l.NormalList = make([]float32, 0, 1024)
	l.UvList = make([]float32, 0, 1024)
	l.ColorList = make([]float32, 0, 1024)

	l.MeshList = make([]mrender_mesh.Mesh, 0, 128)
	l.IndexList = make([]uint16, 1024)
}

func (l *StaticMeshLayer) Build() {
	// Clear before start
	l.VertexList = l.VertexList[:0]
	l.UvList = l.UvList[:0]
	l.NormalList = l.NormalList[:0]
	l.ColorList = l.ColorList[:0]
	l.IndexList = l.IndexList[:0]

	lastMaxIndex := uint16(0)
	for i := 0; i < len(l.MeshList); i++ {
		mesh := l.MeshList[i]

		// Copy vertex
		for j := 0; j < len(mesh.Vertices); j++ {
			v := mesh.Vertices[j]
			l.VertexList = append(l.VertexList, v.X, v.Y, v.Z)

			// Copy normal
			n := mesh.Normal[j]
			l.NormalList = append(l.NormalList, n.X, n.Y, n.Z)
		}

		// Copy index
		maxIndex := lastMaxIndex
		for j := 0; j < len(mesh.Indices); j++ {
			iv := mesh.Indices[j] + maxIndex
			if iv > lastMaxIndex {
				lastMaxIndex = iv
			}
			l.IndexList = append(l.IndexList, iv)
		}
		lastMaxIndex += 1

		// Copy uv
		for j := 0; j < len(mesh.UV); j++ {
			v := mesh.UV[j]
			l.UvList = append(l.UvList, v.X, v.Y)
		}

		// Copy color
		for j := 0; j < len(mesh.Vertices); j++ {
			c := mesh.Color[j]
			l.ColorList = append(l.ColorList, c.R, c.G, c.B, c.A)
		}
	}
}

func (l *StaticMeshLayer) Render() {
	//l.Camera.ApplyMatrix()
	if l.IsChanged {
		l.Build()
		l.IsChanged = false
	}
}

func (l *StaticMeshLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	normalHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.NormalList))
	uvHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.UvList))
	indexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.IndexList))
	colorHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ColorList))

	if l.state == nil {
		l.state = map[string]any{
			"vertexPointer": vertexHeader.Data,
			"normalPointer": normalHeader.Data,
			"uvPointer":     uvHeader.Data,
			"indexPointer":  indexHeader.Data,
			"colorPointer":  colorHeader.Data,

			/*"vertexAmount":   l.VertexAmount,
			"normalAmount":   l.VertexAmount,
			"uvAmount":       l.UvAmount,
			"indexAmount":    l.IndexAmount,
			"positionAmount": l.VertexAmount,
			"rotationAmount": l.VertexAmount,
			"scaleAmount":    l.VertexAmount,
			"colorAmount":    l.ColorAmount,*/

			//"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		/*l.state["vertexAmount"] = l.VertexAmount
		l.state["normalAmount"] = l.VertexAmount
		l.state["uvAmount"] = l.UvAmount
		l.state["indexAmount"] = l.IndexAmount
		l.state["positionAmount"] = l.VertexAmount
		l.state["rotationAmount"] = l.VertexAmount
		l.state["scaleAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.ColorAmount*/
	}

	return l.state
}
