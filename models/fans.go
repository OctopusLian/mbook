package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Fans struct {
	Id       int //PK
	MemberId int
	FansId   int `orm:"index"` //粉丝id
}

type FansData struct {
	MemberId int
	Nickname string
	Avatar   string
	Account  string
}

func (m *Fans) TableName() string {
	return TNFans()
}

// 多字段唯一键
func (m *Fans) TableUnique() [][]string {
	return [][]string{
		[]string{"MemberId", "FansId"},
	}
}

//查询粉丝
func (m *Fans) FansList(mid, page, pageSize int) (fans []FansData, total int64, err error) {
	o := GetOrm("uar")
	total, _ = o.QueryTable(TNFans()).Filter("member_id", mid).Count() //用户粉丝总数
	if total > 0 {
		// sql := fmt.Sprintf(
		// 	"select m.member_id member_id,m.avatar,m.account,m.nickname from "+TNMembers()+" m left join "+TNFans()+" f on m.member_id=f.fans_id where f.member_id=?  order by f.id desc limit %v offset %v",
		// 	pageSize, (page-1)*pageSize,
		// )
		// _, err = o.Raw(sql, mid).QueryRows(&fans)

		sql := "select fans_id from md_fans where member_id=? order by id desc limit %v offset %v"
		sql = fmt.Sprintf(sql, pageSize, (page-1)*pageSize)
		fmt.Println(sql)
		var tmpFans []Fans
		if _, err = o.Raw(sql, mid).QueryRows(&tmpFans); nil == err {
			fansIds := []string{}
			for _, v := range tmpFans {
				fansIds = append(fansIds, strconv.Itoa(v.FansId))
			}
			fansIdStr := strings.Join(fansIds, ",")
			sql = "select member_id,account,avatar,nickname from md_members where member_id in(" + fansIdStr + ")"
			_, err = GetOrm("r").Raw(sql).QueryRows(&fans)
		}

	}
	return
}

//查询关注的人
func (m *Fans) FollowList(fansId, page, pageSize int) (fans []FansData, total int64, err error) {
	o := GetOrm("uar")
	total, _ = o.QueryTable(TNFans()).Filter("fans_id", fansId).Count() //关注总数
	if total > 0 {
		// sql := fmt.Sprintf(
		// 	"select m.member_id member_id,m.avatar,m.account,m.nickname from "+TNMembers()+" m left join "+TNFans()+" f on m.member_id=f.member_id where f.fans_id=?  order by f.id desc limit %v offset %v",
		// 	pageSize, (page-1)*pageSize,
		// )
		// _, err = o.Raw(sql, fansId).QueryRows(&fans)

		sql := "select member_id from md_fans where fans_id=? order by id desc limit %v offset %v"
		sql = fmt.Sprintf(sql, pageSize, (page-1)*pageSize)
		fmt.Println(sql)
		var tmpFans []Fans
		if _, err = o.Raw(sql, fansId).QueryRows(&tmpFans); nil == err {
			memberIds := []string{}
			for _, v := range tmpFans {
				memberIds = append(memberIds, strconv.Itoa(v.MemberId))
			}
			memberIdStr := strings.Join(memberIds, ",")
			sql = "select member_id,account,avatar,nickname from md_members where member_id in(" + memberIdStr + ")"
			_, err = GetOrm("r").Raw(sql).QueryRows(&fans)
		}
	}
	return
}

//查询是否存在关注关系
func (m *Fans) Relation(mid, fansId interface{}) (ok bool) {
	var fans Fans
	GetOrm("uar").QueryTable(TNFans()).Filter("member_id", mid).Filter("fans_id", fansId).One(&fans)
	return fans.Id != 0
}

//关注或取消关注
func (m *Fans) FollowOrCancel(mid, fansId int) (cancel bool, err error) {
	var fans Fans
	o := GetOrm("uaw")
	qs := o.QueryTable(TNFans()).Filter("member_id", mid).Filter("fans_id", fansId)
	qs.One(&fans)
	if fans.Id > 0 { //取消关注
		_, err = qs.Delete()
		cancel = true
	} else { //关注
		fans.MemberId = mid
		fans.FansId = fansId
		_, err = o.Insert(&fans)
	}
	return
}
