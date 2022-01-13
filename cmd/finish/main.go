package main

import (
	"autoEncode"
	"flag"
	"fmt"
	"log"
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
		fmt.Printf("%v-%v\n", autoEncode.VERSION, autoEncode.REVISION)
		return
	}

	if name == "" {
		log.Fatal("must -f")
	}

	utf8, err := autoEncode.ToShiftJIS(name)
	if err != nil {
		log.Fatal(err)
	}

	factory, err := autoEncode.NewEncodeFactory()
	if err != nil {
		log.Fatal(err)
	}
	if err := factory.Finish(utf8); err != nil {
		log.Fatal(err)
	}
}
