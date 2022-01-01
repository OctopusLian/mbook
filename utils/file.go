package utils

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
)

func FormatBytes(size int64) string {
	units := []string{" B", " KB", " MB", " GB", " TB"}

	s := float64(size)

	i := 0

	for ; s >= 1024 && i < 4; i++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f%s", s, units[i])
}

//从md的html文件中提取文章标题（从h1-h6）
func ParseTitleFromMdHtml(html string) (title string) {
	h := map[string]string{
		"h1": "h1",
		"H1": "H1",
		"h2": "h2",
		"H2": "H2",
		"h3": "h3",
		"H3": "H3",
		"h4": "h4",
		"H4": "H4",
		"h5": "h5",
		"h6": "h6",
	}
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(html)); err == nil {
		for _, tag := range h {
			if title := doc.Find(tag).First().Text(); len(title) > 0 {
				return title
			}
		}
	} else {
		beego.Error(err.Error())
	}
	return "空标题文档"
}
