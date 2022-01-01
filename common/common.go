package common

import (
	"strings"

	"github.com/astaxie/beego"
)

// session
const SessionName = "__mbook_session__"

//正则表达式
const RegexpEmail = `^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`

// 默认PageSize
const PageSize = 20
const RollPage = 4

const WorkingDirectory = "./"

// 用户权限
const (
	// 超级管理员.
	MemberSuperRole = 0
	//普通管理员.
	MemberAdminRole = 1
	//普通用户.
	MemberGeneralRole = 2
)

func Role(role int) string {
	if role == MemberSuperRole {
		return "超级管理员"
	} else if role == MemberAdminRole {
		return "管理员"
	} else if role == MemberGeneralRole {
		return "普通用户"
	} else {
		return ""
	}
}

//图书关系
const (
	// 创始人.
	BookFounder = 0
	//管理
	BookAdmin = 1
	//编辑
	BookEditor = 2
	//普通用户
	BookGeneral = 3
)

func BookRole(role int) string {
	switch role {
	case BookFounder:
		return "创始人"
	case BookAdmin:
		return "管理员"
	case BookEditor:
		return "编辑"
	case BookGeneral:
		return "普通用户"
	default:
		return ""
	}

}

// app_key
func AppKey() string {
	return beego.AppConfig.DefaultString("app_key", "godoc")
}

//默认头像
func DefaultAvatar() string {
	return beego.AppConfig.DefaultString("avatar", "/static/images/headimgurl.jpg")
}

//默认封面
func DefaultCover() string {
	return beego.AppConfig.DefaultString("cover", "/static/images/book.jpg")
}

//获取文件类型
func getFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")
	temp := strings.Split(ext, "|")
	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

//是否允许该类文件类型
func IsAllowedFileExt(ext string) bool {

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := getFileExt()

	for _, item := range exts {
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}
