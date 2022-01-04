package autoEncode

import (
	"fmt"
	"github.com/jszwec/csvutil"
	"io"
	"io/ioutil"
	"os"
)

const (
	Init = iota
	Add
	Added
	Started
	Finish
)

type EncodeState struct {
	Title  string `csv:"Title"`
	Status int    `csv:"Status"`
}

type EncodeStatuses []EncodeState

func NewEncodeStatuses() EncodeStatuses {
	return EncodeStatuses{}
}

func (f *EncodeStatuses) ReadAll(r io.Reader) (err error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := csvutil.Unmarshal(bytes, f); err != nil {
		return err
	}
	return nil
}

func (f *EncodeStatuses) ReadFile(path string) (err error) {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("%v:%v", ErrFileNotFound, path)
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	return f.ReadAll(file)
}

func (f *EncodeStatuses) WriteAll(w io.Writer) (n int, err error) {
	bytes, err := csvutil.Marshal(f)
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (f *EncodeStatuses) WriteFile(path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = f.WriteAll(file)
	if err != nil {
		return err
	}
	return err
}

func (f *EncodeStatuses) Exists(title string) bool {
	for _, record := range *f {
		if record.Title == title {
			return true
		}
	}
	return false
}

func (f *EncodeStatuses) GetStatus(title string) (int, error) {
	for _, record := range *f {
		if record.Title == title {
			return record.Status, nil
		}
	}
	return 0, ErrTargetNotFound
}

func (f *EncodeStatuses) Add(title string) error {
	if f.Exists(title) {
		status, err := f.GetStatus(title)
		if err != nil {
			return err
		}
		if status != Init {
			return fmt.Errorf("%v:%v", ErrAlreadyExists, title)
		}
	}

	*f = append(*f, EncodeState{
		Title:  title,
		Status: Init,
	})
	return nil
}

func (f *EncodeStatuses) Set(title string, i int) error {
	for num, record := range *f {
		if record.Title != title {
			continue
		}

		if record.Status == i {
			return ErrStatusUnchanged
		}

		(*f)[num].Status = i
		return nil
	}
	return ErrZeroRecord
}
