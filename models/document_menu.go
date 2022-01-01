package models

import (
	"bytes"
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
)

//图书目录
type DocumentMenu struct {
	DocumentId   int             `json:"id"`
	DocumentName string          `json:"text"`
	BookIdentify string          `json:"-"`
	Identify     string          `json:"identify"`
	ParentId     interface{}     `json:"parent"`
	Version      int64           `json:"version"`
	State        *highlightState `json:"state,omitempty"` //如果字段为空，则json中不会有该字段
}
type highlightState struct {
	Selected bool `json:"selected"`
	Opened   bool `json:"opened"`
}

func (m *Document) GetMenuHtml(bookId, selectedId int) (string, error) {
	trees, err := m.GetMenu(bookId, selectedId)
	if err != nil {
		return "", err
	}
	parentId := m.highlightNode(trees, selectedId)

	buf := bytes.NewBufferString("")

	m.treeHTML(trees, 0, selectedId, parentId, buf)

	return buf.String(), nil

}

func (m *Document) GetMenu(bookId int, selectedId int, isEdit ...bool) ([]*DocumentMenu, error) {
	trees := make([]*DocumentMenu, 0)
	var docs []*Document

	count, err := GetOrm("r").QueryTable(m).Filter("book_id", bookId).OrderBy("order_sort", "identify").Limit(2000).All(&docs, "document_id", "document_name", "parent_id", "identify", "version")
	if err != nil {
		return trees, err
	}
	book, _ := NewBook().Select("book_id", bookId)
	trees = make([]*DocumentMenu, count)
	for index, item := range docs {
		tree := &DocumentMenu{}
		if selectedId > 0 {
			if selectedId == item.DocumentId {
				tree.State = &highlightState{Selected: true, Opened: true}
			}
		} else {
			if index == 0 {
				tree.State = &highlightState{Selected: true, Opened: true}
			}
		}

		tree.DocumentId = item.DocumentId
		tree.Identify = item.Identify
		tree.Version = item.Version
		tree.BookIdentify = book.Identify
		if item.ParentId > 0 {
			tree.ParentId = item.ParentId
		} else {
			tree.ParentId = "#"
		}
		idf := item.Identify
		if idf == "" {
			idf = strconv.Itoa(item.DocumentId)
		}
		if len(isEdit) > 0 && isEdit[0] == true {
			tree.DocumentName = item.DocumentName + "<small class='text-danger'>(" + idf + ")</small>"
		} else {
			tree.DocumentName = item.DocumentName
		}

		trees[index] = tree
	}

	return trees, nil
}

func (m *Document) highlightNode(array []*DocumentMenu, parentId int) int {
	for _, item := range array {
		if _, ok := item.ParentId.(string); ok && item.DocumentId == parentId {
			return item.DocumentId
		} else if pid, ok := item.ParentId.(int); ok && item.DocumentId == parentId {
			if pid == parentId {
				return 0
			}
			return m.highlightNode(array, pid)
		}
	}
	return 0
}

//生成目录HTML
func (m *Document) treeHTML(array []*DocumentMenu, parentId int, selectedId int, selectedParentId int, buf *bytes.Buffer) {
	buf.WriteString("<ul>")

	for _, item := range array {
		pid := 0

		if p, ok := item.ParentId.(int); ok {
			pid = p
		}
		if pid == parentId {

			selected := ""
			if item.DocumentId == selectedId {
				selected = ` class="jstree-clicked"`
			}
			selectedLi := ""
			if item.DocumentId == selectedParentId {
				selectedLi = ` class="jstree-open"`
			}
			buf.WriteString("<li id=\"")
			buf.WriteString(strconv.Itoa(item.DocumentId))
			buf.WriteString("\"")
			buf.WriteString(selectedLi)
			buf.WriteString("><a href=\"")
			if item.Identify != "" {
				uri := beego.URLFor("DocumentController.Read", ":key", item.BookIdentify, ":id", item.Identify)
				buf.WriteString(uri)
			} else {
				uri := beego.URLFor("DocumentController.Read", ":key", item.BookIdentify, ":id", item.DocumentId)
				buf.WriteString(uri)
			}
			buf.WriteString("\" title=\"")
			buf.WriteString(template.HTMLEscapeString(item.DocumentName) + "\"")
			buf.WriteString(selected + ">")
			buf.WriteString(template.HTMLEscapeString(item.DocumentName) + "</a>")

			for _, sub := range array {
				if p, ok := sub.ParentId.(int); ok && p == item.DocumentId {
					m.treeHTML(array, p, selectedId, selectedParentId, buf)
					break
				}
			}
			buf.WriteString("</li>")
		}
	}
	buf.WriteString("</ul>")
}
