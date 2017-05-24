package router

import (
	"html/template"
	"path/filepath"
)

// FIXME This should be a constant, but then the tests fails
var TemplateDir = "templates"

const ModeratorDir = "admin"

var homeTemplate,
	searchTemplate,
	faqTemplate,
	uploadTemplate,
	viewTemplate,
	viewRegisterTemplate,
	viewLoginTemplate,
	viewRegisterSuccessTemplate,
	viewVerifySuccessTemplate,
	viewProfileTemplate,
	viewProfileNotifTemplate,
	viewProfileEditTemplate,
	viewUserDeleteTemplate,
	userTorrentEd,
	notFoundTemplate,
	changeLanguageTemplate,
	databaseDumpTemplate *template.Template

var panelIndex,
	panelTorrentList,
	panelUserList,
	panelCommentList,
	panelTorrentEd,
	panelTorrentReportList,
	panelTorrentReassign *template.Template

type templateLoader struct {
	templ     **template.Template
	file      string
	indexFile string
	name      string
}

// ReloadTemplates reloads templates on runtime
func ReloadTemplates() {
	pubTempls := []templateLoader{
		{
			templ: &databaseDumpTemplate,
			name:  "dump",
			file:  "dumps.html",
		},
		{
			templ: &homeTemplate,
			name:  "home",
			file:  "home.html",
		},
		{
			templ: &searchTemplate,
			name:  "search",
			file:  "home.html",
		},
		{
			templ: &uploadTemplate,
			name:  "upload",
			file:  "upload.html",
		},
		{
			templ: &faqTemplate,
			name:  "FAQ",
			file:  "FAQ.html",
		},
		{
			templ: &viewTemplate,
			name:  "view",
			file:  "view.html",
		},
		{
			templ: &viewRegisterTemplate,
			name:  "user_register",
			file:  filepath.Join("user", "register.html"),
		},
		{
			templ: &viewRegisterSuccessTemplate,
			name:  "user_register_success",
			file:  filepath.Join("user", "signup_success.html"),
		},
		{
			templ: &viewVerifySuccessTemplate,
			name:  "user_verify_success",
			file:  filepath.Join("user", "verify_success.html"),
		},
		{
			templ: &viewLoginTemplate,
			name:  "user_login",
			file:  filepath.Join("user", "login.html"),
		},
		{
			templ: &viewProfileTemplate,
			name:  "user_profile",
			file:  filepath.Join("user", "profile.html"),
		},
		{
			templ: &viewProfileNotifTemplate,
			name:  "user_profile",
			file:  filepath.Join("user", "profile_notifications.html"),
		},
		{
			templ: &viewProfileEditTemplate,
			name:  "user_profile",
			file:  filepath.Join("user", "profile_edit.html"),
		},
		{
			templ: &viewUserDeleteTemplate,
			name:  "user_delete",
			file:  filepath.Join("user", "delete_success.html"),
		},
		{
			templ: &userTorrentEd,
			name:  "user_torrent_edit",
			file:  filepath.Join("user", "torrent_edit.html"),
		},
		{
			templ: &notFoundTemplate,
			name:  "404",
			file:  "404.html",
		},
		{
			templ: &changeLanguageTemplate,
			name:  "change_language",
			file:  "change_language.html",
		},
	}
	for idx := range pubTempls {
		pubTempls[idx].indexFile = filepath.Join(TemplateDir, "index.html")
	}

	modTempls := []templateLoader{
		{
			templ: &panelTorrentList,
			name:  "torrentlist",
			file:  filepath.Join(ModeratorDir, "torrentlist.html"),
		},
		{
			templ: &panelUserList,
			name:  "userlist",
			file:  filepath.Join(ModeratorDir, "userlist.html"),
		},
		{
			templ: &panelCommentList,
			name:  "commentlist",
			file:  filepath.Join(ModeratorDir, "commentlist.html"),
		},
		{
			templ: &panelIndex,
			name:  "indexPanel",
			file:  filepath.Join(ModeratorDir, "panelindex.html"),
		},
		{
			templ: &panelTorrentEd,
			name:  "torrent_ed",
			file:  filepath.Join(ModeratorDir, "paneltorrentedit.html"),
		},
		{
			templ: &panelTorrentReportList,
			name:  "torrent_report",
			file:  filepath.Join(ModeratorDir, "torrent_report.html"),
		},
		{
			templ: &panelTorrentReassign,
			name:  "torrent_reassign",
			file:  filepath.Join(ModeratorDir, "reassign.html"),
		},
	}

	for idx := range modTempls {
		modTempls[idx].indexFile = filepath.Join(TemplateDir, "admin_index.html")
	}

	templs := make([]templateLoader, 0, len(modTempls)+len(pubTempls))
	templs = append(templs, pubTempls...)
	templs = append(templs, modTempls...)

	for _, templ := range templs {
		t := template.Must(template.New(templ.name).Funcs(FuncMap).ParseFiles(templ.indexFile, filepath.Join(TemplateDir, templ.file)))
		t = template.Must(t.ParseGlob(filepath.Join(TemplateDir, "_*.html")))
		*templ.templ = t
	}
}
