package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/metainfo"
)

const (
	DatabaseDumpPath = "./public/dumps"
	GPGPublicKeyPath = "./public/gpg/gpg.key"
)

func DatabaseDumpHandler(w http.ResponseWriter, r *http.Request) {
	// db params url
	var err error
	// TODO Use config from cli
	files, err := filepath.Glob(filepath.Join(DatabaseDumpPath, "*.torrent"))
	var dumpsJson []model.DatabaseDumpJSON
	// TODO Filter *.torrent files
	for _, f := range files {
		// TODO Use config from cli
		file, err := os.Open(f)
		if err != nil {
			continue
		}
		var tf metainfo.TorrentFile
		err = tf.Decode(file)
		if err != nil {
			log.CheckError(err)
			fmt.Println(err)
			continue
		}
		dump := model.DatabaseDump{
			Date:        time.Now(),
			Filesize:    int64(tf.TotalSize()),
			Name:        tf.TorrentName(),
			TorrentLink: "/dbdumps/" + file.Name()}
		dumpsJson = append(dumpsJson, dump.ToJSON())
	}

	// TODO Remove ?
	navigationTorrents := Navigation{0, 0, 0, "search_page"}
	common := NewCommonVariables(r)
	common.Navigation = navigationTorrents
	dtv := DatabaseDumpTemplateVariables{common, dumpsJson, "/gpg/gpg.pub"}
	err = databaseDumpTemplate.ExecuteTemplate(w, "index.html", dtv)
	if err != nil {
		log.Errorf("DatabaseDump(): %s", err)
	}
}
