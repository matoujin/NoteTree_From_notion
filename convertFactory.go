package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"sort"
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

type Block struct {
	CreateTime string `json:"create_time" `
	EditTime   string `json:"edit_time" `
	Title      string `json:"title" `
}

/*type notes struct {
	Title            string
	has_children     bool
	node_id          string
	createTime       string
	last_edited_time string
	parent_note_id   string
}*/

var BlockList = make([]Block, 0)

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

	block := Block{
		Title:      Notelist[blockId].Title,
		CreateTime: Notelist[blockId].CreateTime[0:7],
		EditTime:   Notelist[blockId].Last_edited_time[0:7],
	}
	BlockList = append(BlockList, block)

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
	CalenFactory()

	return MaptoJson(dirBlockchain)
}

func MaptoJson(dirBlockchain Dirblock) string {
	//fmt.Println("Decoding---------------")
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

func CalenFactory() {

	filePath := "calen.md"
	calenFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	datelist := make([]string, 0)
	datelist = append(datelist, BlockList[0].CreateTime)
	//create datalist
	for _, value := range BlockList {
		flag := false
	loop1:
		for _, date := range datelist {
			if date == value.CreateTime {
				fmt.Println("New date added...")
				flag = true
				break loop1
			}
		}
		if flag == false {
			datelist = append(datelist, value.CreateTime)
		}
	}

	sort.Strings(datelist)

	//year
	yearlist := make([]string, 0)
	yearlist = append(yearlist, datelist[0][0:4])
	for _, value := range datelist {
		flag := false
	loop2:
		for _, year := range yearlist {
			if value[0:4] == year {
				flag = true
				break loop2
			}
		}
		if flag == false {
			yearlist = append(yearlist, value[0:4])
		}
	}

	for _, year := range yearlist {
		yearcount := 0
		yearString := ""
		monthString := ""
		for _, date1 := range datelist {
			if date1[0:4] == year {
				monthcount := 0
				contentString := ""
				for _, block := range BlockList {
					if block.CreateTime == date1 {
						contentString = contentString +
							"\t" + block.Title + " " + "修改时间" + block.EditTime + "\n\n"
						monthcount++
					}
				}
				yearcount += monthcount
				monthString = monthString + "### " + date1 + " 共(" + strconv.Itoa(monthcount) + ")篇" + "\n\n" + contentString
			}
		}
		yearString = "## " + year + " 共(" + strconv.Itoa(yearcount) + ")篇" + "\n\n" + monthString + "\n\n---\n\n"
		calenFile.WriteString(yearString)
	}

	calenFile.Close()

}
