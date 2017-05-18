package router

import (
	"os"
	"path/filepath"
	"net/http"
	"fmt"
	"time"

	"github.com/NyaanPantsu/nyaa/model"
	"github.com/NyaanPantsu/nyaa/util/languages"
	"github.com/NyaanPantsu/nyaa/util/log"
	"github.com/NyaanPantsu/nyaa/util/metainfo"
	"github.com/gorilla/mux"
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
	languages.SetTranslationFromRequest(databaseDumpTemplate, r)
	dtv := DatabaseDumpTemplateVariables{dumpsJson, "/gpg/gpg.pub", NewSearchForm(), navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}
	err = databaseDumpTemplate.ExecuteTemplate(w, "index.html", dtv)
	if err != nil {
		log.Errorf("DatabaseDump(): %s", err)
	}
}

