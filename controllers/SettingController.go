package controllers

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"ziyoubiancheng/mbook/common"
	"ziyoubiancheng/mbook/models"
	"ziyoubiancheng/mbook/utils/graphics"
	"ziyoubiancheng/mbook/utils/store"

	"github.com/astaxie/beego"
)

type SettingController struct {
	BaseController
}

//基本信息
func (c *SettingController) Index() {
	if c.Ctx.Input.IsPost() {
		email := strings.TrimSpace(c.GetString("email", ""))
		phone := strings.TrimSpace(c.GetString("phone"))
		description := strings.TrimSpace(c.GetString("description"))
		if email == "" {
			c.JsonResult(1, "邮箱不能为空")
		}
		member := c.Member
		member.Email = email
		member.Phone = phone
		member.Description = description
		if err := member.Update("email", "phone", "description"); err != nil {
			c.JsonResult(1, "提交信息错误")
		}
		c.SetMember(*member)
		c.JsonResult(0, "ok")
	}
	c.Data["SettingBasic"] = true
	c.TplName = "setting/index.html"
}

//上传头像
func (c *SettingController) Upload() {
	file, moreFile, err := c.GetFile("image-file")
	if err != nil {
		c.JsonResult(1, "文件异常")
	}
	defer file.Close()

	ext := filepath.Ext(moreFile.Filename)
	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		c.JsonResult(1, "图片格式异常")
	}

	x1, _ := strconv.ParseFloat(c.GetString("x"), 10)
	y1, _ := strconv.ParseFloat(c.GetString("y"), 10)
	w1, _ := strconv.ParseFloat(c.GetString("width"), 10)
	h1, _ := strconv.ParseFloat(c.GetString("height"), 10)

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	fileName := strconv.FormatInt(time.Now().UnixNano(), 16)
	filePath := filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("201911"), fileName+ext)

	path := filepath.Dir(filePath)
	os.MkdirAll(path, os.ModePerm)
	err = c.SaveToFile("image-file", filePath)
	if err != nil {
		c.JsonResult(1, "保存失败")
	}

	//剪切图片
	subImg, err := graphics.ImageCopyFromFile(filePath, x, y, width, height)
	if err != nil {
		c.JsonResult(1, "剪切失败")
	}
	os.Remove(filePath)

	filePath = filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("201911"), fileName+ext)
	graphics.ImageResizeSaveFile(subImg, 120, 120, filePath)
	err = graphics.SaveImage(filePath, subImg)
	if err != nil {
		c.JsonResult(1, "保存文件失败")
	}

	url := "/" + strings.Replace(strings.TrimPrefix(filePath, common.WorkingDirectory), "\\", "/", -1)
	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}

	if member, err := models.NewMember().Find(c.Member.MemberId); err == nil {
		avatar := member.Avatar
		member.Avatar = url
		err = member.Update()
		if err != nil {
			c.JsonResult(1, "保存信息失败")
		}
		if strings.HasPrefix(avatar, "/uploads/") {
			os.Remove(filepath.Join(common.WorkingDirectory, avatar))
		}
		c.SetMember(*member)
	}

	if err := store.SaveToLocal("."+url, strings.TrimLeft(url, "./")); err != nil {
		beego.Error(err.Error())
	} else {
		url = "/" + strings.TrimLeft(url, "./")
	}

	c.JsonResult(0, "ok", url)
}
