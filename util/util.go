package util

import (
	"os"
	"github.com/samber/lo"
)

func GetCwd() string {
	return lo.Must1(os.Getwd())
}
