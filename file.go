package autoEncode

import (
	"fmt"
	"io"
	"os"
)

func checkIsFile(str string) error {
	file, err := os.Stat(str)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%v:%v", ErrFileNotFound, str)
		}
		return fmt.Errorf("%v:%v", ErrException, str)
	}

	if file.IsDir() {
		return fmt.Errorf("target is dir:%v", str)
	}
	return nil
}

func Copy(src, dst string) error {
	if err := checkIsFile(src); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
