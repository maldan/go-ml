package ml_template

import (
	"bytes"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	"path"
	"text/template"
)

func Execute(input string, output string, data any) error {
	// Prepare
	tmpl, err := template.New(path.Base(input)).ParseFiles(input)
	if err != nil {
		return err
	}

	// Execute
	buf := new(bytes.Buffer)
	err2 := tmpl.Execute(buf, data)
	if err2 != nil {
		return err2
	}

	// Output
	err3 := ml_file.New(output).Write(buf.Bytes())
	return err3
}
