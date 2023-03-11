package mdb_editor

import (
	"fmt"
	"github.com/Knetic/govaluate"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	"github.com/maldan/go-ml/db/mdb"
	ml_convert "github.com/maldan/go-ml/util/convert"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"log"
	"os"
	"reflect"
)

type UIInfo[T any] struct {
	DB              *mdb.DataTable[T]
	Elements        []ui.Drawable
	TableInfo       *widgets.Paragraph
	PrintInfo       *widgets.Table
	IsShowTableMode bool
	HideColumns     []string
}

func Search[T any](h *UIInfo[T], query string) mdb.SearchResult[T] {
	expression, err := govaluate.NewEvaluableExpression(query)
	if err != nil {
		return mdb.SearchResult[T]{}
	}

	result := h.DB.FindBy(mdb.ArgsFind[T]{
		Where: func(v *T) bool {
			parameters := ml_convert.StructToMap(v)
			r, _ := expression.Evaluate(parameters)
			if r == nil {
				return false
			}
			return r.(bool)
		},
	})
	return result
}

func HandleCommand[T any](h *UIInfo[T], cmd string) {
	// Table mode
	if h.IsShowTableMode {
		ui.Clear()

		tw, _ := ui.TerminalDimensions()
		h.PrintInfo = widgets.NewTable()
		h.PrintInfo.SetRect(0, 0, tw, 20)
		h.PrintInfo.Rows = [][]string{}

		// Select
		result := Search[T](h, cmd)

		// Prepare keys
		ss := h.DB.Container.GetStruct()
		keys := ml_slice.GetKeys(ss)

		vKeys := make([]string, 0)
		vKeys = append(vKeys, "#")
		vKeys = append(vKeys, keys...)

		h.PrintInfo.Rows = append(h.PrintInfo.Rows, vKeys)

		for i := 0; i < result.Count; i++ {
			valueOf := reflect.ValueOf(result.Result[i])
			values := make([]string, 0)
			values = append(values, fmt.Sprintf("%v", i))

			for _, key := range keys {
				field := valueOf.FieldByName(key)
				values = append(values, fmt.Sprintf("%v", field.Interface()))
			}
			h.PrintInfo.Rows = append(h.PrintInfo.Rows, values)
		}

		h.PrintInfo.ColumnWidths = make([]int, len(vKeys))
		for i := 0; i < len(vKeys); i++ {
			h.PrintInfo.ColumnWidths[i] = 20
		}
		h.PrintInfo.ColumnWidths[0] = 5

		return
	}

	if cmd == "show table" {
		h.IsShowTableMode = true

		/*h.PrintInfo = widgets.NewTable()
		h.PrintInfo.SetRect(0, 0, tw, 20)
		h.PrintInfo.Rows = [][]string{}

		result := h.DB.FindBy(mdb.ArgsFind[T]{
			Limit: 10,
			Where: func(*T) bool {
				return true
			},
		})
		ss := h.DB.Container.GetStruct()
		keys := ml_slice.GetKeys(ss)

		vKeys := make([]string, 0)
		vKeys = append(vKeys, "#")
		vKeys = append(vKeys, keys...)

		h.PrintInfo.Rows = append(h.PrintInfo.Rows, vKeys)

		for i := 0; i < result.Count; i++ {
			valueOf := reflect.ValueOf(result.Result[i])
			values := make([]string, 0)
			values = append(values, fmt.Sprintf("%v", i))

			for _, key := range keys {
				field := valueOf.FieldByName(key)
				values = append(values, fmt.Sprintf("%v", field.Interface()))
			}
			h.PrintInfo.Rows = append(h.PrintInfo.Rows, values)
		}

		h.PrintInfo.ColumnWidths = make([]int, len(vKeys))
		for i := 0; i < len(vKeys); i++ {
			h.PrintInfo.ColumnWidths[i] = 20
		}
		h.PrintInfo.ColumnWidths[0] = 5*/
	}
	if cmd == "total" {
		h.TableInfo = widgets.NewParagraph()

		result := h.DB.FindBy(mdb.ArgsFind[T]{
			Where: func(*T) bool {
				return true
			},
		})
		info := fmt.Sprintf("Total: %v\n", result.Count)

		h.TableInfo.Text = info
		h.TableInfo.SetRect(0, 0, 30, 10)
	}
	if cmd == "info" {
		h.TableInfo = widgets.NewParagraph()

		info := fmt.Sprintf("Path: %v\n", h.DB.Path)
		info += fmt.Sprintf("Name: %v\n", h.DB.Name)
		info += fmt.Sprintf("Version: %v\n", h.DB.Header.Version)
		info += fmt.Sprintf("AI: %v\n", h.DB.Header.AutoIncrement)

		h.TableInfo.Text = info
		h.TableInfo.SetRect(0, 0, 30, 10)
	}
	if cmd == "hide" {
		ui.Clear()
		h.TableInfo = nil
	}
	if cmd == "exit" {
		ui.Close()
		os.Exit(0)
	}
}

func Start[T any]() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Open db
	userDb := mdb.New[T]("./db", "tags", &gosn_driver.Container{})

	cmd := ""

	// Paragraph
	p := widgets.NewParagraph()
	p.Text = "Hello World!"
	w, h := ui.TerminalDimensions()
	p.SetRect(0, h-3, w, h)

	uiInfo := UIInfo[T]{
		DB: userDb,
	}
	uiInfo.Elements = append(uiInfo.Elements, p)

	refresh := func() {
		p.Text = "Cmd: " + cmd
		final := make([]ui.Drawable, 0)
		final = append(final, uiInfo.Elements...)
		if uiInfo.TableInfo != nil {
			final = append(final, uiInfo.TableInfo)
		}
		if uiInfo.PrintInfo != nil {
			final = append(final, uiInfo.PrintInfo)
		}

		for i := 0; i < len(final); i++ {
			ui.Render(final[i])
		}
	}
	refresh()

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<C-c>" {
				break
			}
			if e.ID == "<Enter>" {
				HandleCommand(&uiInfo, cmd)
				cmd = ""
			} else if e.ID == "<Space>" {
				cmd += " "
			} else if e.ID == "<C-<Backspace>>" {
				if len(cmd) > 0 {
					cmd = cmd[0 : len(cmd)-1]
				}
			} else {
				cmd += e.ID
			}

			refresh()
		}
		if e.Type == ui.ResizeEvent {
			w, h := ui.TerminalDimensions()
			p.SetRect(0, h-3, w, h)
			refresh()
		}
	}
}
