package main

import (
	"autoEncode"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	isVersion                              bool
	inDir, outDir, exeDir, tmpDir, profile string
)

func init() {
	flag.StringVar(&inDir, "i", "", "エンコード対象が格納されたフォルダ")
	flag.StringVar(&outDir, "o", "", "エンコード後のファイルを格納するフォルダ")
	flag.StringVar(&tmpDir, "t", "", "処理対象ファイルを格納するローカルフォルダ")
	flag.StringVar(&exeDir, "e", "", "amatsukazeのexeが格納されたフォルダ")
	flag.StringVar(&profile, "p", "", "amatsukazeで使用するプロファイル名")
	flag.BoolVar(&isVersion, "v", false, "バージョン表示")

	dir := "work"
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	fp, _ := os.Create(filepath.Join(dir, "error.log"))
	writer := io.MultiWriter(os.Stderr, fp)
	autoEncode.Logger = log.New(writer, "INFO", log.LstdFlags|log.Lshortfile)
}

func main() {
	flag.Parse()
	if isVersion {
		fmt.Printf("%v-%v\n", autoEncode.VERSION, autoEncode.REVISION)
		return
	}

	const (
		waitTime = 1 * time.Hour
	)

	factory, err := autoEncode.NewEncodeFactory()
	if err != nil {
		log.Fatal(err)
	}

	if err := factory.Set(inDir, outDir, tmpDir, exeDir, profile); err != nil {
		log.Fatal(err)
	}

	for {
		if err := factory.Add(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(waitTime)
	}
}
