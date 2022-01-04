package autoEncode

import (
	"bytes"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"strings"
)

func transformEncoding(rawReader io.Reader, trans transform.Transformer) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(rawReader, trans))
	if err == nil {
		return string(ret), nil
	} else {
		return "", err
	}
}

// FromShiftJIS Convert a string encoding from ShiftJIS to UTF-8
func FromShiftJIS(str string) (string, error) {
	return transformEncoding(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
}

// ToShiftJIS Convert a string encoding from UTF-8 to ShiftJIS
func ToShiftJIS(str string) (string, error) {
	return transformEncoding(strings.NewReader(str), japanese.ShiftJIS.NewEncoder())
}

// BytesFromShiftJIS Convert an array of bytes (a valid ShiftJIS string) to a UTF-8 string
func BytesFromShiftJIS(b []byte) (string, error) {
	return transformEncoding(bytes.NewReader(b), japanese.ShiftJIS.NewDecoder())
}

// BytesToShiftJIS Convert an array of bytes (a valid UTF-8 string) to a ShiftJIS string
func BytesToShiftJIS(b []byte) (string, error) {
	return transformEncoding(bytes.NewReader(b), japanese.ShiftJIS.NewEncoder())
}
