package ml_gosn

/*
type NameToId map[string]uint8
type IdToName map[uint8]string

func (n *NameToId) FromStruct(v any) {
	n.Add(GetNameList(v)...)
}

func (n *NameToId) Add(names ...string) {
	for _, nm := range names {
		(*n)[nm] = uint8(len(*n)) + 1
	}
}

func (n *NameToId) Invert() IdToName {
	out := IdToName{}

	for k, v := range *n {
		out[v] = k
	}

	return out
}

// GetNameList return unique list of names from given struct
func GetNameList(v any) []string {
	m := map[string]any{}
	l := make([]string, 0)

	// Get all fields names
	__getFieldList(v, &l)
	for _, name := range l {
		m[name] = true
	}

	// Get unique list
	out := make([]string, 0)
	for name, _ := range m {
		out = append(out, name)
	}
	return out
}

func __getFieldList(v any, out *[]string) {
	typeOf := reflect.TypeOf(v)

	// Skip non struct or non slice types
	if !(typeOf.Kind() == reflect.Struct || typeOf.Kind() == reflect.Slice) {
		return
	}

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		// Skip private
		if string(field.Name[0]) == strings.ToLower(string(field.Name[0])) {
			continue
		}

		*out = append(*out, field.Name)

		if field.Type.Kind() == reflect.Struct {
			__getFieldList(reflect.New(field.Type).Elem().Interface(), out)
		}
		if field.Type.Kind() == reflect.Slice {
			__getFieldList(reflect.New(field.Type.Elem()).Elem().Interface(), out)
		}
	}
}

func NameListToIdList(fieldList string, nameToId NameToId) []uint8 {
	// Field list
	fieldList2 := ml_slice.Map(strings.Split(fieldList, ","), func(t string) string {
		return strings.Trim(t, " ")
	})
	if len(fieldList2) == 0 {
		return []uint8{}
	}

	fieldIdList := make([]uint8, 0)
	for _, v := range fieldList2 {
		id, ok := nameToId[v]
		if !ok {
			panic(fmt.Sprintf("field %v not found", id))
		}
		fieldIdList = append(fieldIdList, id)
	}
	return fieldIdList
}
*/
