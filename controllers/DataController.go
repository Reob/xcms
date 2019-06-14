package controllers

import (
	"fmt"
	"strconv"

	//"github.com/astaxie/beego/orm"
	//"github.com/bitly/go-simplejson"
	"github.com/ziyoubiancheng/xcms/consts"
	"github.com/ziyoubiancheng/xcms/models"
)

type DataController struct {
	BaseController
	Mid int
}

func (c *DataController) Prepare() {
	c.BaseController.Prepare()

	midstr := c.Ctx.Input.Param(":mid")
	mid, err := strconv.Atoi(midstr)
	c.Data["Mid"] = midstr
	if nil != err || mid <= 0 {
		//TODO error page
		c.setTpl()
	}
	c.Mid = mid
}

func (c *DataController) Index() {
	sj := models.MenuFormatStruct(c.Mid)
	//fmt.Println(sj.Get("schema"))
	if sj != nil {
		title := make(map[string]string)
		titlemap := sj.Get("schema")
		for k, _ := range titlemap.MustMap() {
			stype := titlemap.GetPath(k, "type").MustString()
			fmt.Println(k)
			fmt.Println(stype)
			if "object" != stype && "array" != stype {
				title[k] = titlemap.GetPath(k, "title").MustString()
			}
		}

		fmt.Println(title)
		c.Data["Title"] = title
	}

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "data/footerjs.html"
	c.setTpl()

}

func (c *DataController) List() {

	c.listJsonResult(consts.JRCodeFailed, "nil", 0, nil)
}

func (c *DataController) Add() {
	c.initForm()

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "data/footerjs_add.html"
	c.setTpl("data/add.html", "common/layout_jfedit.html")
}
func (c *DataController) AddDo() {
	fmt.Println(c.Mid)
	fmt.Println("+++++++++")
	fmt.Println(string(c.Ctx.Input.RequestBody))
	fmt.Println("---------")
	c.listJsonResult(consts.JRCodeFailed, "nil", 0, nil)
}

func (c *DataController) Edit() {
	c.initForm()

}
func (c *DataController) EditDo() {

}

func (c *DataController) DeleteDo() {

}

func (c *DataController) initForm() {
	format := models.MenuFormatStruct(c.Mid)
	schemaMap := format.Get("schema")
	formArray := format.Get("form")

	//添加通用Form
	fa := formArray.MustArray()
	if len(fa) <= 0 {
		var tmpArray []map[string]string
		tmpArray = append(tmpArray, map[string]string{"key": "parent"})
		tmpArray = append(tmpArray, map[string]string{"key": "name"})
		tmpArray = append(tmpArray, map[string]string{"key": "seq"})
		tmpArray = append(tmpArray, map[string]string{"key": "status"})
		for k, _ := range schemaMap.MustMap() {
			tmpArray = append(tmpArray, map[string]string{"key": k})
		}
		tmpArray = append(tmpArray, map[string]string{"type": "submit", "title": "提交"})

		c.Data["Form"] = tmpArray
	} else {
		var tmpArray []interface{}
		tmpArray = append(tmpArray, map[string]string{"key": "parent"})
		tmpArray = append(tmpArray, map[string]string{"key": "name"})
		tmpArray = append(tmpArray, map[string]string{"key": "seq"})
		tmpArray = append(tmpArray, map[string]string{"key": "status"})
		var haveSubmit bool = false
		for _, v := range formArray.MustArray() {
			tmpArray = append(tmpArray, v)
			if "submit" == v["type"] {
				haveSubmit = true
			}
		}
		if false == haveSubmit {
			tmpArray = append(tmpArray, map[string]string{"type": "submit", "title": "提交"})
		}
		c.Data["Form"] = tmpArray
	}

	//添加通用Schema
	schemaMap.SetPath([]string{"parent", "type"}, "integer")
	schemaMap.SetPath([]string{"parent", "title"}, "上级数据")

	schemaMap.SetPath([]string{"name", "type"}, "string")
	schemaMap.SetPath([]string{"name", "title"}, "名称")

	schemaMap.SetPath([]string{"seq", "type"}, "integer")
	schemaMap.SetPath([]string{"seq", "title"}, "排序(倒序)")

	schemaMap.SetPath([]string{"status", "type"}, "integer")
	schemaMap.SetPath([]string{"status", "title"}, "状态")
	schemaMap.SetPath([]string{"status", "enum"}, []int{0, 1})

	c.Data["Schema"] = schemaMap.MustMap()

	fmt.Println(schemaMap)
}
