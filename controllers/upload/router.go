package uploadController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/upload", UploadHandler)
	router.Get().Any("/upload/status/:id", multiUploadStatus)
}
