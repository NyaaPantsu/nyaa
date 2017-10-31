package torrentController

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/upload"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/gin-gonic/gin"
)

// DownloadTorrent : Controller for downloading a torrent
func DownloadTorrent(c *gin.Context) {
	hash := c.Param("hash")
	messages := msg.GetMessages(c)

	torrent, err := torrents.FindRawByHash(hash)

	if err != nil {
		messages.AddError("errors", "No torrent with such hash")
		//File not found, send 404
		variables := templates.Commonvariables(c)
		templates.Render(c, "errors/torrent_file_missing.jet.html", variables)
		return
	}

	if c.Query("js_query") != "" {
		exists := true
		generating := false

		if len(config.Get().Torrents.FileStorage) == 0 {
			exists = false
		} else {
			Openfile, err := os.Open(fmt.Sprintf("%s%c%s.torrent", config.Get().Torrents.FileStorage, os.PathSeparator, hash))
			defer Openfile.Close()
			if err != nil {
				exists = false
				generating = true

				var trackers []string
				if torrent.Trackers == "" {
					trackers = config.Get().Torrents.Trackers.Default
				} else {
					trackers = torrent.GetTrackersArray()
				}
				magnet := format.InfoHashToMagnet(strings.TrimSpace(torrent.Hash), torrent.Name, trackers...)
				if upload.GenerateTorrent(magnet) != nil {
					//Error during the generation
					generating = false
				}
			}
		}
		c.JSON(200, gin.H{ // Better to use gin for that, less code
			"exists":     exists,
			"generating": generating,
		})
		return
	}

	if len(config.Get().Torrents.FileStorage) == 0 { // if no FileStorage configured, you still can display the magnet link
		messages.AddError("errors", "We do not store torrents file")
		//File not found, send 404
		variables := templates.Commonvariables(c)
		var trackers []string
		if torrent.Trackers == "" {
			trackers = config.Get().Torrents.Trackers.Default
		} else {
			trackers = torrent.GetTrackersArray()
		}
		magnet := format.InfoHashToMagnet(strings.TrimSpace(torrent.Hash), torrent.Name, trackers...)
		variables.Set("magnet", magnet)
		templates.Render(c, "errors/torrent_file_missing.jet.html", variables)
		return
	}

	//Check if file exists and open
	Openfile, err := os.Open(torrent.GetPath())
	if err != nil {
		//File not found, send 404
		variables := templates.Commonvariables(c)
		var trackers []string
		if torrent.Trackers == "" {
			trackers = config.Get().Torrents.Trackers.Default
		} else {
			trackers = torrent.GetTrackersArray()
		}
		magnet := format.InfoHashToMagnet(strings.TrimSpace(torrent.Hash), torrent.Name, trackers...)
		variables.Set("magnet", magnet)
		if upload.GenerateTorrent(magnet) != nil {
			messages.AddError("errors", "Could not generate torrent file")
		}
		templates.Render(c, "errors/torrent_file_missing.jet.html", variables)
		return
	}
	defer Openfile.Close() //Close after function return

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.torrent\"", torrent.Name))
	c.Header("Content-Type", "application/x-bittorrent")
	c.Header("Content-Length", FileSize)
	//Send the file
	// We reset the offset to 0
	Openfile.Seek(0, 0)
	io.Copy(c.Writer, Openfile) //'Copy' the file to the client
}
