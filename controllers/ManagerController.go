package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"ziyoubiancheng/mbook/models"
	"ziyoubiancheng/mbook/utils"
	"ziyoubiancheng/mbook/utils/store"

	"github.com/astaxie/beego"
)

type ManagerController struct {
	BaseController
}

func (c *ManagerController) Prepare() {
	c.BaseController.Prepare()
	if !c.Member.IsAdministrator() {
		c.Abort("404")
	}
}

//分类管理
func (c *ManagerController) Category() {
	cate := new(models.Category)
	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		//新增分类
		pid, _ := c.GetInt("pid")
		if err := cate.InsertMulti(pid, c.GetString("cates")); err != nil {
			c.JsonResult(1, "新增失败："+err.Error())
		}
		c.JsonResult(0, "新增成功")
	}

	//查询所有分类
	cates, err := cate.GetCates(-1, -1)
	if err != nil {
		beego.Error(err)
	}

	var parents []models.Category
	for idx, item := range cates {
		if strings.TrimSpace(item.Icon) == "" {
			item.Icon = "/static/images/icon.png"
		} else {
			item.Icon = utils.ShowImg(item.Icon)
		}
		if item.Pid == 0 {
			parents = append(parents, item)
		}
		cates[idx] = item
	}

	c.Data["Parents"] = parents
	c.Data["Cates"] = cates
	c.Data["IsCategory"] = true
	c.TplName = "manager/category.html"
}

//更新分类字段内容
func (c *ManagerController) UpdateCate() {
	field := c.GetString("field")
	val := c.GetString("value")
	id, _ := c.GetInt("id")
	if err := new(models.Category).UpdateField(id, field, val); err != nil {
		c.JsonResult(1, "更新失败："+err.Error())
	}
	c.JsonResult(0, "更新成功")
}

//删除分类
func (c *ManagerController) DelCate() {
	var err error
	if id, _ := c.GetInt("id"); id > 0 {
		err = new(models.Category).Delete(id)
	}
	if err != nil {
		c.JsonResult(1, err.Error())
	}
	c.JsonResult(0, "删除成功")
}

//更新分类的图标
func (c *ManagerController) UpdateCateIcon() {
	var err error
	id, _ := c.GetInt("id")
	if id == 0 {
		c.JsonResult(1, "参数不正确")
	}
	category := new(models.Category)
	if cate := category.Find(id); cate.Id > 0 {
		cate.Icon = strings.TrimLeft(cate.Icon, "/")
		f, h, err1 := c.GetFile("icon")
		if err1 != nil {
			err = err1
		}
		defer f.Close()

		tmpFile := fmt.Sprintf("uploads/icons/%v%v"+filepath.Ext(h.Filename), id, time.Now().Unix())
		os.MkdirAll(filepath.Dir(tmpFile), os.ModePerm)
		if err = c.SaveToFile("icon", tmpFile); err == nil {
			store.DeleteLocalFiles(cate.Icon)
			err = category.UpdateField(cate.Id, "icon", "/"+tmpFile)
		}
	}

	if err != nil {
		c.JsonResult(1, err.Error())
	}
	c.JsonResult(0, "更新成功")
}
