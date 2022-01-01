package models

import (
	"fmt"
	"mbook/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Book struct {
	BookId         int       `orm:"pk;auto" json:"book_id"`
	BookName       string    `orm:"size(500)" json:"book_name"`       //名称
	Identify       string    `orm:"size(100);unique" json:"identify"` //唯一标识
	OrderIndex     int       `orm:"default(0)" json:"order_index"`
	Description    string    `orm:"size(1000)" json:"description"`       //图书描述
	Cover          string    `orm:"size(1000)" json:"cover"`             //封面地址
	Editor         string    `orm:"size(50)" json:"editor"`              //编辑器类型: "markdown"
	Status         int       `orm:"default(0)" json:"status"`            //状态:0 正常 ; 1 已删除
	PrivatelyOwned int       `orm:"default(0)" json:"privately_owned"`   // 是否私有: 0 公开 ; 1 私有
	PrivateToken   string    `orm:"size(500);null" json:"private_token"` // 私有图书访问Token
	MemberId       int       `orm:"size(100)" json:"member_id"`
	CreateTime     time.Time `orm:"type(datetime);auto_now_add" json:"create_time"` //创建时间
	ModifyTime     time.Time `orm:"type(datetime);auto_now_add" json:"modify_time"`
	ReleaseTime    time.Time `orm:"type(datetime);" json:"release_time"` //发布时间
	DocCount       int       `json:"doc_count"`                          //文档数量
	CommentCount   int       `orm:"type(int)" json:"comment_count"`
	Vcnt           int       `orm:"default(0)" json:"vcnt"`              //阅读次数
	Collection     int       `orm:"column(star);default(0)" json:"star"` //收藏次数
	Score          int       `orm:"default(40)" json:"score"`            //评分
	CntScore       int       //评分人数
	CntComment     int       //评论人数
	Author         string    `orm:"size(50)"`                      //来源
	AuthorURL      string    `orm:"column(author_url);size(1000)"` //来源链接
}

//orm 回调
func (m *Book) TableName() string {
	return TNBook()
}

func NewBook() *Book {
	return &Book{}
}

func (m *Book) HomeData(pageIndex, pageSize int, cid int, fields ...string) (books []Book, totalCount int, err error) {
	if len(fields) == 0 {
		fields = append(fields, "book_id", "book_name", "identify", "cover", "order_index")
	}
	fieldStr := "b." + strings.Join(fields, ",b.")

	sqlFmt := "select %v from " + TNBook() + " b left join " + TNBookCategory() + " c on b.book_id=c.book_id where c.category_id=" + strconv.Itoa(cid)

	sql := fmt.Sprintf(sqlFmt, fieldStr)
	sqlCount := fmt.Sprintf(sqlFmt, "count(*) cnt")
	fmt.Println(sql)
	fmt.Println(sqlCount)
	o := GetOrm("r")
	var params []orm.Params
	if _, err := o.Raw(sqlCount).Values(&params); err == nil {
		if len(params) > 0 {
			totalCount, _ = strconv.Atoi(params[0]["cnt"].(string))
		}
	}
	_, err = o.Raw(sql).QueryRows(&books)

	return
}

func (m *Book) SearchBook(wd string, page, size int) (books []Book, cnt int, err error) {
	sqlFmt := "select %v from md_books where book_name like ? or description like ? order by star desc"
	sql := fmt.Sprintf(sqlFmt, "book_id")
	sqlCount := fmt.Sprintf(sqlFmt, "count(book_id) cnt")

	wd = "%" + wd + "%"

	o := GetOrm("r")
	var count struct{ Cnt int }
	err = o.Raw(sqlCount, wd, wd).QueryRow(&count)
	if count.Cnt > 0 {
		cnt = count.Cnt
		_, err = o.Raw(sql+" limit ? offset ?", wd, wd, size, (page-1)*size).QueryRows(&books)
	}

	return
}

func (m *Book) GetBooksByIds(ids []int, fields ...string) (books []Book, err error) {
	if len(ids) == 0 {
		return
	}

	var bs []Book
	var idArr []interface{}

	for _, i := range ids {
		idArr = append(idArr, i)
	}

	rows, err := GetOrm("r").QueryTable(TNBook()).Filter("book_id__in", idArr).All(&bs, fields...)
	if rows > 0 {
		bookMap := make(map[interface{}]Book)
		for _, book := range bs {
			bookMap[book.BookId] = book
		}
		for _, i := range ids {
			if book, ok := bookMap[i]; ok {
				books = append(books, book)
			}
		}
	}

	return
}

//Insert
func (m *Book) Insert() (err error) {
	if _, err = GetOrm("w").Insert(m); err != nil {
		return
	}

	relationship := Relationship{BookId: m.BookId, MemberId: m.MemberId, RoleId: 0}
	if err = relationship.Insert(); err != nil {
		return err
	}

	document := Document{BookId: m.BookId, DocumentName: "空白文档", Identify: "blank", MemberId: m.MemberId}
	var id int64
	if id, err = document.InsertOrUpdate(); err == nil {
		documentstore := DocumentStore{DocumentId: int(id), Markdown: ""}
		err = documentstore.InsertOrUpdate()
	}
	return err
}

//Update
func (m *Book) Update(cols ...string) (err error) {
	bk := NewBook()
	bk.BookId = m.BookId
	o := GetOrm("w")
	if err = o.Read(bk); err != nil {
		return err
	}
	_, err = o.Update(m, cols...)
	return err
}

func (m *Book) Select(field string, value interface{}, cols ...string) (book *Book, err error) {
	o := GetOrm("r")
	if len(cols) == 0 {
		err = o.QueryTable(m.TableName()).Filter(field, value).One(m)
	} else {
		err = o.QueryTable(m.TableName()).Filter(field, value).One(m, cols...)
	}
	return m, err
}

func (m *Book) SelectPage(pageIndex, pageSize, memberId int, PrivatelyOwned int) (books []*BookData, totalCount int, err error) {
	o := GetOrm("r")
	sql1 := "select count(b.book_id) as total_count from " + TNBook() + " as b left join " +
		TNRelationship() + " as r on b.book_id=r.book_id and r.member_id = ? where r.relationship_id > 0  and b.privately_owned=" + strconv.Itoa(PrivatelyOwned)

	err = o.Raw(sql1, memberId).QueryRow(&totalCount)
	if err != nil {
		return
	}
	offset := (pageIndex - 1) * pageSize
	sql2 := "select book.*,rel.member_id,rel.role_id,m.account as create_name from " + TNBook() + " as book" +
		" left join " + TNRelationship() + " as rel on book.book_id=rel.book_id and rel.member_id = ?" +
		" left join " + TNRelationship() + " as rel1 on book.book_id=rel1.book_id  and rel1.role_id=0" +
		" left join " + TNMembers() + " as m on rel1.member_id=m.member_id " +
		" where rel.relationship_id > 0 %v order by book.book_id desc limit " + fmt.Sprintf("%d,%d", offset, pageSize)
	sql2 = fmt.Sprintf(sql2, " and book.privately_owned="+strconv.Itoa(PrivatelyOwned))

	_, err = o.Raw(sql2, memberId).QueryRows(&books)
	if err != nil {
		return
	}
	return
}

func (book *Book) ToBookData() (m *BookData) {
	m = &BookData{}
	m.BookId = book.BookId
	m.BookName = book.BookName
	m.Identify = book.Identify
	m.OrderIndex = book.OrderIndex
	m.Description = strings.Replace(book.Description, "\r\n", "<br/>", -1)
	m.PrivatelyOwned = book.PrivatelyOwned
	m.PrivateToken = book.PrivateToken
	m.DocCount = book.DocCount
	m.CommentCount = book.CommentCount
	m.CreateTime = book.CreateTime
	m.ModifyTime = book.ModifyTime
	m.Cover = book.Cover
	m.MemberId = book.MemberId
	m.Status = book.Status
	m.Editor = book.Editor
	m.Vcnt = book.Vcnt
	m.Collection = book.Collection
	m.Score = book.Score
	m.ScoreFloat = utils.ScoreFloat(book.Score)
	m.CntScore = book.CntScore
	m.CntComment = book.CntComment
	m.Author = book.Author
	m.AuthorURL = book.AuthorURL
	if book.Editor == "" {
		m.Editor = "markdown"
	}
	return m
}

//更新文档数量
func (m *Book) RefreshDocumentCount(bookId int) {
	o := GetOrm("w")
	docCount, err := o.QueryTable(TNDocuments()).Filter("book_id", bookId).Count()
	if err == nil {
		temp := NewBook()
		temp.BookId = bookId
		temp.DocCount = int(docCount)
		o.Update(temp, "doc_count")
	} else {
		beego.Error(err)
	}
}
