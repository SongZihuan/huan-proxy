package main

import (
	"github.com/SongZihuan/huan-proxy/src/cmd/define"
	"github.com/SongZihuan/huan-proxy/src/mainfunc"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

var v1Main define.MainFunc = mainfunc.MainV1

func main() {
	utils.Exit(_main())
}

func _main() int {
	return v1Main()
}
