package ml_template

import (
	"bytes"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	"path"
	"text/template"
)

func ExecuteFromString(inputData string, data any) (string, error) {
	// Prepare
	tmpl, err := template.New("template").Parse(inputData)
	if err != nil {
		return "", err
	}

	// Execute
	buf := new(bytes.Buffer)
	err2 := tmpl.Execute(buf, data)
	if err2 != nil {
		return "", err2
	}

	// Output
	return string(buf.Bytes()), nil
}

func Execute(inputPath string, output string, data any) error {
	// Prepare
	tmpl, err := template.New(path.Base(inputPath)).ParseFiles(inputPath)
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
