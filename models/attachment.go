package models

import (
	"time"
)

type AttachmentData struct {
	Attachment
	IsExist       bool
	BookName      string
	DocumentName  string
	FileShortSize string
	Account       string
	LocalHttpPath string
}

type Attachment struct {
	AttachmentId int `orm:"pk;auto" json:"attachment_id"`
	BookId       int ` json:"book_id"`
	DocumentId   int ` json:"doc_id"`
	Name         string
	Path         string    `orm:"size(2000)" json:"file_path"`
	Size         float64   `orm:"type(float)" json:"file_size"`
	Ext          string    `orm:"size(50)" json:"file_ext"`
	HttpPath     string    `orm:"size(2000)" json:"http_path"`
	CreateTime   time.Time `orm:"type(datetime);auto_now_add" json:"create_time"`
	CreateAt     int       `orm:"type(int)" json:"create_at"`
}

//orm回调TableName
func (m *Attachment) TableName() string {
	return TNAttachment()
}

func NewAttachment() *Attachment {
	return &Attachment{}
}

func (m *Attachment) Insert() error {
	_, err := GetOrm("w").Insert(m)
	return err
}

func (m *Attachment) Update() error {
	_, err := GetOrm("w").Update(m)
	return err
}

func (m *Attachment) SelectByDocumentId(docId int) (attaches []*Attachment, err error) {
	_, err = GetOrm("r").QueryTable(m.TableName()).Filter("document_id", docId).OrderBy("-attachment_id").All(&attaches)
	return
}
