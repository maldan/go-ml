package mbackup

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	Id        string `json:"id"`
	IsVirtual bool   `json:"isVirtual"`

	Src         []string `json:"src"`
	SrcVirtual  []string `json:"srcVirtual"`
	RemoveQueue []string `json:"removeQueue"`
	Dst         []string `json:"dst"`

	Status string `json:"status"`
	Error  string `json:"error"`
	Period string `json:"period"`

	NextRun   time.Time              `json:"nextRun"`
	LastRun   time.Time              `json:"lastRun"`
	BeforeRun func(task *Task) error `json:"-"`
	AfterRun  func(task *Task) error `json:"-"`
	OnError   func(task *Task)       `json:"-"`
}

func (t *Task) Exec(args ...string) error {
	so, se, err := Exec(args...)
	if so != "" {
		fmt.Printf("[TASK EXEC STDOUT] - %v\n", so)
	}
	if se != "" {
		fmt.Printf("[TASK EXEC STDERR] - %v\n", se)
		return errors.New(se)
	}
	if err != nil {
		fmt.Printf("[TASK EXEC ERR] - %v\n", err)
		return err
	}
	return nil
}

func (t *Task) DoRsync() {
	t.Status = "progress"

	sources := t.Src
	if t.IsVirtual {
		sources = t.SrcVirtual
	}

	for _, source := range sources {
		destinations := t.GetDestination()
		for _, destination := range destinations {
			fmt.Printf("[BACKUP TASK RSYNC] From: %v | To: %v\n", source, destination)
			so, se, err := Exec("rsync", "-raz", source, destination)

			if so != "" {
				fmt.Printf("[BACKUP STDOUT] - %v\n", so)
			}
			if se != "" {
				t.Status = "error"
				t.Error = source + "\n" + se
				fmt.Printf("[BACKUP STDERR] - %v\n", se)
			}
			if err != nil {
				t.Status = "error"
				t.Error += "\n" + err.Error()
				fmt.Printf("[BACKUP ERR] - %v\n", err)
			}

			// Break after error
			if t.Status == "error" {
				break
			}
		}
	}
}

func (t *Task) Start() {
	t.Status = "start"
	t.Error = ""
}

func (t *Task) Done() {
	t.Status = "done"
}

func (t *Task) Clean() {
	for _, name := range t.RemoveQueue {
		if name == "" {
			continue
		}

		// Security reason, we can clean only file in tmp folder
		if !strings.HasPrefix(name, "/tmp/") {
			continue
		}

		fmt.Printf("[CLEAN TASK] %v\n", name)

		// Open file
		file, err := os.Open(name)
		if err != nil {
			continue
		}

		fmt.Printf("[CLEAN TASK OPEN] %v\n", name)

		fileInfo, err := file.Stat()
		isDir := false
		if err != nil {
			file.Close()
			continue
		} else {
			isDir = fileInfo.IsDir()
			file.Close()
		}

		fmt.Printf("[CLEAN TASK STATS] %v\n", name)

		if isDir {
			err2 := os.RemoveAll(name)
			if err2 != nil {
				t.Status = "error"
				t.Error = err2.Error()
				continue
			}
		} else {
			err2 := os.Remove(name)
			if err2 != nil {
				t.Status = "error"
				t.Error = err2.Error()
				continue
			}
		}

		// Remove file
		fmt.Printf("[BACKUP TASK CLEAN] Remove: %v\n", name)
	}

	t.RemoveQueue = []string{}
	t.SrcVirtual = []string{}
}

func (t *Task) CreateTempDir() (string, error) {
	// Create temp dir
	tmpDir := fmt.Sprintf("%v/tmp_dir_%v/", os.TempDir(), time.Now().UnixNano())
	err2 := os.MkdirAll(tmpDir, 0777)
	if err2 != nil {
		return "", err2
	}
	return tmpDir, nil
}

func (t *Task) CopyFilesToTmp(from string, ignore []string) (string, error) {
	absFrom, err := filepath.Abs(from)
	if err != nil {
		return "", err
	}

	// Get recursively all files
	list, err := FSListAll(from)
	if err != nil {
		return "", err
	}

	// Filter all files
	list = FilterBy(list, func(t *FileInfo) bool {
		if Includes(ignore, t.Name) {
			return false
		}
		return true
	})

	// Create temp dir
	tmpDir := fmt.Sprintf("%v/tmp_dir_%v/", os.TempDir(), time.Now().UnixNano())
	err2 := os.MkdirAll(tmpDir, 0777)
	if err2 != nil {
		return "", err2
	}

	// Copy files from list to temp dir
	for _, file := range list {
		// dir
		rel := path.Dir(strings.Replace(file.FullPath, absFrom, "", 1))

		// create relative
		os.MkdirAll(tmpDir+"/"+rel, 0777)

		err3 := t.Exec("cp", file.FullPath, tmpDir+"/"+rel)
		if err3 != nil {
			return "", err3
		}
	}

	return tmpDir, nil
}

func (t *Task) CopyFilesToTmpFn(from string, filter func(t *FileInfo) bool) (string, error) {
	absFrom, err := filepath.Abs(from)
	if err != nil {
		return "", err
	}

	// Get recursively all files
	list, err := FSListAll(from)
	if err != nil {
		return "", err
	}

	// Filter all files
	list = FilterBy(list, filter)

	// Create temp dir
	tmpDir := fmt.Sprintf("%v/tmp_dir_%v/", os.TempDir(), time.Now().UnixNano())
	err2 := os.MkdirAll(tmpDir, 0777)
	if err2 != nil {
		return "", err2
	}

	// Copy files from list to temp dir
	for _, file := range list {
		// dir
		rel := path.Dir(strings.Replace(file.FullPath, absFrom, "", 1))

		// create relative
		os.MkdirAll(tmpDir+"/"+rel, 0777)

		err3 := t.Exec("cp", file.FullPath, tmpDir+"/"+rel)
		if err3 != nil {
			return "", err3
		}
	}

	return tmpDir, nil
}

func (t *Task) CompressDir(from string, archivePrefix string) (string, error) {
	// Compress temp folder
	name := fmt.Sprintf("%v/%v_%v.tar.gz", os.TempDir(), SafeName(archivePrefix), time.Now().Format("2006-01-02_15_04_05"))
	err3 := t.Exec("tar", "-czf", name, "-C", from, ".")
	if err3 != nil {
		return "", err3
	}
	return name, nil
}

func (t *Task) EncryptFile(from string, password string) (string, error) {
	err := t.Exec("openssl", "enc", "-aes-256-cbc", "-salt", "-in", from, "-out", from+".enc", "-k", password, "-pbkdf2")
	if err != nil {
		return "", err
	}
	return from + ".enc", nil
}

func (t *Task) DecryptFile(from string, password string) (string, error) {
	if !strings.HasSuffix(from, ".enc") {
		return "", errors.New("not ends with .enc extension")
	}
	to := from[0 : len(from)-4]

	err := t.Exec("openssl", "enc", "-d", "-aes-256-cbc", "-in", from, "-out", to, "-k", password, "-pbkdf2")
	if err != nil {
		return "", err
	}
	return to, nil
}

func (t *Task) GetDestination() []string {
	out := make([]string, 0)
	for i := 0; i < len(t.Dst); i++ {
		dst := strings.ReplaceAll(t.Dst[i], "%date%", time.Now().Format("2006-01-02"))
		out = append(out, dst)
	}
	return out
}

func (t *Task) IsReady() bool {
	return time.Now().Unix() >= t.NextRun.Unix()
}
