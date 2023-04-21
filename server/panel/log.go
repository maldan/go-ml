package ms_panel

import (
	"github.com/maldan/go-ml/db/mdb"
	ms_log "github.com/maldan/go-ml/server/log"
	ml_string "github.com/maldan/go-ml/util/string"
	ml_time "github.com/maldan/go-ml/util/time"
)

type Log struct {
	Path string
}

type ArgsOffset struct {
	Page           int    `json:"page"`
	Kind           string `json:"kind"`
	Created        string `json:"created"`
	TimezoneOffset int    `json:"timezoneOffset"`
}

func (l Log) GetList(args ArgsOffset) any {
	fieldList := make([]string, 0)
	whereFn := make([]func(t *ms_log.LogBody) bool, 0)

	if args.Kind != "" {
		fieldList = append(fieldList, "Kind")
		whereFn = append(whereFn, func(t *ms_log.LogBody) bool {
			return t.Kind == args.Kind
		})
	}
	if args.Created != "" {
		fieldList = append(fieldList, "Created")

		// No time
		created := ml_time.DateTime{}
		created.FromString(args.Created)
		created.SetTimezoneOffset(args.TimezoneOffset)
		created = created.In(0)

		// Full
		if len(ml_string.OnlyDigit(args.Created)) > 12 {
			whereFn = append(whereFn, func(t *ms_log.LogBody) bool {
				return t.Created.Equal(created)
			})
		} else
		// Have date and hour and minutes
		if len(ml_string.OnlyDigit(args.Created)) == 12 {
			whereFn = append(whereFn, func(t *ms_log.LogBody) bool {
				dateIn := created.In(t.Created.TimezoneOffset())
				return t.Created.EqualDate(created) && t.Created.Hour() == dateIn.Hour() && t.Created.Minute() == dateIn.Minute()
			})
		} else
		// Only date
		// Have date and hour
		if len(ml_string.OnlyDigit(args.Created)) == 10 {
			whereFn = append(whereFn, func(t *ms_log.LogBody) bool {
				dateIn := created.In(t.Created.TimezoneOffset())
				return t.Created.EqualDate(created) && t.Created.Hour() == dateIn.Hour()
			})
		} else
		// Only date
		if len(ml_string.OnlyDigit(args.Created)) <= 8 {
			whereFn = append(whereFn, func(t *ms_log.LogBody) bool {
				return t.Created.EqualDate(created)
			})
		}
	}

	// Final where
	finalWhere := func(t any) bool {
		for i := 0; i < len(whereFn); i++ {
			if !whereFn[i](t.(*ms_log.LogBody)) {
				return false
			}
		}
		return true
	}

	return ms_log.LogDB.FindBy(mdb.ArgsFind{
		// FieldList: strings.Join(fieldList, ","),
		Where: finalWhere,
		// Limit:     10,
	})

	/*// Open file
	f, err := os.OpenFile(l.Path, os.O_RDONLY, 0777)
	ms_error.FatalIfError(err)

	// Get size
	info, err := f.Stat()
	ms_error.FatalIfError(err)

	blockSize := 1024 * 4
	if int(info.Size()) < blockSize {
		blockSize = int(info.Size())
	}

	// Read lines
	b := make([]byte, blockSize)
	f.ReadAt(b, info.Size()-int64((args.Page)*blockSize))
	ss := string(b)
	lines := strings.Split(ss, "\n")

	return map[string]any{
		"lines": ml_slice.Reverse(lines),
		"total": info.Size() / int64(blockSize),
	}*/
}
