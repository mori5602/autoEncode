package main

import (
	"autoEncode"
	"flag"
	"fmt"
	"log"
	"path/filepath"
)

var (
	isVersion bool
	name      string
)

func init() {
	flag.StringVar(&name, "f", "", "対象ファイル")
	flag.BoolVar(&isVersion, "v", false, "バージョン表示")
}

func main() {
	flag.Parse()
	if isVersion {
		fmt.Printf("%v:%v\n", autoEncode.VERSION, autoEncode.REVISION)
		return
	}

	if name == "" {
		log.Fatal("must -f")
	}

	title := filepath.Base(name)
	factory, err := autoEncode.NewEncodeFactory()
	if err != nil {
		log.Fatal(err)
	}

	if err := factory.Finish(title); err != nil {
		log.Fatal(err)
	}
}
