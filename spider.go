package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"time"
)

type notes struct {
	Title            string
	Has_children     bool
	node_id          string
	CreateTime       string
	Last_edited_time string
	Parent_note_id   string
}

var Notelist = map[string]*notes{}
var token = "secret_XXXXXXXXXXXX"

func SBlockChildren(blockId string) {

	request := gorequest.New()
	res, body, _ := request.
		Get("https://api.notion.com/v1/blocks/"+blockId+"/children?page_size=100").
		Set("Notion-Version", "2022-06-28").
		Set("Authorization", "Bearer "+token).
		Set("accept", "application/json").
		End()

	if res.StatusCode != 200 {
		fmt.Println(blockId, "with ErrorCode", res.StatusCode)
	} else {
		js, _ := simplejson.NewJson([]byte(body))
		results, _ := js.Get("results").Array()
		for _, value := range results {
			childBlock, _ := value.(map[string]interface{})
			//只记录子页面
			cTitleMap := childBlock["child_page"]
			if cTitleMap != nil {
				cTitle := mapGet(cTitleMap)
				cCreatedTime := childBlock["created_time"]
				cHasChildren := childBlock["has_children"]
				cBlockId_ := childBlock["id"]
				cLastEditedTime := childBlock["last_edited_time"]

				cBlockId := cBlockId_.(string)
				Notelist[cBlockId] = new(notes)
				initStruct(Notelist[cBlockId])
				Notelist[cBlockId].Parent_note_id = blockId
				Notelist[cBlockId].Title = cTitle.(string)
				Notelist[cBlockId].node_id = cBlockId
				Notelist[cBlockId].Has_children = cHasChildren.(bool)
				Notelist[cBlockId].Last_edited_time = timeDecoder(cLastEditedTime.(string))
				Notelist[cBlockId].CreateTime = timeDecoder(cCreatedTime.(string))

				if Notelist[cBlockId].Has_children == true {
					go NewTask(cBlockId)
				}
			}
		}
	}
}

func mapGet(tMap interface{}) interface{} {
	titleMap := tMap.(map[string]interface{})
	title, _ := titleMap["title"]
	return title
}

func initStruct(note *notes) {

	note.Parent_note_id = ""
	note.Has_children = false
	note.CreateTime = ""
	note.Title = ""
	note.node_id = ""
	note.Last_edited_time = ""
}

func SBlock(blockId string) {

	Notelist[blockId] = new(notes)
	initStruct(Notelist[blockId])
	Notelist[blockId].node_id = blockId

	request := gorequest.New()
	res, body, _ := request.
		Get("https://api.notion.com/v1/blocks/"+blockId).
		Set("Notion-Version", "2022-06-28").
		Set("Authorization", "Bearer "+token).
		Set("accept", "application/json").
		End()

	if res.StatusCode != 200 {
		fmt.Println(blockId, "with ErrorCode", res.StatusCode)
	} else {
		cTime := gjson.Get(body, "created_time")
		Notelist[blockId].CreateTime = timeDecoder(cTime.String())

		eTime := gjson.Get(body, "last_edited_time")
		Notelist[blockId].Last_edited_time = timeDecoder(eTime.String())

		tit := gjson.Get(body, "child_page.title")
		Notelist[blockId].Title = tit.String()

		parentType := gjson.Get(body, "parent.type").String()
		//!!!database not select into page!!!
		if parentType == "page_id" {
			Notelist[blockId].Parent_note_id = gjson.Get(body, "parent.page_id").String()
		} else if parentType == "block_id" {
			Notelist[blockId].Parent_note_id = gjson.Get(body, "block.page_id").String()
		} else {
			Notelist[blockId].Parent_note_id = "Top_page"
		}

		hasChildren := gjson.Get(body, "has_children").Bool()
		Notelist[blockId].Has_children = hasChildren
		if hasChildren == true {
			go NewTask(blockId)
		}
	}
}

func timeDecoder(tm string) string {
	to, _ := time.Parse("2006-01-02T15:04:05Z", tm)
	tm1 := to.Format("2006-01-02")
	return tm1
}
