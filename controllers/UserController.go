package controllers

import (
	"ziyoubiancheng/mbook/common"
	"ziyoubiancheng/mbook/models"
	"ziyoubiancheng/mbook/utils"

	"github.com/astaxie/beego"
)

type UserController struct {
	BaseController
	UcenterMember models.Member
}

func (c *UserController) Prepare() {
	c.BaseController.Prepare()

	username := c.GetString(":username")
	c.UcenterMember, _ = new(models.Member).GetByUsername(username)
	if c.UcenterMember.MemberId == 0 {
		c.Abort("404")
		return
	}
	c.Data["IsSelf"] = c.UcenterMember.MemberId == c.Member.MemberId
	c.Data["User"] = c.UcenterMember
	c.Data["Tab"] = "share"
}

//首页
func (c *UserController) Index() {
	page, _ := c.GetInt("page")
	pageSize := 10
	if page < 1 {
		page = 1
	}
	books, totalCount, _ := models.NewBook().SelectPage(page, pageSize, c.UcenterMember.MemberId, 0)
	c.Data["Books"] = books

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, totalCount, pageSize, page, beego.URLFor("UserController.Index", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Total"] = totalCount
	c.TplName = "user/index.html"
}

//收藏
func (c *UserController) Collection() {
	page, _ := c.GetInt("page")
	pageSize := 10
	if page < 1 {
		page = 1
	}

	totalCount, books, _ := new(models.Collection).List(c.UcenterMember.MemberId, page, pageSize)
	c.Data["Books"] = books

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Collection", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Total"] = totalCount
	c.Data["Tab"] = "collection"
	c.TplName = "user/collection.html"
}

//关注
func (c *UserController) Follow() {
	page, _ := c.GetInt("page")
	pageSize := 18
	if page < 1 {
		page = 1
	}
	fans, totalCount, _ := new(models.Fans).FollowList(c.UcenterMember.MemberId, page, pageSize)
	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Follow", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Fans"] = fans
	c.Data["Tab"] = "follow"
	c.TplName = "user/fans.html"
}

//粉丝和关注
func (c *UserController) Fans() {
	page, _ := c.GetInt("page")
	pageSize := 18
	if page < 1 {
		page = 1
	}
	fans, totalCount, _ := new(models.Fans).FansList(c.UcenterMember.MemberId, page, pageSize)
	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Fans", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Fans"] = fans
	c.Data["Tab"] = "fans"
	c.TplName = "user/fans.html"
}
