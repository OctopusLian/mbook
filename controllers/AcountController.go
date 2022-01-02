package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"ziyoubiancheng/mbook/common"
	"ziyoubiancheng/mbook/models"
	"ziyoubiancheng/mbook/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
)

type AccountController struct {
	BaseController
}

var cpt *captcha.Captcha

func init() {
	// use beego cache system store the captcha data
	fc := &cache.FileCache{CachePath: "./cache/captcha"}
	cpt = captcha.NewWithFilter("/captcha/", fc)
}

//注册
func (c *AccountController) Regist() {
	var (
		nickname  string      //昵称
		avatar    string      //头像的http链接地址
		email     string      //邮箱地址
		username  string      //用户名
		id        interface{} //用户id
		captchaOn bool        //是否开启了验证码
	)

	//如果开启了验证码
	if v, ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v, "true") {
		captchaOn = true
		c.Data["CaptchaOn"] = captchaOn
	}

	c.Data["Nickname"] = nickname
	c.Data["Avatar"] = avatar
	c.Data["Email"] = email
	c.Data["Username"] = username
	c.Data["Id"] = id
	c.Data["RandomStr"] = time.Now().Unix()
	c.SetSession("auth", fmt.Sprintf("%v-%v", "email", id)) //存储标识，以标记是哪个用户，在完善用户信息的时候跟传递过来的auth和id进行校验
	c.TplName = "account/bind.html"

}

//登录
func (c *AccountController) Login() {
	var remember CookieRemember
	//验证cookie
	if cookie, ok := c.GetSecureCookie(common.AppKey(), "login"); ok {
		if err := utils.Decode(cookie, &remember); err == nil {
			if err = c.login(remember.MemberId); err == nil {
				c.Redirect(beego.URLFor("HomeController.Index"), 302)
				return
			}
		}
	}
	c.TplName = "account/login.html"

	if c.Ctx.Input.IsPost() {
		account := c.GetString("account")
		password := c.GetString("password")
		member, err := models.NewMember().Login(account, password)
		fmt.Println(err)
		if err != nil {
			c.JsonResult(1, "登录失败", nil)
		}
		member.LastLoginTime = time.Now()
		member.Update()
		c.SetMember(*member)
		remember.MemberId = member.MemberId
		remember.Account = member.Account
		remember.Time = time.Now()
		v, err := utils.Encode(remember)
		if err == nil {
			c.SetSecureCookie(common.AppKey(), "login", v, 24*3600*365)
		}
		c.JsonResult(0, "ok")
	}

	c.Data["RandomStr"] = time.Now().Unix()
}

//注册
func (c *AccountController) DoRegist() {
	var err error
	account := c.GetString("account")
	nickname := strings.TrimSpace(c.GetString("nickname"))
	password1 := c.GetString("password1")
	password2 := c.GetString("password2")
	email := c.GetString("email")

	member := models.NewMember()

	if password1 != password2 {
		c.JsonResult(1, "登录密码与确认密码不一致")
	}

	if l := strings.Count(password1, ""); password1 == "" || l > 20 || l < 6 {
		c.JsonResult(1, "密码必须在6-20个字符之间")
	}

	if ok, err := regexp.MatchString(common.RegexpEmail, email); !ok || err != nil || email == "" {
		c.JsonResult(1, "邮箱格式错误")
	}
	if l := strings.Count(nickname, "") - 1; l < 2 || l > 20 {
		c.JsonResult(1, "用户昵称限制在2-20个字符")
	}

	member.Account = account
	member.Nickname = nickname
	member.Password = password1
	if account == "admin" || account == "administrator" {
		member.Role = common.MemberSuperRole
	} else {
		member.Role = common.MemberGeneralRole
	}
	member.Avatar = common.DefaultAvatar()
	member.CreateAt = 0
	member.Email = email
	member.Status = 0
	if err := member.Add(); err != nil {
		beego.Error(err)
		c.JsonResult(1, err.Error())
	}

	if err = c.login(member.MemberId); err != nil {
		beego.Error(err.Error())
		c.JsonResult(1, err.Error())
	}

	c.JsonResult(0, "注册成功")
}

//退出登录
func (c *AccountController) Logout() {
	c.SetMember(models.Member{})
	c.SetSecureCookie(common.AppKey(), "login", "", -3600)
	c.Redirect(beego.URLFor("AccountController.Login"), 302)
}

/*
* 私有函数
 */
//封装一个内部调用的函数，login
func (c *AccountController) login(memberId int) (err error) {
	member, err := models.NewMember().Find(memberId)
	if member.MemberId == 0 {
		return errors.New("用户不存在")
	}
	//如果没有数据
	if err != nil {
		return err
	}
	member.LastLoginTime = time.Now()
	member.Update()
	c.SetMember(*member)
	var remember CookieRemember
	remember.MemberId = member.MemberId
	remember.Account = member.Account
	remember.Time = time.Now()
	v, err := utils.Encode(remember)
	if err == nil {
		c.SetSecureCookie(common.AppKey(), "login", v, 24*3600*365)
	}
	return err
}
