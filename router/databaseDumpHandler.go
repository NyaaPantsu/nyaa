package router

import (
    "io/ioutil"
	"os"
	"net/http"
	"fmt"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/metainfo"
	"github.com/gorilla/mux"
)

func DatabaseDumpHandler(w http.ResponseWriter, r *http.Request) {
	// db params url
	var err error
	// TODO Use config from cli
    files, _ := ioutil.ReadDir(config.DefaultDatabaseDumpPath)
	if len(files) <= 0 {
		return
	}
	var dumpsJson []model.DatabaseDumpJSON
	// TODO Filter *.torrent files
    for _, f := range files {
		// TODO Use config from cli
		file, err := os.Open(config.DefaultDatabaseDumpPath + f.Name())
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
			TorrentLink: "/dbdumps/" + f.Name()}
		dumpsJson = append(dumpsJson, dump.ToJSON())
    }

	// TODO Remove ?
	navigationTorrents := Navigation{0, 0, 0, "search_page"}
	languages.SetTranslationFromRequest(databaseDumpTemplate, r, "en-us")
	dtv := DatabaseDumpTemplateVariables{dumpsJson, "/gpg/gpg.pub", NewSearchForm(), navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}
	err = databaseDumpTemplate.ExecuteTemplate(w, "index.html", dtv)
	if err != nil {
		log.Errorf("DatabaseDump(): %s", err)
	}
}

