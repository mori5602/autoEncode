package autoEncode

import (
	"path/filepath"
	"testing"
)

func TestHash(t *testing.T) {
	path := filepath.Join("testdata", "sample.txt")
	hash, err := Hash(path)
	if err != nil {
		t.Errorf("err:%v", err)
	}

	t.Log(hash)
}
