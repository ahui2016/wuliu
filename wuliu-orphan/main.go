package main

import (
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
)

func main() {
	util.MustInWuliu()
	findOrphans()
}

func findOrphans() {
	fileOrphans, metaOrphans := lo.Must2(util.FindOrphans())
	fmt.Println("file-orphans:")
	util.PrintList(fileOrphans)
	fmt.Println()
	fmt.Println("metadata-orphans:")
	util.PrintList(metaOrphans)
}
