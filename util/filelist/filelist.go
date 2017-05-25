package filelist

import (
	"bytes"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/bradfitz/slice"
	"html/template"
	"strconv"
	"strings"
)

type FileListFolder struct {
	Folders map[string]*FileListFolder
	Files   []model.File
}

func FileListToFolder(fileList []model.File) (out *FileListFolder) {
	out = &FileListFolder{
		Folders: make(map[string]*FileListFolder),
		Files:   make([]model.File, 0),
	}

	pathsToFolders := make(map[string][]model.File)

	for _, file := range fileList {
		pathArray := file.Path()

		if len(pathArray) > 1 {
			pathStrippedFile := model.File{
				ID:           file.ID,
				TorrentID:    file.TorrentID,
				BencodedPath: "",
				Filesize:     file.Filesize,
			}
			pathStrippedFile.SetPath(pathArray[1:])
			pathsToFolders[pathArray[0]] = append(pathsToFolders[pathArray[0]], pathStrippedFile)
		} else {
			out.Files = append(out.Files, file)
		}
	}

	for folderName, folderFiles := range pathsToFolders {
		out.Folders[folderName] = FileListToFolder(folderFiles)
	}

	return
}

func (f *FileListFolder) TotalSize() (out int64) {
	out = 0
	for _, folder := range f.Folders {
		out += folder.TotalSize()
	}

	for _, file := range f.Files {
		out += file.Filesize
	}
	return
}

type folderFormatData struct {
	Data             interface{}
	FolderName       string
	TotalSize        int64
	NestLevel        uint
	ParentIdentifier string
	Identifier       string
}

type fileFormatData struct {
	Data             interface{}
	Filename         string
	Filesize         int64
	NestLevel        uint
	ParentIdentifier string
}

func execTemplateToHTML(tmpl *template.Template, data interface{}) (out template.HTML, err error) {
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return
	}
	out = template.HTML(buf.String())
	return
}

func (f *FileListFolder) makeFolderTreeView(folderTmpl *template.Template, fileTmpl *template.Template, nestLevel uint, identifier string, data interface{}) (output template.HTML, err error) {
	output = template.HTML("")
	var tmp template.HTML
	var folderNames []string // need this for sorting
	for folderName, _ := range f.Folders {
		folderNames = append(folderNames, folderName)
	}

	slice.Sort(folderNames, func(i, j int) bool {
		return strings.ToLower(folderNames[i]) < strings.ToLower(folderNames[j])
	})

	for i, folderName := range folderNames {
		folder := f.Folders[folderName]
		childIdentifier := identifier + "_d" + strconv.Itoa(i)
		// To the folder, our identifier is their parent identifier, and our child identifier is their own identifier.
		tmp, err = execTemplateToHTML(folderTmpl, folderFormatData{data, folderName, folder.TotalSize(), nestLevel, identifier, childIdentifier})
		if err != nil {
			return
		}
		output += tmp

		tmp, err = folder.makeFolderTreeView(folderTmpl, fileTmpl, nestLevel+1, childIdentifier, data)
		if err != nil {
			return
		}
		output += tmp
	}

	slice.Sort(f.Files, func(i, j int) bool {
		return strings.ToLower(f.Files[i].Filename()) < strings.ToLower(f.Files[j].Filename())
	})
	for _, file := range f.Files {
		tmp, err = execTemplateToHTML(fileTmpl, fileFormatData{data, file.Filename(), file.Filesize, nestLevel, identifier})
		if err != nil {
			return
		}
		output += tmp
	}

	return
}

func (f *FileListFolder) MakeFolderTreeView(folderFormat string, fileFormat string, funcMap template.FuncMap, data interface{}) (output template.HTML, err error) {
	folderTmpl, err := template.New("folderTemplate").Funcs(funcMap).Parse(folderFormat)
	if err != nil {
		return
	}
	fileTmpl, err := template.New("fileTemplate").Funcs(funcMap).Parse(fileFormat)
	if err != nil {
		return
	}

	output, err = f.makeFolderTreeView(folderTmpl, fileTmpl, 0, "root", data)
	return
}
