package main

import (
	"fmt"
	"github.com/maldan/go-ml/db/mdb"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_response "github.com/maldan/go-ml/server/response"
	ml_crypto "github.com/maldan/go-ml/util/crypto"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_time "github.com/maldan/go-ml/util/time"
	"time"
)

type User struct{}
type Template struct{}

type ArgsX struct {
	A int    `json:"a"`
	B int    `json:"b"`
	X string `json:"x"`
}

type Args2 struct {
	Authorization string           `json:"authorization"`
	A             string           `json:"a"`
	B             int              `json:"b"`
	File          ml_file.File     `json:"file"`
	XFile         ml_file.File     `json:"xfile"`
	Created       ml_time.DateTime `json:"created"`
}

type Gasofeal struct {
	Name    string    `json:"name"`
	Lox     string    `json:"lox"`
	Gas     string    `json:"gas"`
	A       string    `json:"a"`
	B       int       `json:"b"`
	Created time.Time `json:"created"`
}

func (g Gasofeal) ToResponse() map[string]any {
	return map[string]any{
		"a": 1,
	}
}

func (u User) GetGav() ms_response.File {
	return ms_response.File{
		Path: "../test.json",
	}
}

func (u User) GetIndex(ctx *ms_handler.Context, x ArgsX) int {
	return x.A + x.B
}

func (u User) GetIndex2() Gasofeal {
	return Gasofeal{Name: "Gas", Lox: "bas", Gas: "ZZAS"}
}

func (u User) GetIndex3() any {
	return DataBase["tags"].FindBy(mdb.ArgsFind{
		WhereExpression: "Width == 820",
		Where: func(any2 any) bool {
			return any2.(*ImageDescription).Width == 820
		},
	}).Result
}

func (u User) PostIndex3() {
	DataBase["x"].Insert(Gasofeal{
		Name:    ml_crypto.UID(12),
		Gas:     ml_crypto.UID(12),
		Created: time.Now(),
	})
}

func (u User) PostIndex(x Args2) {
	fmt.Printf("%v\n", x.File.Size())
	fmt.Printf("%v\n", x.XFile.Size())
	// ml_console.PrettyPrint(x)
	// fmt.Printf("%+v\n", x.Created)
}
