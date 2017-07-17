package templates

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
	contextvariables := ContextTest{
		"dumps.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("GPGLink", "test")
			variables.Set("ListDumps", []models.DatabaseDumpJSON{fakeDB.ToJSON(), fakeDB.ToJSON()})
			return variables
		},
		"activities.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Models", []models.Activity{*fakeActivity})
			return variables
		},
		"listing.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Models", []models.TorrentJSON{fakeTorrent.ToJSON(), fakeTorrent.ToJSON()})
			return variables
		},
		"edit.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", fakeTorrentRequest)
			variables.Set("Languages", publicSettings.Languages{*fakeLanguage, *fakeLanguage})
			return variables
		},
		"upload.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", fakeTorrentRequest)
			return variables
		},
		"view.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Torrent", fakeTorrent.ToJSON())
			variables.Set("CaptchaID", "xxxxxx")
			variables.Set("RootFolder", filelist.FileListToFolder(fakeTorrent.FileList, "root"))
			return variables
		},
		"settings.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", &LanguagesJSONResponse{"test", publicSettings.Languages{*fakeLanguage, *fakeLanguage}})
			return variables
		},
		"login.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", fakeLogin)
			return variables
		},
		"register.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", fakeRegistration)
			return variables
		},
		"index.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Torrents", []models.Torrent{*fakeTorrent, *fakeTorrent})
			variables.Set("Users", []models.User{*fakeUser, *fakeUser})
			variables.Set("Comments", []models.Comment{*fakeComment, *fakeComment})
			variables.Set("TorrentReports", []models.TorrentReportJSON{fakeReport.ToJSON(), fakeReport.ToJSON()})
			return variables
		},
		"paneltorrentedit.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", *fakeTorrent)
			return variables
		},
		"reassign.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Form", torrentValidator.ReassignForm{1, "", "", []uint{1, 1}})
			return variables
		},
		"torrentlist.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Models", []models.Torrent{*fakeTorrent, *fakeTorrent})
			return variables
		},
		"userlist.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Models", []models.User{*fakeUser, *fakeUser})
			return variables
		},
		"commentlist.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Models", []models.Comment{*fakeComment, *fakeComment})
			return variables
		},
		"torrent_report.jet.html": func(variables jet.VarMap) jet.VarMap {
			variables.Set("Models", []models.TorrentReportJSON{fakeReport.ToJSON(), fakeReport.ToJSON()})
			return variables
		},
		"report.jet.html": func(variables jet.VarMap) jet.VarMap {
			type form struct {
				ID        int
				CaptchaID string
			}
			variables.Set("Form", form{1, "test"})
			return variables
		},
	}

	fmt.Printf("\nTesting Folder: %s\n", dir)
	view := jet.NewHTMLSet(path.Join("..", TemplateDir))
	files, err := ioutil.ReadDir(path.Join("..", TemplateDir) + dir)
	if err != nil {
		t.Errorf("Couldn't find the folder %s", path.Join("..", TemplateDir)+dir)
	}
	if len(files) == 0 {
		t.Errorf("Couldn't find any files in folder %s", path.Join("..", TemplateDir)+dir)
	}
	for _, f := range files {
		variables := mockupCommonvariables(t)
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
			if contextvariables[f.Name()] != nil {
				variables = contextvariables[f.Name()](variables)
			}
			if err = template.Execute(buff, variables, nil); err != nil {
				t.Errorf("\nEval error: %q executing %s", err.Error(), template.Name)
			}
			fmt.Print("\tOK\n")
		}
	}
}

func mockupCommonvariables(t *testing.T) jet.VarMap {
	variables := jet.VarMap{}
	variables.Set("Navigation", NewNavigation())
	variables.Set("Search", SearchForm{
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
	variables.Set("T", T)
	variables.Set("Theme", "test")
	variables.Set("Mascot", "test")
	variables.Set("MascotURL", "test")
	variables.Set("User", &models.User{})
	variables.Set("UserProfile", &models.User{})
	variables.Set("URL", &url.URL{})
	variables.Set("CsrfToken", "xxxxxx")
	variables.Set("Config", config.Get())
	variables.Set("Infos", make(map[string][]string))
	variables.Set("Errors", make(map[string][]string))
	variables.Set("UserProfile", &models.User{})
	variables = templateFunctions(variables)
	return variables
}
