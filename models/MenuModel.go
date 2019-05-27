package models

import (
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

func MenuStruct() map[int]MenuTree {
	query := orm.NewOrm().QueryTable(TbNameMenu())
	data := make([]*MenuModel, 0)
	query.OrderBy("parent", "-seq").Limit(1000).All(&data)

	var menu = make(map[int]MenuTree)
	for _, v := range data { //查询出来的数组
		//fmt.Println(v.Mid, v.Parent, v.Name)
		if 0 == v.Parent {
			var tree = new(MenuTree)
			tree.MenuModel = *v
			menu[v.Mid] = *tree
		} else {
			if tmp, ok := menu[v.Parent]; ok {
				tmp.Child = append(tmp.Child, *v)
				menu[v.Parent] = tmp
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
