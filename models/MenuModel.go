package models

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/astaxie/beego/orm"
)

type MenuModel struct {
	Mid    int `orm:"pk;auto"`
	Parent int
	Seq    int
	Name   string `orm:"size(45)"`
	Format string `orm:"size(2048);default({})"`
}

type MenuTree struct {
	MenuModel
	Child []MenuModel
}

func (m *MenuModel) TableName() string {
	return TbNameMenu()
}

func MenuTreeStruct(user UserModel) map[int]MenuTree {
	query := orm.NewOrm().QueryTable(TbNameMenu())
	data := make([]*MenuModel, 0)
	query.OrderBy("parent", "-seq").Limit(1000).All(&data)

	var menu = make(map[int]MenuTree)
	//auth
	if len(user.AuthStr) > 0 {
		var authArr []int
		json.Unmarshal([]byte(user.AuthStr), &authArr)
		sort.Ints(authArr)

		for _, v := range data { //查询出来的数组
			//fmt.Println(v.Mid, v.Parent, v.Name)
			if 0 == v.Parent {
				idx := sort.SearchInts(authArr, v.Mid)
				found := (idx < len(authArr) && authArr[idx] == v.Mid)
				if found {
					var tree = new(MenuTree)
					tree.MenuModel = *v
					menu[v.Mid] = *tree
				}
			} else {
				if tmp, ok := menu[v.Parent]; ok {
					tmp.Child = append(tmp.Child, *v)
					menu[v.Parent] = tmp
				}
			}
		}
	}

	return menu
}

func MenuList() ([]*MenuModel, int64) {
	query := orm.NewOrm().QueryTable(TbNameMenu())
	total, _ := query.Count()
	data := make([]*MenuModel, 0)
	query.OrderBy("parent", "-seq").Limit(1000).All(&data)

	return data, total
}

func ParentMenuList() []*MenuModel {
	query := orm.NewOrm().QueryTable(TbNameMenu()).Filter("parent", 0)
	data := make([]*MenuModel, 0)
	query.OrderBy("-seq").Limit(1000).All(&data)

	return data
}

func MenuStruct(mid int) {
	o := orm.NewOrm()
	menu := MenuModel{Mid: mid}
	o.Read(&menu)
	fmt.Println(menu)

	var jsonmap map[string]interface{}
	err := json.Unmarshal([]byte(menu.Format), &jsonmap)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(jsonmap)
	for k, v := range jsonmap {
		fmt.Print(k + ":")
		fmt.Println(v)
	}
}
