package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"ziyoubiancheng/mbook/common"
	"ziyoubiancheng/mbook/models"
	"ziyoubiancheng/mbook/utils/store"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DocumentController struct {
	BaseController
}

//获取图书内容并判断权限
func (c *DocumentController) getBookData(identify, token string) *models.BookData {
	book, err := models.NewBook().Select("identify", identify)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}

	//私有文档
	if book.PrivatelyOwned == 1 && !c.Member.IsAdministrator() {
		isOk := false
		if c.Member != nil {
			_, err := models.NewRelationship().SelectRoleId(book.BookId, c.Member.MemberId)
			if err == nil {
				isOk = true
			}
		}
		if book.PrivateToken != "" && !isOk {
			if token != "" && strings.EqualFold(token, book.PrivateToken) {
				c.SetSession(identify, token)
			} else if token, ok := c.GetSession(identify).(string); !ok || !strings.EqualFold(token, book.PrivateToken) {
				c.Abort("404")
			}
		} else if !isOk {
			c.Abort("404")
		}
	}

	bookResult := book.ToBookData()
	if c.Member != nil {
		rsh, err := models.NewRelationship().Select(bookResult.BookId, c.Member.MemberId)
		if err == nil {
			bookResult.MemberId = rsh.MemberId
			bookResult.RoleId = rsh.RoleId
			bookResult.RelationshipId = rsh.RelationshipId
		}
	}
	return bookResult
}

//图书目录&详情页
func (c *DocumentController) Index() {
	token := c.GetString("token")
	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}
	tab := strings.ToLower(c.GetString("tab"))

	bookResult := c.getBookData(identify, token)
	if bookResult.BookId == 0 { //没有阅读权限
		c.Redirect(beego.URLFor("HomeController.Index"), 302)
		return
	}

	c.TplName = "document/intro.html"
	c.Data["Book"] = bookResult

	switch tab {
	case "comment", "score":
	default:
		tab = "default"
	}
	c.Data["Tab"] = tab
	c.Data["Menu"], _ = new(models.Document).GetMenuTop(bookResult.BookId)

	c.Data["Comments"], _ = new(models.Comments).BookComments(1, 30, bookResult.BookId)
	c.Data["MyScore"] = new(models.Score).BookScoreByUid(c.Member.MemberId, bookResult.BookId)
}

//阅读器页面
func (c *DocumentController) Read() {
	identify := c.Ctx.Input.Param(":key")
	id := c.GetString(":id")
	token := c.GetString("token")

	if identify == "" || id == "" {
		c.Abort("404")
	}

	//没开启匿名
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}

	bookData := c.getBookData(identify, token)

	doc := models.NewDocument()
	doc, err := doc.SelectByIdentify(bookData.BookId, id) //文档标识
	if err != nil {
		c.Abort("404")
	}

	if doc.BookId != bookData.BookId {
		c.Abort("404")
	}

	if doc.Release != "" {
		query, err := goquery.NewDocumentFromReader(bytes.NewBufferString(doc.Release))
		if err != nil {
			beego.Error(err)
		} else {
			query.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
				if _, ok := contentSelection.Attr("src"); ok {
				}
				if alt, _ := contentSelection.Attr("alt"); alt == "" {
					contentSelection.SetAttr("alt", doc.DocumentName+" - 图"+fmt.Sprint(i+1))
				}
			})
			html, err := query.Find("body").Html()
			if err != nil {
				beego.Error(err)
			} else {
				doc.Release = html
			}
		}
	}

	attach, err := models.NewAttachment().SelectByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	//图书阅读人次+1
	if err := models.IncOrDec(models.TNBook(), "vcnt",
		fmt.Sprintf("book_id=%v", doc.BookId),
		true, 1,
	); err != nil {
		beego.Error(err.Error())
	}

	//文档阅读人次+1
	if err := models.IncOrDec(models.TNDocuments(), "vcnt",
		fmt.Sprintf("document_id=%v", doc.DocumentId),
		true, 1,
	); err != nil {
		beego.Error(err.Error())
	}

	doc.Vcnt = doc.Vcnt + 1

	if c.IsAjax() {
		var data struct {
			Id        int    `json:"doc_id"`
			DocTitle  string `json:"doc_title"`
			Body      string `json:"body"`
			Title     string `json:"title"`
			View      int    `json:"view"`
			UpdatedAt string `json:"updated_at"`
		}
		data.DocTitle = doc.DocumentName
		data.Body = doc.Release
		data.Id = doc.DocumentId
		data.View = doc.Vcnt
		data.UpdatedAt = doc.ModifyTime.Format("2006-01-02 15:04:05")

		c.JsonResult(0, "ok", data)
	}

	tree, err := models.NewDocument().GetMenuHtml(bookData.BookId, doc.DocumentId)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}

	c.Data["Bookmark"] = false
	c.Data["Model"] = bookData
	c.Data["Book"] = bookData
	c.Data["Result"] = template.HTML(tree)
	c.Data["Title"] = doc.DocumentName
	c.Data["DocId"] = doc.DocumentId
	c.Data["Content"] = template.HTML(doc.Release)
	c.Data["View"] = doc.Vcnt
	c.Data["UpdatedAt"] = doc.ModifyTime.Format("2006-01-02 15:04:05")

	//设置模版
	c.TplName = "document/default_read.html"
}

//编辑
func (c *DocumentController) Edit() {
	docId := 0 // 文档id

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}

	bookData := models.NewBookData()

	var err error
	//权限验证
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			c.JsonResult(1, "权限错误")
		}
		bookData = book.ToBookData()
	} else {
		bookData, err = models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)
		if err != nil {
			c.Abort("404")
		}

		if bookData.RoleId == common.BookGeneral {
			c.JsonResult(1, "权限错误")
		}
	}

	c.TplName = "document/markdown_edit_template.html"

	c.Data["Model"] = bookData
	r, _ := json.Marshal(bookData)

	c.Data["ModelResult"] = template.JS(string(r))

	c.Data["Result"] = template.JS("[]")

	// 编辑的文档
	if id := c.GetString(":id"); id != "" {
		if num, _ := strconv.Atoi(id); num > 0 {
			docId = num
		} else { //字符串
			var doc = models.NewDocument()
			models.GetOrm("w").QueryTable(doc).Filter("identify", id).Filter("book_id", bookData.BookId).One(doc, "document_id")
			docId = doc.DocumentId
		}
	}

	trees, err := models.NewDocument().GetMenu(bookData.BookId, docId, true)
	if err != nil {
		beego.Error("GetMenu error : ", err)
	} else {
		if len(trees) > 0 {
			if jsTree, err := json.Marshal(trees); err == nil {
				c.Data["Result"] = template.JS(string(jsTree))
			}
		} else {
			c.Data["Result"] = template.JS("[]")
		}
	}
	c.Data["BaiDuMapKey"] = beego.AppConfig.DefaultString("baidumapkey", "")

}

//创建文档
func (c *DocumentController) Create() {
	identify := c.GetString("identify")        //图书标识
	docIdentify := c.GetString("doc_identify") //新建的文档标识
	docName := c.GetString("doc_name")
	parentId, _ := c.GetInt("parent_id", 0)
	docId, _ := c.GetInt("doc_id", 0)
	bookIdentify := strings.TrimSpace(c.GetString(":key"))
	o := models.GetOrm("w")
	if identify == "" {
		c.JsonResult(1, "参数错误")
	}
	if docName == "" {
		c.JsonResult(1, "文档名为空")
	}
	if docIdentify != "" {
		if bookIdentify == "" {
			c.JsonResult(1, "图书参数错误")
		}

		var book models.Book
		o.QueryTable(models.TNBook()).Filter("Identify", bookIdentify).One(&book, "BookId")
		if book.BookId == 0 {
			c.JsonResult(1, "未找到该图书")
		}

		d, _ := models.NewDocument().SelectByIdentify(book.BookId, docIdentify)
		if d.DocumentId > 0 && d.DocumentId != docId {
			c.JsonResult(1, "文档标识重复")
		}
	} else {
		docIdentify = fmt.Sprintf("date-%v", time.Now().Format("2019.11.02.01.01.05"))
	}

	bookId := 0
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			beego.Error(err)
			c.JsonResult(1, "权限错误")
		}
		bookId = book.BookId
	} else {
		bookData, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)

		if err != nil || bookData.RoleId == common.BookGeneral {
			c.JsonResult(1, "权限错误")
		}
		bookId = bookData.BookId
	}

	if parentId > 0 {
		doc, err := models.NewDocument().SelectByDocId(parentId)
		if err != nil || doc.BookId != bookId {
			c.JsonResult(1, "分类错误")
		}
	}

	document, _ := models.NewDocument().SelectByDocId(docId)

	document.MemberId = c.Member.MemberId
	document.BookId = bookId
	if docIdentify != "" {
		document.Identify = docIdentify
	}
	document.Version = time.Now().Unix()
	document.DocumentName = docName
	document.ParentId = parentId

	documentId, err := document.InsertOrUpdate()
	if err != nil {
		c.JsonResult(1, "保存失败")
	}

	documentStore := models.DocumentStore{DocumentId: int(documentId), Markdown: ""}
	if documentStore.SelectField(documentId, "markdown") == "" {
		if err := documentStore.InsertOrUpdate(); err != nil {
			beego.Error(err)
		}
	}
	c.JsonResult(0, "ok", document)
}

//上传附件
func (c *DocumentController) Upload() {
	identify := c.GetString("identify")
	docId, _ := c.GetInt("doc_id")
	isAttach := true

	if identify == "" {
		c.JsonResult(1, "参数错误")
	}
	name := "editormd-file-file"
	file, moreFile, err := c.GetFile(name)
	if err == http.ErrMissingFile {
		name = "editormd-image-file"
		file, moreFile, err = c.GetFile(name)
		if err == http.ErrMissingFile {
			c.JsonResult(1, "文件错误")
		}
	}
	if err != nil {
		c.JsonResult(1, err.Error())
	}

	defer file.Close()

	ext := filepath.Ext(moreFile.Filename)
	if ext == "" {
		c.JsonResult(1, "文件格式错误")
	}

	if !common.IsAllowedFileExt(ext) {
		c.JsonResult(1, "文件类型错误")
	}

	bookId := 0
	//如果是超级管理员，则不判断权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			c.JsonResult(1, "文档不存在或权限不足")
		}
		bookId = book.BookId
	} else {
		book, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)
		if err != nil {
			if err == orm.ErrNoRows {
				c.JsonResult(1, "权限错误")
			}
			c.JsonResult(6001, err.Error())
		}
		//没有编辑权限
		if book.RoleId != common.BookEditor && book.RoleId != common.BookAdmin && book.RoleId != common.BookFounder {
			c.JsonResult(1, "权限错误")
		}
		bookId = book.BookId
	}

	if docId > 0 {
		doc, err := models.NewDocument().SelectByDocId(docId)
		if err != nil {
			c.JsonResult(1, "获取文档错误")
		}
		if doc.BookId != bookId {
			c.JsonResult(1, "获取文档错误")
		}
	}

	fileName := strconv.FormatInt(time.Now().UnixNano(), 16)
	filePath := filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("200601"), fileName+ext)
	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile(name, filePath)

	if err != nil {
		c.JsonResult(1, "保存文件失败")
	}
	attachment := models.NewAttachment()
	attachment.BookId = bookId
	attachment.Name = moreFile.Filename
	attachment.CreateAt = c.Member.MemberId
	attachment.Ext = ext
	attachment.Path = strings.TrimPrefix(filePath, common.WorkingDirectory)
	attachment.DocumentId = docId

	if fileInfo, err := os.Stat(filePath); err == nil {
		attachment.Size = float64(fileInfo.Size())
	}
	if docId > 0 {
		attachment.DocumentId = docId
	}

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") || strings.EqualFold(ext, ".png") || strings.EqualFold(ext, ".gif") {

		attachment.HttpPath = "/" + strings.Replace(strings.TrimPrefix(filePath, common.WorkingDirectory), "\\", "/", -1)
		if strings.HasPrefix(attachment.HttpPath, "//") {
			attachment.HttpPath = string(attachment.HttpPath[1:])
		}
		isAttach = false
	}

	err = attachment.Insert()

	if err != nil {
		os.Remove(filePath)
		c.JsonResult(1, "文件保存失败")
	}
	if attachment.HttpPath == "" {
		attachment.HttpPath = beego.URLFor("DocumentController.DownloadAttachment", ":key", identify, ":attach_id", attachment.AttachmentId)

		if err := attachment.Update(); err != nil {
			c.JsonResult(1, "保存文件失败")
		}
	}
	osspath := fmt.Sprintf("projects/%v/%v", identify, fileName+filepath.Ext(attachment.HttpPath))

	osspath = "uploads/" + osspath
	if err := store.SaveToLocal("."+attachment.HttpPath, osspath); err != nil {
		beego.Error(err.Error())
	}
	attachment.HttpPath = "/" + osspath

	result := map[string]interface{}{
		"errcode":   0,
		"success":   1,
		"message":   "ok",
		"url":       attachment.HttpPath,
		"alt":       attachment.Name,
		"is_attach": isAttach,
		"attach":    attachment,
	}
	c.Ctx.Output.JSON(result, true, false)
	c.StopRun()
}

//删除
func (c *DocumentController) Delete() {

	identify := c.GetString("identify")
	docId, _ := c.GetInt("doc_id", 0)

	bookId := 0
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			c.JsonResult(1, "权限错误")
		}
		bookId = book.BookId
	} else {
		bookData, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)
		if err != nil || bookData.RoleId == common.BookGeneral {
			c.JsonResult(1, "权限错误")
		}
		bookId = bookData.BookId
	}

	if docId <= 0 {
		c.JsonResult(1, "参数错误")
	}

	doc, err := models.NewDocument().SelectByDocId(docId)
	if err != nil {
		c.JsonResult(1, "删除失败")
	}

	//如果文档所属图书错误
	if doc.BookId != bookId {
		c.JsonResult(1, "参数错误")
	}
	//删除图书下的文档以及子文档
	err = doc.Delete(doc.DocumentId)
	if err != nil {
		beego.Error(err.Error())
		c.JsonResult(1, "删除失败")
	}

	//文档数量统计
	models.NewBook().RefreshDocumentCount(doc.BookId)

	c.JsonResult(0, "ok")
}

//保存文档并返回内容
func (c *DocumentController) Content() {
	identify := c.Ctx.Input.Param(":key")
	docId, err := c.GetInt("doc_id")
	errMsg := "ok"
	if err != nil {
		docId, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	}
	bookId := 0
	//权限验证
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			c.JsonResult(1, "获取内容错误")
		}
		bookId = book.BookId
	} else {
		bookData, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)

		if err != nil || bookData.RoleId == common.BookGeneral {
			c.JsonResult(1, "权限错误")
		}
		bookId = bookData.BookId
	}

	if docId <= 0 {
		c.JsonResult(1, "参数错误")
	}

	documentStore := new(models.DocumentStore)

	if !c.Ctx.Input.IsPost() {
		doc, err := models.NewDocument().SelectByDocId(docId)

		if err != nil {
			c.JsonResult(1, "文档不存在")
		}
		attach, err := models.NewAttachment().SelectByDocumentId(doc.DocumentId)
		if err == nil {
			doc.AttachList = attach
		}

		doc.Release = "" //Ajax请求，之间用markdown渲染，不用release
		doc.Markdown = documentStore.SelectField(doc.DocumentId, "markdown")
		c.JsonResult(0, errMsg, doc)
	}

	//更新文档内容
	markdown := strings.TrimSpace(c.GetString("markdown", ""))
	content := c.GetString("html")

	version, _ := c.GetInt64("version", 0)
	isCover := c.GetString("cover")

	doc, err := models.NewDocument().SelectByDocId(docId)

	if err != nil {
		c.JsonResult(1, "读取文档错误")
	}
	if doc.BookId != bookId {
		c.JsonResult(1, "内部错误")
	}
	if doc.Version != version && !strings.EqualFold(isCover, "yes") {
		c.JsonResult(1, "文档将被覆盖")
	}

	isSummary := false
	isAuto := false

	if markdown == "" && content != "" {
		documentStore.Markdown = content
	} else {
		documentStore.Markdown = markdown
	}
	documentStore.Content = content
	doc.Version = time.Now().Unix()
	if docId, err := doc.InsertOrUpdate(); err != nil {
		c.JsonResult(1, "保存失败")
	} else {
		documentStore.DocumentId = int(docId)
		if err := documentStore.InsertOrUpdate("markdown", "content"); err != nil {
			beego.Error(err)
		}
	}

	if isAuto {
		errMsg = "auto"
	} else if isSummary {
		errMsg = "true"
	}

	doc.Release = ""
	c.JsonResult(0, errMsg, doc)
}

//阅读页内搜索
func (c *DocumentController) Search() {
	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")
	keyword := strings.TrimSpace(c.GetString("keyword"))

	if identify == "" {
		c.JsonResult(1, "参数错误")
	}
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}
	bookData := c.getBookData(identify, token)
	docs, _, err := models.NewDocumentSearch().SearchDocument(keyword, bookData.BookId, 1, 10000)
	if err != nil {
		beego.Error(err)
		c.JsonResult(1, "搜索结果错误")
	}
	c.JsonResult(0, keyword, docs)
}
