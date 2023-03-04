package ml_net_test

import (
	"fmt"
	ml_net "github.com/maldan/go-ml/util/io/net"
	"testing"
)

type Todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestGet(t *testing.T) {
	response := ml_net.Get[[]Todo]("https://jsonplaceholder.typicode.com/todos", ml_net.RequestOptions{})
	todoList, _ := response.Unpack()

	fmt.Printf("%v\n", todoList)

	if response.StatusCode != 200 {
		t.Errorf("Can't request")
	}
	if len(todoList) == 0 {
		t.Errorf("No data")
	}
	if todoList[0].Id != 1 {
		t.Errorf("Can't parse json")
	}
}

func TestPost(t *testing.T) {
	response := ml_net.Request[any]("https://jsonplaceholder.typicode.com/posts", "POST", ml_net.RequestOptions{
		Data: &Post{
			UserId: 1,
			Id:     1,
			Title:  "A",
			Body:   "B",
		},
	})
	response.Close()

	if response.StatusCode != 201 {
		t.Errorf("Can't create")
	}
}
