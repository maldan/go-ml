package gosn_driver

import (
	"github.com/maldan/go-ml/util/encode/gosn"
)

type Container struct {
	// NameToId ml_gosn.NameToId
}

/*func (g *Container) GetHeader() []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, uint8(len(g.NameToId)))

	// Write name to id
	for name, id := range g.NameToId {
		// Name length
		bytes = append(bytes, uint8(len(name)))

		// Name
		bytes = append(bytes, name...)

		// Id
		bytes = append(bytes, id)
	}

	return bytes
}

func (g *Container) SetHeader(bytes []byte) {
	offset := 0

	// Amount
	amount := int(bytes[offset])
	offset += 1

	// Fill maps
	for i := 0; i < amount; i++ {
		// Name length
		nameLen := int(bytes[offset])
		offset += 1

		// Name
		name := string(bytes[offset : offset+nameLen])
		offset += nameLen

		// Fill map
		g.NameToId[name] = bytes[offset]
		offset += 1
	}
}*/

/*func (g *Container) Prepare(v any) {
	nid := ml_gosn.NameToId{}
	nid.FromStruct(v)
	_ = ml_gosn.MarshalExt(v, nid)
	g.NameToId = nid
}*/

func (g *Container) Marshal(v any) []byte {
	return ml_gosn.Marshal(v)
}

func (g *Container) Unmarshall(b []byte, out any) {
	ml_gosn.Unmarshall(b, out)
}

/*func (g *Container) GetMapper(fieldList string, out any) any {
	return NewMapper(g.NameToId, fieldList, out)
}*/

/*func (g *Container) GetStruct() map[string]string {
	out := map[string]string{}
	for k, _ := range g.NameToId {
		out[k] = ""
	}
	return out
}
*/
