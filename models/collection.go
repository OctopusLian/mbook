package models

import (
	"fmt"
	"strconv"
	"strings"
)

type CollectionData struct {
	BookId      int    `json:"book_id"`
	BookName    string `json:"book_name"`
	Identify    string `json:"identify"`
	Description string `json:"description"`
	DocCount    int    `json:"doc_count"`
	Cover       string `json:"cover"`
	MemberId    int    `json:"member_id"`
	Nickname    string `json:"user_name"`
	Vcnt        int    `json:"vcnt"`
	Collection  int    `json:"star"`
	Score       int    `json:"score"`
	CntComment  int    `json:"cnt_comment"`
	CntScore    int    `json:"cnt_score"`
	ScoreFloat  string `json:"score_float"`
	OrderIndex  int    `json:"order_index"`
}

type Collection struct {
	Id       int
	MemberId int `orm:"index"`
	BookId   int
}

func (m *Collection) TableName() string {
	return TNCollection()
}

// 多字段唯一键
func (m *Collection) TableUnique() [][]string {
	return [][]string{
		[]string{"MemberId", "BookId"},
	}
}

//收藏或取消收藏
//@param            uid         用户id
//@param            bid         书籍id
//@return           cancel      是否是取消收藏
func (m *Collection) Collection(uid, bid int) (cancel bool, err error) {
	var star = Collection{MemberId: uid, BookId: bid}
	o := GetOrm("uaw")
	qs := o.QueryTable(TNCollection())
	o.Read(&star, "MemberId", "BookId")
	if star.Id > 0 { //取消收藏
		if _, err = qs.Filter("id", star.Id).Delete(); err == nil {
			IncOrDec(TNBook(), "star", fmt.Sprintf("book_id=%v and star>0", bid), false, 1)
		}
		cancel = true
	} else { //添加收藏
		cancel = false
		if _, err = o.Insert(&star); err == nil {
			//收藏计数+1
			IncOrDec(TNBook(), "star", "book_id="+strconv.Itoa(bid), true, 1)
		}
	}
	return
}

//是否收藏了文档
func (m *Collection) DoesCollection(uid, bid interface{}) bool {
	var star Collection
	star.MemberId, _ = strconv.Atoi(fmt.Sprintf("%v", uid))
	star.BookId, _ = strconv.Atoi(fmt.Sprintf("%v", bid))
	GetOrm("uar").Read(&star, "MemberId", "BookId")
	if star.Id > 0 {
		return true
	}
	return false
}

//获取收藏列表，查询图书信息
func (m *Collection) List(mid, p, listRows int) (cnt int64, books []CollectionData, err error) {
	o := GetOrm("uar")
	filter := o.QueryTable(TNCollection()).Filter("member_id", mid)
	if cnt, _ = filter.Count(); cnt > 0 {
		// sql := "select b.*,m.nickname from " + TNBook() + " b left join " + TNCollection() + " s on s.book_id=b.book_id left join " + TNMembers() + " m on m.member_id=b.member_id where s.member_id=? order by id desc limit %v offset %v"
		// sql = fmt.Sprintf(sql, listRows, (p-1)*listRows)
		// _, err = o.Raw(sql, mid).QueryRows(&books)

		sql := "select bid from " + TNCollection() + " where member_id=? order by id desc limit %v offset %v"
		sql = fmt.Sprintf(sql, listRows, (p-1)*listRows)
		var stars []Collection
		_, err = o.Raw(sql, mid).QueryRows(&stars)
		if nil == err {
			bids := []string{}
			for _, v := range stars {
				bids = append(bids, strconv.Itoa(v.BookId))
			}
			bidstr := strings.Join(bids, ",")

			sql = "select b.*,m.nickname from md_books b left join md_members m on m.member_id=b.member_id where b.book_id in (" + bidstr + ")"
			_, err = GetOrm("r").Raw(sql).QueryRows(&books)
		}
	}
	return
}
