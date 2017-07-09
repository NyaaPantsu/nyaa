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

	"github.com/CloudyKit/jet"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
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
	config.Parse()
	return
}()

func TestTemplates(t *testing.T) {
	//var View = jet.NewHTMLSet(TemplateDir)
	fmt.Print("JetTest Template started\n")

	walkDirTest("/", t)

}

type ContextTest map[string]func(jet.VarMap) jet.VarMap

func walkDirTest(dir string, t *testing.T) {
	contextVars := ContextTest{
		"dumps.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("GPGLink", "")
			vars.Set("ListDumps", []models.DatabaseDumpJSON{})
			return vars
		},
		"activities.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.Activity{})
			return vars
		},
		"listing.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.TorrentJSON{})
			return vars
		},
		"edit.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", &torrentValidator.TorrentRequest{})
			vars.Set("Languages", publicSettings.Languages{{"", ""}})
			return vars
		},
		"upload.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", &torrentValidator.TorrentRequest{})
			return vars
		},
		"view.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Torrent", &models.TorrentJSON{})
			vars.Set("CaptchaID", "xxxxxx")
			return vars
		},
		"settings.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", &LanguagesJSONResponse{"", publicSettings.Languages{{"", ""}}})
			return vars
		},
		"login.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", &userValidator.LoginForm{})
			return vars
		},
		"register.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", &userValidator.RegistrationForm{})
			return vars
		},
		"index.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Torrents", []models.Torrent{})
			vars.Set("Users", []models.User{})
			vars.Set("Comments", []models.Comment{})
			vars.Set("TorrentReports", []models.TorrentReportJSON{})
			return vars
		},
		"paneltorrentedit.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", models.Torrent{})
			return vars
		},
		"reassign.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Form", ReassignForm{})
			return vars
		},
		"torrentlist.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.Torrent{})
			return vars
		},
		"userlist.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.User{})
			return vars
		},
		"commentlist.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.Comment{})
			return vars
		},
		"torrent_report.jet.html": func(vars jet.VarMap) jet.VarMap {
			vars.Set("Models", []models.TorrentReportJSON{})
			return vars
		},
		"report.jet.html": func(vars jet.VarMap) jet.VarMap {
			type form struct {
				ID        int
				CaptchaID string
			}
			vars.Set("Form", form{1, ""})
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
		DateType:         "",
		MinSize:          "",
		MaxSize:          "",
		FromDate:         "",
		ToDate:           "",
	})
	conf := config.Conf.I18n
	conf.Directory = path.Join("..", conf.Directory)
	var retriever publicSettings.UserRetriever // not required during initialization

	err := publicSettings.InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
	Ts, _, err := publicSettings.TfuncAndLanguageWithFallback("en-us", "", "")
	if err != nil {
		t.Error("Couldn't load language files!")
	}
	var T publicSettings.TemplateTfunc
	T = func(id string, args ...interface{}) template.HTML {
		return template.HTML(fmt.Sprintf(Ts(id), args...))
	}
	vars.Set("T", T)
	vars.Set("Theme", "")
	vars.Set("Mascot", "")
	vars.Set("MascotURL", "")
	vars.Set("User", &models.User{})
	vars.Set("UserProfile", &models.User{})
	vars.Set("URL", &url.URL{})
	vars.Set("CsrfToken", "xxxxxx")
	vars.Set("Config", config.Conf)
	vars.Set("Infos", make(map[string][]string))
	vars.Set("Errors", make(map[string][]string))
	vars.Set("UserProfile", &models.User{})
	return vars
}
