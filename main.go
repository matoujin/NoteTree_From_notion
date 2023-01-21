package main

import (
	"fmt"
)

var blockId = "XXXXXXXXXXX"
var waittime = 50

func main() {

	SBlock(blockId)
	InitPool(waittime)
	fmt.Println("Building tree......\n if Break with 'panic: send on closed channel' Please waittime more!!")
	dir := Headblock(blockId)
	fmt.Println("Done")
	fmt.Println("Creating file......")
	f := initFile()
	MarkdownFactory(dir, f)

}
