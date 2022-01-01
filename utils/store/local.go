package store

import (
	"os"
	"path/filepath"
	"strings"
)

//删除文件
//@param           object                     文件对象
//@param           IsPreview                  是否是预览的Local
func DeleteLocalFiles(object ...string) error {
	for _, file := range object {
		os.Remove(strings.TrimLeft(file, "/"))
	}
	return nil
}

//保存文件
//@param            tmpfile          临时文件
//@param            save             存储文件，不建议与临时文件相同，特别是IsDel参数值为true的时候
//@param            IsDel            文件上传后，是否删除临时文件
func SaveToLocal(tmpfile, save string) (err error) {
	save = strings.TrimLeft(save, "/")
	//"./a.png"与"a.png"是相同路径
	if strings.HasPrefix(tmpfile, "./") || strings.HasPrefix(save, "./") {
		tmpfile = strings.TrimPrefix(tmpfile, "./")
		save = strings.TrimPrefix(save, "./")
	}
	if strings.ToLower(tmpfile) != strings.ToLower(save) { //不是相同文件路径
		os.MkdirAll(filepath.Dir(save), os.ModePerm)
		err = os.Rename(tmpfile, save)
	}
	return
}
