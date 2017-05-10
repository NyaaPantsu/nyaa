// hurry mod panel to get it faaaaaaaaaaaast

package router

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"fmt"

	"github.com/ewhal/nyaa/service/comment"
	"github.com/ewhal/nyaa/service/moderation"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/service/torrent/form"
	"github.com/ewhal/nyaa/util/search"
	"github.com/ewhal/nyaa/service/user"
	form "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/gorilla/mux"
)

var panelIndex, panelTorrentList, panelUserList, panelCommentList, panelTorrentEd, torrentReportTemplate *template.Template

func init() {
	panelTorrentList = template.Must(template.New("torrentlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/torrentlist.html")))
	panelTorrentList = template.Must(panelTorrentList.ParseGlob(filepath.Join("templates", "_*.html")))
	panelUserList = template.Must(template.New("userlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/userlist.html")))
		panelUserList = template.Must(panelUserList.ParseGlob(filepath.Join("templates", "_*.html")))
	panelCommentList = template.Must(template.New("commentlist").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/commentlist.html")))
		panelCommentList = template.Must(panelCommentList.ParseGlob(filepath.Join("templates", "_*.html")))
	panelIndex = template.Must(template.New("indexPanel").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/panelindex.html")))
		panelIndex = template.Must(panelIndex.ParseGlob(filepath.Join("templates", "_*.html")))
	panelTorrentEd = template.Must(template.New("torrent_ed").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/paneltorrentedit.html")))
	panelTorrentEd = 	template.Must(panelTorrentEd.ParseGlob(filepath.Join("templates", "_*.html")))
	torrentReportTemplate = template.Must(template.New("torrent_report").Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "admin_index.html"), filepath.Join(TemplateDir, "admin/torrent_report.html")))
	torrentReportTemplate = template.Must(torrentReportTemplate.ParseGlob(filepath.Join("templates", "_*.html")))
}

func IndexModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		offset := 10

		torrents, _, _ := torrentService.GetAllTorrents(offset, 0)
		users, _ := userService.RetrieveUsersForAdmin(offset, 0)
		comments, _ := commentService.GetAllComments(offset, 0, "", "")
		torrentReports, _, _ := moderationService.GetTorrentReports(offset,0, "", "")
		languages.SetTranslationFromRequest(panelIndex, r, "en-us")
		htv := PanelIndexVbs{torrents, torrentReports,  users, comments, NewSearchForm(), currentUser, r.URL}
		_ = panelIndex.ExecuteTemplate(w, "admin_index.html", htv)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func TorrentsListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		vars := mux.Vars(r)
		page, _ := strconv.Atoi(vars["page"])
		offset := 100

		searchParam, torrents, _, err := search.SearchByQuery(r, page)
		searchForm := SearchForm{
			SearchParam:        searchParam,
			Category:           searchParam.Category.String(),
			HideAdvancedSearch: false,
		}

		languages.SetTranslationFromRequest(panelTorrentList, r, "en-us")
		htv := PanelTorrentListVbs{torrents, searchForm, Navigation{int(searchParam.Max), offset, page, "mod_tlist_page"}, currentUser, r.URL}
		err = panelTorrentList.ExecuteTemplate(w, "admin_index.html", htv)
		fmt.Println(err)
	} else {

		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func UsersListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		vars := mux.Vars(r)
		page, _ := strconv.Atoi(vars["page"])
		offset := 100

		users, nbUsers := userService.RetrieveUsersForAdmin(offset, page*offset)
		languages.SetTranslationFromRequest(panelUserList, r, "en-us")
		htv := PanelUserListVbs{users, NewSearchForm(), Navigation{nbUsers, offset, page, "mod_ulist_page"}, currentUser, r.URL}
		err := panelUserList.ExecuteTemplate(w, "admin_index.html", htv)
		fmt.Println(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}
func CommentsListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		vars := mux.Vars(r)
		page, _ := strconv.Atoi(vars["page"])
		offset := 100
		userid := r.URL.Query().Get("userid")
		var conditions string
		var values []interface{}
		if (userid != "") {
			conditions = "user_id = ?"
			values = append(values, userid)
		}
		comments, nbComments := commentService.GetAllComments(offset, page * offset, conditions, values...)
		languages.SetTranslationFromRequest(panelCommentList, r, "en-us")
		htv := PanelCommentListVbs{comments, NewSearchForm(), Navigation{nbComments, offset, page, "mod_clist_page"}, currentUser, r.URL}
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
		htv := PanelTorrentEdVbs{torrent, NewSearchForm(), currentUser}
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
		htv := PanelTorrentEdVbs{torrent, NewSearchForm(), currentUser}
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
		url, _ := Router.Get("mod_clist").URL()
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
		url, _ := Router.Get("mod_tlist").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}