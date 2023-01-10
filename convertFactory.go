package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strconv"
	"strings"
)

type Dirblock struct {
	Title      string              `json:"title" `
	Url        string              `json:"url" `
	CreateTime string              `json:"create_time" `
	EditTime   string              `json:"edit_time" `
	Child      map[string]Dirblock `json:",omitempy" `
	HasChild   bool                `json:"has_child"`
	Depth      int                 `json:"depth"`
}

/*type notes struct {
	Title            string
	has_children     bool
	node_id          string
	createTime       string
	last_edited_time string
	parent_note_id   string
}*/

func convert(blockId string) Dirblock {
	newBlock := Dirblock{
		Title:      Notelist[blockId].Title,
		Url:        "https://glen-mat-fd9.notion.site/" + blockId,
		CreateTime: Notelist[blockId].CreateTime,
		EditTime:   Notelist[blockId].Last_edited_time,
		Child:      make(map[string]Dirblock),
		HasChild:   Notelist[blockId].Has_children,
		Depth:      1,
	}

	return newBlock
}

func addChild(blockId string, dirBlockchain Dirblock) Dirblock {

	for id, childblock := range Notelist {

		if childblock.Parent_note_id == blockId {
			cvtchild := convert(id)

			cvtchild.Depth = dirBlockchain.Depth + 1
			dirBlockchain.Child[id] = cvtchild
			if dirBlockchain.HasChild == true {
				cvtchild = addChild(id, cvtchild)
			}
		}
	}
	return dirBlockchain
}

func Headblock(blockId string) string {

	dirBlockchain := convert(blockId)
	dirBlockchain = addChild(blockId, dirBlockchain)

	return MaptoJson(dirBlockchain)
}

func MaptoJson(dirBlockchain Dirblock) string {
	fmt.Println("Decoding---------------")
	buf, err := json.MarshalIndent(dirBlockchain, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	return string(buf)
}

func MarkdownFactory(jsons string, f *os.File) {

	f.WriteString(
		strings.Repeat("#", 1) + " " + gjson.Get(jsons, "title").String() + "\n\n" +
			gjson.Get(jsons, "url").String() + "\n\n" +
			"创建时间" + gjson.Get(jsons, "create_time").String() + "\n\n" +
			"修改时间" + gjson.Get(jsons, "edit_time").String() + "\n\n",
	)

	if gjson.Get(jsons, "has_child").Bool() == true {
		cjson := gjson.Get(jsons, "Child")
		cjson.ForEach(func(_, value gjson.Result) bool {
			value1 := value.String()
			f.WriteString(
				"> " +
					strings.Repeat("#", getdepth(value)) + " " + gjson.Get(value1, "title").String() + "\n\n" +
					gjson.Get(value1, "url").String() + "\n\n" +
					"创建时间" + gjson.Get(value1, "create_time").String() + "\n\n" +
					"修改时间" + gjson.Get(value1, "edit_time").String() + "\n\n",
			)
			markdownLoop(value1, f)
			return true
		})

	}
	f.Close()
}

func markdownLoop(child string, f *os.File) {

	if gjson.Get(child, "has_child").Bool() == true {
		cjson := gjson.Get(child, "Child")
		cjson.ForEach(func(_, c gjson.Result) bool {
			c1 := c.String()
			f.WriteString(
				strings.Repeat("#", getdepth(c)) + " " + gjson.Get(c1, "title").String() + "\n\n" +
					gjson.Get(c1, "url").String() + "\n\n" +
					"创建时间" + gjson.Get(c1, "create_time").String() + "\n\n" +
					"修改时间" + gjson.Get(c1, "edit_time").String() + "\n\n",
			)
			markdownLoop(c1, f)
			return true
		})
	}
}
func getdepth(child gjson.Result) int {
	str := gjson.Get(child.String(), "depth").String()
	res, _ := strconv.Atoi(str)
	return res
}
