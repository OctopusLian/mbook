package models

//文档编辑
type DocumentStore struct {
	DocumentId int    `orm:"pk;auto;column(document_id)"`
	Markdown   string `orm:"type(text);"` //markdown内容
	Content    string `orm:"type(text);"` //html内容
}

func (m *DocumentStore) TableName() string {
	return TNDocumentStore()
}

//编辑表内容
func (m *DocumentStore) SelectField(docId interface{}, field string) string {
	var ds = DocumentStore{}
	if field != "markdown" {
		field = "content"
	}
	GetOrm("r").QueryTable(TNDocumentStore()).Filter("document_id", docId).One(&ds, field)
	if field == "content" {
		return ds.Content
	}
	return ds.Markdown
}

//插入或者更新
func (m *DocumentStore) InsertOrUpdate(fields ...string) (err error) {
	o := GetOrm("w")
	var one DocumentStore
	o.QueryTable(TNDocumentStore()).Filter("document_id", m.DocumentId).One(&one, "document_id")

	if one.DocumentId > 0 {
		_, err = o.Update(m, fields...)
	} else {
		_, err = o.Insert(m)
	}
	return
}

//删除记录
func (m *DocumentStore) Delete(docId ...interface{}) {
	if len(docId) > 0 {
		GetOrm("w").QueryTable(TNDocumentStore()).Filter("document_id__in", docId...).Delete()
	}
}
