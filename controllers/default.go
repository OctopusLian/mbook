package controllers

import (
	"fmt"
	"ziyoubiancheng/mbook/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) test() {
	var maps []orm.Params
	models.GetOrm("w").Raw("select * from md_documents").Values(&maps)
	fmt.Println(maps)

	fmt.Println("--------")
	o := models.GetOrm("w")
	o.Using("r")
	o.Raw("select * from md_documents").Values(&maps)
	fmt.Println(maps)

	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) InitData() {
	c.TplName = "index.tpl"
}
