package filelist;

import (
	"testing"
	"html/template"
	"github.com/NyaaPantsu/nyaa/model"
)

func makeDummyFile(path ...string) (file model.File) {
	file.SetPath(path)
	return
}

func dashes(n uint) (out string) {
	var i uint
	for i = 0; i < n; i++ {
		out += "-"
	}
	return
}

func TestFilelist(T *testing.T) {
	files := []model.File{
		makeDummyFile("A", "B", "C.txt"),
		makeDummyFile("A", "C", "C.txt"),
		makeDummyFile("B.txt"),
	}
	expected := "A\n"+
	            "-B\n"+
	            "--C.txt\n"+
	            "-C\n"+
                "--C.txt\n"+
                "B.txt\n"

	filelist := FileListToFolder(files)

	out, err := filelist.MakeFolderTreeView("{{dashes .NestLevel}}{{.FolderName}}\n", "{{dashes .NestLevel}}{{.Filename}}\n", map[string]interface{}{
		"dashes": dashes,
	}, nil)
	if err != nil {
		T.Fatalf("%v", err)
		return
	}
	if out != template.HTML(expected) {
		T.Fatalf("Error: expected %s, got %s", expected, out)
		return
	}
}

