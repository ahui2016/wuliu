package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
)

var (
	checkFlag = flag.Bool("check", false, "check orphans")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	if *checkFlag {
		findOrphans()
		return
	}

	flag.Usage()
}

func findOrphans() {
	fileOrphans, metaOrphans := lo.Must2(util.FindOrphans())
	fmt.Print("file-orphans:")
	if len(fileOrphans) == 0 {
		fmt.Print(" (none)")
	}
	fmt.Println()
	util.PrintList(fileOrphans)

	fmt.Print("metadata-orphans:")
	if len(metaOrphans) == 0 {
		fmt.Print(" (none)")
	}
	fmt.Println()
	util.PrintListWithSuffix(metaOrphans, ".json")
}
