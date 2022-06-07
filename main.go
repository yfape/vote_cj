package main

import (
	"fmt"
	"vote_cj/app"
	_ "vote_cj/config"
)

func main() {
	fmt.Println("开始允许")
	app.Run()
}