// hurry mod panel to get it faaaaaaaaaaaast

package router

import (
	"fmt"
	"net/http"
	"strconv"
	"html/template"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/service/user/permission"
	// "github.com/ewhal/nyaa/util/languages"
	// "github.com/ewhal/nyaa/util/modelHelper"
	"github.com/gorilla/mux"
)

var panelCommentList *template.Template
func init() {
	panelCommentList = template.Must(template.New("commentlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/commentlist.html")))
}

func IndexModPanel(w http.ResponseWriter, r *http.Request) {}
func TorrentsListPanel(w http.ResponseWriter, r *http.Request) {}
func UsersListPanel(w http.ResponseWriter, r *http.Request) {}
func CommentsListPanel(w http.ResponseWriter, r *http.Request) {
	page,_ := strconv.Atoi(r.URL.Query().Get("p"))
	offset := 100

	comments := commentService.GetAllTorrents(page*offset, offset)
	languages.SetTranslationFromRequest(panelCommentList, r, "en-us")
	htv := PanelCommentListVbs{comments}
	err := panelCommentList.ExecuteTemplate(w, "index.html", htv)

}
func TorrentEditModPanel(w http.ResponseWriter, r *http.Request) {
	// Todo 
}
func CommentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if (HasAdmin(currentUser)) {
		err := form.NewErrors()
		_, _ := userService.DeleteComment(id)
		url, _ := Router.Get("mod_comment_list")
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	}
}
func TorrentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if (HasAdmin(currentUser)) {
		err := form.NewErrors()
		_, _ := torrentService.DeleteTorrent(id)
		url, _ := Router.Get("mod_torrent_list")
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	}
}