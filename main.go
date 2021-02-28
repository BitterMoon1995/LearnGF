package main

import (
	_ "LearnGF/boot"
	_ "LearnGF/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
