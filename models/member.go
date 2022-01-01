package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
	" mbook/common"
	" mbook/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Member struct {
	MemberId      int       `orm:"pk;auto" json:"member_id"`
	Account       string    `orm:"size(30);unique" json:"account"`
	Nickname      string    `orm:"size(30);unique" json:"nickname"`
	Password      string    ` json:"-"`
	Description   string    `orm:"size(640)" json:"description"`
	Email         string    `orm:"size(100);unique" json:"email"`
	Phone         string    `orm:"size(20);null;default(null)" json:"phone"`
	Avatar        string    `json:"avatar"`
	Role          int       `orm:"default(1)" json:"role"`
	RoleName      string    `orm:"-" json:"role_name"`
	Status        int       `orm:"default(0)" json:"status"`
	CreateTime    time.Time `orm:"type(datetime);auto_now_add" json:"create_time"`
	CreateAt      int       `json:"create_at"`
	LastLoginTime time.Time `orm:"type(datetime);null" json:"last_login_time"`
}

func (m *Member) TableName() string {
	return TNMembers()
}

func NewMember() *Member {
	return &Member{}
}

// 添加用
func (m *Member) Add() error {
	if m.Email == "" {
		return errors.New("请填写邮箱")
	}
	if ok, err := regexp.MatchString(common.RegexpEmail, m.Email); !ok || err != nil || m.Email == "" {
		return errors.New("邮箱格式错误")
	}
	if l := strings.Count(m.Password, ""); l < 6 || l >= 20 {
		return errors.New("密码请输入6-20个字符")
	}

	cond := orm.NewCondition().Or("email", m.Email).Or("nickname", m.Nickname).Or("account", m.Account)
	var one Member
	o := GetOrm("w")
	if o.QueryTable(m.TableName()).SetCond(cond).One(&one, "member_id", "nickname", "account", "email"); one.MemberId > 0 {
		if one.Nickname == m.Nickname {
			return errors.New("昵称已存在")
		}
		if one.Email == m.Email {
			return errors.New("邮箱已存在")
		}
		if one.Account == m.Account {
			return errors.New("用户已存在")
		}
	}

	hash, err := utils.PasswordHash(m.Password)

	if err != nil {
		return err
	}

	m.Password = hash
	_, err = o.Insert(m)

	if err != nil {
		return err
	}
	m.RoleName = common.Role(m.Role)
	return nil
}

func (m *Member) Update(cols ...string) error {
	if m.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if _, err := GetOrm("w").Update(m, cols...); err != nil {
		return err
	}
	return nil
}

func (m *Member) Find(id int) (*Member, error) {
	m.MemberId = id
	if err := GetOrm("s").Read(m); err != nil {
		return m, err
	}
	m.RoleName = common.Role(m.Role)
	return m, nil
}

//登录
func (m *Member) Login(account string, password string) (*Member, error) {
	member := &Member{}
	err := GetOrm("r").QueryTable(m.TableName()).Filter("account", account).Filter("status", 0).One(member)

	if err != nil {
		return member, errors.New("用户不存在")
	}

	ok, err := utils.PasswordVerify(member.Password, password)
	if ok && err == nil {
		m.RoleName = common.Role(m.Role)
		return member, nil
	}

	return member, errors.New("密码错误")
}

func (m *Member) IsAdministrator() bool {
	if m == nil || m.MemberId <= 0 {
		return false
	}
	return m.Role == 0 || m.Role == 1
}

//获取用户名
func (m *Member) GetUsernameByUid(id interface{}) string {
	var user Member
	GetOrm("r").QueryTable(TNMembers()).Filter("member_id", id).One(&user, "account")
	return user.Account
}

//获取昵称
func (m *Member) GetNicknameByUid(id interface{}) string {
	var user Member
	if err := GetOrm("s").QueryTable(TNMembers()).Filter("member_id", id).One(&user, "nickname"); err != nil {
		beego.Error(err.Error())
	}

	return user.Nickname
}

//根据用户名获取用户信息
func (m *Member) GetByUsername(username string) (member Member, err error) {
	err = GetOrm("r").QueryTable(TNMembers()).Filter("account", username).One(&member)
	return
}
