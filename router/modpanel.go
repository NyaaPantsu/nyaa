// hurry mod panel to get it faaaaaaaaaaaast

package router

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"fmt"

	"github.com/ewhal/nyaa/service/comment"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/service/torrent/form"
	"github.com/ewhal/nyaa/service/user"
	form "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/modelHelper"
)

var panelIndex, panelTorrentList, panelUserList, panelCommentList, panelTorrentEd *template.Template

func init() {
	panelTorrentList = template.Must(template.New("torrentlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/torrentlist.html")))
	panelUserList = template.Must(template.New("userlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/userlist.html")))
	panelCommentList = template.Must(template.New("commentlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/commentlist.html")))
	panelIndex = template.Must(template.New("indexPanel").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/panelindex.html")))
	panelTorrentEd = template.Must(template.New("indexPanel").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/paneltorrentedit.html")))
}

func IndexModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		offset := 10

		torrents, _, _ := torrentService.GetAllTorrents(0, offset)
		users := userService.RetrieveUsersForAdmin(0, offset)
		comments := commentService.GetAllComments(0, offset)
		languages.SetTranslationFromRequest(panelIndex, r, "en-us")
		htv := PanelIndexVbs{torrents, users, comments}
		_ = panelIndex.ExecuteTemplate(w, "admin_index.html", htv)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func TorrentsListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		page, _ := strconv.Atoi(r.URL.Query().Get("p"))
		offset := 100

		torrents, _, _ := torrentService.GetAllTorrents(offset, page * offset)
		languages.SetTranslationFromRequest(panelTorrentList, r, "en-us")
		htv := PanelTorrentListVbs{torrents}
		err := panelTorrentList.ExecuteTemplate(w, "admin_index.html", htv)
		fmt.Println(err)
	} else {

		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func UsersListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		page, _ := strconv.Atoi(r.URL.Query().Get("p"))
		offset := 100

		users := userService.RetrieveUsersForAdmin(offset, page*offset)
		languages.SetTranslationFromRequest(panelUserList, r, "en-us")
		htv := PanelUserListVbs{users}
		err := panelUserList.ExecuteTemplate(w, "admin_index.html", htv)
		fmt.Println(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func CommentsListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		page, _ := strconv.Atoi(r.URL.Query().Get("p"))
		offset := 100

		comments := commentService.GetAllComments(offset, page * offset)
		languages.SetTranslationFromRequest(panelCommentList, r, "en-us")
		htv := PanelCommentListVbs{comments}
		err := panelCommentList.ExecuteTemplate(w, "admin_index.html", htv)
		fmt.Println(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}

}
func TorrentEditModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		id := r.URL.Query().Get("id")
		torrent, _ := torrentService.GetTorrentById(id)
		languages.SetTranslationFromRequest(panelTorrentEd, r, "en-us")
		htv := PanelTorrentEdVbs{torrent}
		err := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
		fmt.Println(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}

}
func TorrentPostEditModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		b := torrentform.PanelPost{}
		err := form.NewErrors()
		infos := form.NewInfos()
		modelHelper.BindValueForm(&b, r)
		err = modelHelper.ValidateForm(&b, err)
		id := r.URL.Query().Get("id")
		torrent, _ := torrentService.GetTorrentById(id)
		if torrent.ID > 0 {
			modelHelper.AssignValue(&torrent, &b)
			if len(err) == 0 {
				_, errorT := torrentService.UpdateTorrent(torrent)
				if errorT != nil {
					err["errors"] = append(err["errors"], errorT.Error())
				}
				if len(err) == 0 {
					infos["infos"] = append(infos["infos"], "torrent_updated")
				}
			}
		}
		languages.SetTranslationFromRequest(panelTorrentEd, r, "en-us")
		htv := PanelTorrentEdVbs{torrent}
		_ = panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func CommentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	id := r.URL.Query().Get("id")

	if userPermission.HasAdmin(currentUser) {
		_ = form.NewErrors()
		_, _ = userService.DeleteComment(id)
		url, _ := Router.Get("mod_comment_list").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func TorrentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	id := r.URL.Query().Get("id")
	if userPermission.HasAdmin(currentUser) {
		_ = form.NewErrors()
		_, _ = torrentService.DeleteTorrent(id)
		url, _ := Router.Get("mod_torrent_list").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}
