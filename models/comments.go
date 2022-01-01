package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

/*
*
*	评论
*
 */

//评论表
type Comments struct {
	Id         int
	Uid        int       `orm:"index"` //用户id
	BookId     int       `orm:"index"` //文档项目id
	Content    string    //评论内容
	TimeCreate time.Time //评论时间
}

// func (m *Comments) TableName() string {
// 	return TNComments()
// }

//评论内容
type BookCommentsResult struct {
	Uid        int       `json:"uid"`
	Score      int       `json:"score"`
	Avatar     string    `json:"avatar"`
	Nickname   string    `json:"nickname"`
	Content    string    `json:"content"`
	TimeCreate time.Time `json:"time_create"` //评论时间
}

func (m *Comments) AddComments(uid, bookId int, content string) (err error) {
	var comment Comments
	//1.限制评论频率
	second := 10
	sql := `select id from ` + TNComments(bookId) + ` where uid=? and time_create>? order by id desc`

	o := GetOrm("uaw")
	o.Raw(sql, uid, time.Now().Add(-time.Duration(second)*time.Second)).QueryRow(&comment)
	if comment.Id > 0 {
		return errors.New(fmt.Sprintf("您距离上次发表评论时间小于 %v 秒，请歇会儿再发。", second))
	}
	fmt.Println(comment)
	//2.插入评论数据
	sql = `insert into ` + TNComments(bookId) + `(uid,book_id,content,time_create) values(?,?,?,?)`
	_, err = o.Raw(sql, uid, bookId, content, time.Now()).Exec()
	if err != nil {
		beego.Error(err.Error())
		err = errors.New("发表评论失败")
		return
	}

	//3.评论数+1
	sql = `update ` + TNBook() + ` set cnt_comment=cnt_comment+1 where book_id=?`
	o.Raw(sql, bookId)

	return
}

//评论内容
func (m *Comments) BookComments(page, size, bookId int) (comments []BookCommentsResult, err error) {
	// sql := `select c.content,s.score,c.uid,c.time_create,m.avatar,m.nickname from ` + TNComments(bookId) + ` c left join ` + TNMembers() + ` m on m.member_id=c.uid left join ` + TNScore() + ` s on s.uid=c.uid and s.book_id=c.book_id where c.book_id=? order by c.id desc limit %v offset %v`
	// sql = fmt.Sprintf(sql, size, (page-1)*size)
	// _, err = GetOrm("w").Raw(sql, bookId).QueryRows(&comments)
	// return

	o := GetOrm("uar")

	sql := `select book_id,uid,content,time_create from ` + TNComments(bookId) + ` where book_id=? limit %v offset %v`
	sql = fmt.Sprintf(sql, size, (page-1)*size)
	_, err = o.Raw(sql, bookId).QueryRows(&comments)
	if nil != err {
		return
	}

	//头像昵称
	uids := []string{}
	for _, v := range comments {
		uids = append(uids, strconv.Itoa(v.Uid))
	}
	uidstr := strings.Join(uids, ",")
	sql = `select member_id,avatar,nickname from md_members where member_id in(` + uidstr + `)`
	members := []Member{}
	_, err = GetOrm("r").Raw(sql).QueryRows(&members)
	if nil != err {
		return
	}
	memberMap := make(map[int]Member)
	for _, member := range members {
		memberMap[member.MemberId] = member
	}
	for k, v := range comments {
		comments[k].Avatar = memberMap[v.Uid].Avatar
		comments[k].Nickname = memberMap[v.Uid].Nickname
	}

	//评分
	sql = `select uid,score from md_score where book_id=? and uid in(` + uidstr + `)`
	scores := []Score{}
	_, err = o.Raw(sql, bookId).QueryRows(&scores)
	if nil != err {
		return
	}
	scoreMap := make(map[int]Score)
	for _, score := range scores {
		scoreMap[score.Uid] = score
	}
	for k, v := range comments {
		comments[k].Score = scoreMap[v.Uid].Score
	}

	return
}

/*
*
*	评分
*
 */

//评分表
type Score struct {
	Id         int
	BookId     int
	Uid        int
	Score      int //评分
	TimeCreate time.Time
}

func (m *Score) TableName() string {
	return TNScore()
}

// 多字段唯一键
func (m *Score) TableUnique() [][]string {
	return [][]string{
		[]string{"Uid", "BookId"},
	}
}

//评分内容
type BookScoresResult struct {
	Avatar     string    `json:"avatar"`
	Nickname   string    `json:"nickname"`
	Score      string    `json:"score"`
	TimeCreate time.Time `json:"time_create"` //评论时间
}

//获取评分内容
func (m *Score) BookScores(p, listRows, bookId int) (scores []BookScoresResult, err error) {
	sql := `select s.score,s.time_create,m.avatar,m.nickname from ` + TNScore() + ` s left join ` + TNMembers() + ` m on m.member_id=s.uid where s.book_id=? order by s.id desc limit %v offset %v`
	sql = fmt.Sprintf(sql, listRows, (p-1)*listRows)
	_, err = GetOrm("uar").Raw(sql, bookId).QueryRows(&scores)
	return
}

//查询用户对文档的评分
func (m *Score) BookScoreByUid(uid, bookId interface{}) int {
	var score Score
	GetOrm("uar").QueryTable(TNScore()).Filter("uid", uid).Filter("book_id", bookId).One(&score, "score")
	return score.Score
}

//添加评论内容

//添加评分
//score的值只能是1-5，然后需要对scorex10，50则表示5.0分
func (m *Score) AddScore(uid, bookId, score int) (err error) {
	//查询评分是否已存在
	o := GetOrm("w")
	var scoreObj = Score{Uid: uid, BookId: bookId}
	o.Read(&scoreObj, "uid", "book_id")
	if scoreObj.Id > 0 { //评分已存在
		err = errors.New("您已给当前文档打过分了")
		return
	}

	//评分不存在，添加评分记录
	score = score * 10
	scoreObj.Score = score
	scoreObj.TimeCreate = time.Now()
	GetOrm("uaw").Insert(&scoreObj)
	if scoreObj.Id > 0 { //评分添加成功，更行当前书籍项目的评分
		//评分人数+1
		var book = Book{BookId: bookId}
		o.Read(&book, "book_id")
		if book.CntScore == 0 {
			book.CntScore = 1
			book.Score = 0
		} else {
			book.CntScore = book.CntScore + 1
		}
		book.Score = (book.Score*(book.CntScore-1) + score) / book.CntScore
		_, err = o.Update(&book, "cnt_score", "score")
		if err != nil {
			beego.Error(err.Error())
			err = errors.New("评分失败，内部错误")
		}
	}
	return
}
