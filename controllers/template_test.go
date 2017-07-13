package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"path"
	"testing"

	"strings"

	"time"

	"github.com/CloudyKit/jet"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
)

// run before router/init.go:init()
var _ = func() (_ struct{}) {
	gin.SetMode(gin.TestMode)
	config.ConfigPath = path.Join("..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", config.DefaultConfigPath)
	config.Reload()
	return
}()

func TestTemplates(t *testing.T) {
	//var View = jet.NewHTMLSet(TemplateDir)
	fmt.Print("JetTest Template started\n")

	walkDirTest("/", t)

}

type ContextTest map[string]func(jet.VarMap) jet.VarMap

func walkDirTest(dir string, t *testing.T) {
	fakeUser := &models.User{1, "test", "test", "test", 1, time.Now(), time.Now(), "test", time.Now(), "en", "test", "test", "test", "test", []models.User{}, []models.User{}, "test", []models.Torrent{}, []models.Notification{}, 1, models.UserSettings{}}
	fakeComment := &models.Comment{1, 1, 1, "test", time.Now(), time.Now(), nil, &models.Torrent{}, fakeUser}
	fakeScrapeData := &models.Scrape{1, 0, 0, 10, time.Now()}
	fakeFile := &models.File{1, 1, "l12:somefile.mp4e", 3}
	fakeLanguages := []string{"fr", "en"}
	fakeTorrent := &models.Torrent{1, "test", "test", 3, 12, 1, false, time.Now(), 1, 0, 3, "test", "test", "test", "test", "test", nil, fakeUser, "test", []models.OldComment{}, []models.Comment{*fakeComment, *fakeComment}, fakeScrapeData, []models.File{*fakeFile}, fakeLanguages}
	fakeActivity := &models.Activity{1, "t", "e", "s", 1, fakeUser}
	fakeDB := &models.DatabaseDump{time.Now(), 3, "test", "test"}
	fakeLanguage := &publicSettings.Language{"English", "en", "en-us"}
	fakeTorrentRequest := &torrentValidator.TorrentRequest{Name: "test", Magnet: "", Category: "", Remake: false, Description: "", Status: 1, Hidden: false, CaptchaID: "", WebsiteLink: "", SubCategory: 0, Languages: nil, Infohash: "", SubCategoryID: 0, CategoryID: 0, Filesize: 0, Filepath: "", FileList: nil, Trackers: nil}
	fakeLogin := &userValidator.LoginForm{"test", "test", "/"}
	fakeRegistration := &userValidator.RegistrationForm{"test", "", "test", "test", "xxxx", "1"}
	fakeReport := &models.TorrentReport{1, "test", 1, 1, time.Now(), fakeTorrent, fakeUser}
	contextVars := ContextTest{
		"dumps.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("GPGLink", "test")
			vars.Set("ListDumps", []models.DatabaseDumpJSON{fakeDB.ToJSON(), fakeDB.ToJSON()})
			return vars
		},
		"activities.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.Activity{*fakeActivity})
			return vars
		},
		"listing.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.TorrentJSON{fakeTorrent.ToJSON(), fakeTorrent.ToJSON()})
			return vars
		},
		"edit.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", fakeTorrentRequest)
			vars.Set("Languages", publicSettings.Languages{*fakeLanguage, *fakeLanguage})
			return vars
		},
		"upload.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", fakeTorrentRequest)
			return vars
		},
		"view.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Torrent", fakeTorrent.ToJSON())
			vars.Set("CaptchaID", "xxxxxx")
			vars.Set("RootFolder", filelist.FileListToFolder(fakeTorrent.FileList, "root"))
			return vars
		},
		"settings.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", &LanguagesJSONResponse{"test", publicSettings.Languages{*fakeLanguage, *fakeLanguage}})
			return vars
		},
		"login.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", fakeLogin)
			return vars
		},
		"register.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", fakeRegistration)
			return vars
		},
		"index.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Torrents", []models.Torrent{*fakeTorrent, *fakeTorrent})
			vars.Set("Users", []models.User{*fakeUser, *fakeUser})
			vars.Set("Comments", []models.Comment{*fakeComment, *fakeComment})
			vars.Set("TorrentReports", []models.TorrentReportJSON{fakeReport.ToJSON(), fakeReport.ToJSON()})
			return vars
		},
		"paneltorrentedit.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", *fakeTorrent)
			return vars
		},
		"reassign.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", ReassignForm{1, "", "", []uint{1, 1}})
			return vars
		},
		"torrentlist.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.Torrent{*fakeTorrent, *fakeTorrent})
			return vars
		},
		"userlist.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.User{*fakeUser, *fakeUser})
			return vars
		},
		"commentlist.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.Comment{*fakeComment, *fakeComment})
			return vars
		},
		"torrent_report.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.TorrentReportJSON{fakeReport.ToJSON(), fakeReport.ToJSON()})
			return vars
		},
		"report.jet.html": func(vars jet.VarMap) jet.VarMap {
			type form struct {
				ID        int
				CaptchaID string
			}
			vars.Set("Form", form{1, "test"})
			return vars
		},
	}

	fmt.Printf("\nTesting Folder: %s\n", dir)
	view := jet.NewHTMLSet(path.Join("..", TemplateDir))
	vars := mockupCommonVars(t)
	files, err := ioutil.ReadDir(path.Join("..", TemplateDir) + dir)
	if err != nil {
		t.Errorf("Couldn't find the folder %s", path.Join("..", TemplateDir)+dir)
	}
	if len(files) == 0 {
		t.Errorf("Couldn't find any files in folder %s", path.Join("..", TemplateDir)+dir)
	}
	for _, f := range files {
		if f.IsDir() {
			walkDirTest(dir+f.Name()+"/", t)
			continue
		}
		if strings.Contains(f.Name(), ".jet.html") {
			template, err := view.GetTemplate(dir + f.Name())
			fmt.Printf("\tJetTest Template of: %s", dir+f.Name())
			if err != nil {
				t.Errorf("\nParsing error: %s %s", err.Error(), dir+f.Name())
			}
			buff := bytes.NewBuffer(nil)
			if contextVars[f.Name()] != nil {
				vars = contextVars[f.Name()](vars)
			}
			if err = template.Execute(buff, vars, nil); err != nil {
				t.Errorf("\nEval error: %q executing %s", err.Error(), template.Name)
			}
			fmt.Print("\tOK\n")
		}
	}
}

func mockupCommonVars(t *testing.T) jet.VarMap {
	vars.Set("Navigation", newNavigation())
	vars.Set("Search", searchForm{
		Category:         "_",
		ShowItemsPerPage: true,
		SizeType:         "b",
		DateType:         "test",
		MinSize:          "test",
		MaxSize:          "test",
		FromDate:         "test",
		ToDate:           "test",
	})
	conf := config.Get().I18n
	conf.Directory = path.Join("..", conf.Directory)
	var retriever publicSettings.UserRetriever // not required during initialization

	err := publicSettings.InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
	Ts, _, err := publicSettings.TfuncAndLanguageWithFallback("en-us", "test", "test")
	if err != nil {
		t.Error("Couldn't load language files!")
	}
	var T publicSettings.TemplateTfunc
	T = func(id string, args ...interface{}) template.HTML {
		return template.HTML(fmt.Sprintf(Ts(id), args...))
	}
	vars.Set("T", T)
	vars.Set("Theme", "test")
	vars.Set("Mascot", "test")
	vars.Set("MascotURL", "test")
	vars.Set("User", &models.User{})
	vars.Set("UserProfile", &models.User{})
	vars.Set("URL", &url.URL{})
	vars.Set("CsrfToken", "xxxxxx")
	vars.Set("Config", config.Get())
	vars.Set("Infos", make(map[string][]string))
	vars.Set("Errors", make(map[string][]string))
	vars.Set("UserProfile", &models.User{})
	return vars
}
