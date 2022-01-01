package html2text

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 将HTML转成符合elasticsearch搜索的文本
func Html2Text(htmlStr string) string {
	var tags = []string{
		"</p>", "</div>", "</code>", "</span>", "</pre>", "</blockquote>",
		"</h1>", "</h2>", "</h3>", "</h4>", "</h5>", "</h6>", "</td>", "</th>",
		"</i>", "</b>", "</strong>", "</a>", "</li>",
	}
	for _, tag := range tags {
		htmlStr = strings.Replace(htmlStr, tag, tag+" ", -1)
	}

	htmlStr = strings.Replace(htmlStr, "\n", " ", -1)

	gq, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}
	return gq.Text()
}
