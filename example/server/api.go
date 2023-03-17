package main

import (
	"github.com/maldan/go-ml/db/mdb"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_response "github.com/maldan/go-ml/server/response"
	ml_crypto "github.com/maldan/go-ml/util/crypto"
	ml_console "github.com/maldan/go-ml/util/io/console"
	ml_time "github.com/maldan/go-ml/util/time"
)

type User struct{}
type Template struct{}

type ArgsX struct {
	A int    `json:"a"`
	B int    `json:"b"`
	X string `json:"x"`
}

type Args2 struct {
	A       string           `json:"a"`
	B       int              `json:"b"`
	Created ml_time.DateTime `json:"created"`
}

type Gasofeal struct {
	Name    string           `json:"name"`
	Lox     string           `json:"lox"`
	Gas     string           `json:"gas"`
	A       string           `json:"a"`
	B       int              `json:"b"`
	Created ml_time.DateTime `json:"created"`
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
	}).Result
}

func (u User) PostIndex3() {
	DataBase["x"].Insert(Gasofeal{
		Name:    ml_crypto.UID(12),
		Gas:     ml_crypto.UID(12),
		Created: ml_time.Now(),
	})
}

func (u User) PostIndex(x Args2) any {
	ml_console.PrettyPrint(x)
	return x
	// fmt.Printf("%+v\n", x.Created)
}
