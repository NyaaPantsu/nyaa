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
	// DatabaseDumpPath : Location of database dumps
	DatabaseDumpPath = "./public/dumps"
	// GPGPublicKeyPath : Location of the GPG key
	GPGPublicKeyPath = "./public/gpg/gpg.key"
)

// DatabaseDumpHandler : Controller for getting database dumps
func DatabaseDumpHandler(w http.ResponseWriter, r *http.Request) {
	// db params url
	var err error
	// TODO Use config from cli
	files, _ := filepath.Glob(filepath.Join(DatabaseDumpPath, "*.torrent"))
	defer r.Body.Close()
	var dumpsJSON []model.DatabaseDumpJSON
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
		dumpsJSON = append(dumpsJSON, dump.ToJSON())
	}

	// TODO Remove ?
	navigationTorrents := navigation{0, 0, 0, "search_page"}
	common := newCommonVariables(r)
	common.Navigation = navigationTorrents
	dtv := databaseDumpTemplateVariables{common, dumpsJSON, "/gpg/gpg.pub"}
	err = databaseDumpTemplate.ExecuteTemplate(w, "index.html", dtv)
	if err != nil {
		log.Errorf("DatabaseDump(): %s", err)
	}
}
