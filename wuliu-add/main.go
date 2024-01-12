package main

import (
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
)

func main() {
	util.MustInWuliu()
	findNewFiles()
}

func findNewFiles() {
	names := lo.Must(util.FindNewFiles())
	for _, name := range names {
		fmt.Println(name)
	}
}
