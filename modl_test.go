package autoEncode

import (
	"path/filepath"
	"testing"
)

func TestEncodeStatuses_ReadFile(t *testing.T) {
	path := filepath.Join("testdata", "encode_status.csv")
	statuses := NewEncodeStatuses()
	if err := statuses.ReadFile(path); err != nil {
		t.Fatal(err)
	}

	t.Log(statuses)
}

func TestEncodeStatuses_WriteFile(t *testing.T) {
	path := filepath.Join("testdata", "encode_status.csv")
	statuses := NewEncodeStatuses()
	if err := statuses.ReadFile(path); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join("testdata", "encode_status_out.csv")
	if err := statuses.WriteFile(outPath); err != nil {
		t.Fatal(err)
	}
}
