package main

import (
	"fmt"
)

var blockId = "XXXXXXXXXXXXXXXX"

func main() {

	SBlock(blockId)
	InitPool(50)
	fmt.Println("Building tree......\n if Break with 'panic: send on closed channel' Please waittime more!!")
	dir := Headblock(blockId)
	fmt.Println("Done")
	fmt.Println("Creating file......")

	f := initFile()
	MarkdownFactory(dir, f)

}
