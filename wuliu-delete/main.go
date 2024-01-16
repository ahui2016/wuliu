package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"os"
)

var (
	newFlag = flag.String("newjson", false, "really do add files")
)

func main() {
	flag.Parse()
}