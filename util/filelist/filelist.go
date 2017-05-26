package filelist

import (
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/bradfitz/slice"
	"strings"
)

type FileListFolder struct {
	Folders    []*FileListFolder
	Files      []model.File
	FolderName string
}

func FileListToFolder(fileList []model.File, folderName string) (out *FileListFolder) {
	out = &FileListFolder{
		Folders:    make([]*FileListFolder, 0),
		Files:      make([]model.File, 0),
		FolderName: folderName,
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
		out.Folders = append(out.Folders, FileListToFolder(folderFiles, folderName))
	}

	// Do some sorting
	slice.Sort(out.Folders, func(i, j int) bool {
		return strings.ToLower(out.Folders[i].FolderName) < strings.ToLower(out.Folders[j].FolderName)
	})
	slice.Sort(out.Files, func(i, j int) bool {
		return strings.ToLower(out.Files[i].Filename()) < strings.ToLower(out.Files[i].Filename())
	})
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

