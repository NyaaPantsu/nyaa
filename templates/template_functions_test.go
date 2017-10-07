package templates

import (
	"fmt"
	"html/template"
	"net/url"
	"path"
	"testing"

	"time"

	"reflect"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

// run before router/init.go:init()
var _ = func() (_ struct{}) {
	categories.InitCategories()
	return
}()

func TestGetRawQuery(t *testing.T) {
	var tests = []map[string]string{
		{
			"test":     "",
			"expected": "",
		},
		{
			"test":     "http://lol.co/",
			"expected": "",
		},
		{
			"test":     "lol.co",
			"expected": "",
		},
		{
			"test":     "lol.co?",
			"expected": "",
		},
		{
			"test":     "lol.co?why",
			"expected": "why",
		},
		{
			"test":     "https://lol.co?why",
			"expected": "why",
		},
	}

	for _, test := range tests {
		url, _ := url.Parse(test["test"])
		value := getRawQuery(url)
		if value != test["expected"] {
			t.Errorf("Unexpected value from the function getRawQuery, got '%s', wanted '%s' for '%s'", value, test["expected"], test["test"])
		}
	}
}

func TestGenSearchWithOrdering(t *testing.T) {
	var tests = []map[string]string{
		{
			"test":     "",
			"mode":     "2",
			"expected": "/search?order=true&sort=2",
		},
		{
			"test":     "http://lol.co/?s=why&sort=1",
			"mode":     "2",
			"expected": "/search?order=false&s=why&sort=2",
		},
		{
			"test":     "http://lol.co/?s=why&sort=1",
			"mode":     "1",
			"expected": "/search?order=true&s=why&sort=1",
		},
		{
			"test":     "http://lol.co/?s=why&sort=1&order=true",
			"mode":     "1",
			"expected": "/search?order=false&s=why&sort=1",
		},
		{
			"test":     "http://lol.co/?s=why&sort=1&order=false",
			"mode":     "1",
			"expected": "/search?order=true&s=why&sort=1",
		},
		{
			"test":     "http://lol.co/?s=why&sort=1&order=false",
			"mode":     "2",
			"expected": "/search?order=false&s=why&sort=2",
		},
		{
			"test":     "http://lol.co/?s=why&sort=1&order=true",
			"mode":     "2",
			"expected": "/search?order=false&s=why&sort=2",
		},
	}

	for _, test := range tests {
		url, _ := url.Parse(test["test"])
		value := genSearchWithOrdering(url, test["mode"])
		if value != test["expected"] {
			t.Errorf("Unexpected value from the function genSearchWithOrdering, got '%s', wanted '%s' for '%s' and '%s'", value, test["expected"], test["test"], test["mode"])
		}
	}
}

func TestgenSearchWithCategory(t *testing.T) {
	var tests = []map[string]string{
		{
			"test":     "",
			"mode":     "1_",
			"expected": "/search?c=1_",
		},
	}

	for _, test := range tests {
		url, _ := url.Parse(test["test"])
		value := genSearchWithCategory(url, test["mode"])
		if value != test["expected"] {
			t.Errorf("Unexpected value from the function genSearchWithCategory, got '%s', wanted '%s' for '%s' and '%s'", value, test["expected"], test["test"], test["mode"])
		}
	}
}

func TestFlagCode(t *testing.T) {
	var tests = []map[string]string{
		{
			"test":     "",
			"expected": "und",
		},
		{
			"test":     "es",
			"expected": "es",
		},
		{
			"test":     "lol",
			"expected": "lol",
		},
		{
			"test":     "fr-fr",
			"expected": "fr",
		},
		{
			"test":     "fr-lol",
			"expected": "lol",
		},
		{
			"test":     "ca-es",
			"expected": "ca",
		},
		{
			"test":     "es-mx",
			"expected": "es",
		},
	}

	for _, test := range tests {
		value := flagCode(test["test"])
		if value != test["expected"] {
			t.Errorf("Unexpected value from the function flagCode, got '%s', wanted '%s' for '%s'", value, test["expected"], test["test"])
		}
	}
}

func TestGetAvatar(t *testing.T) {
	var tests = []struct {
		Test     string
		Size     int
		Expected string
	}{
		{
			Test:     "",
			Size:     0,
			Expected: "https://www.gravatar.com/avatar/?s=0",
		},
		{
			Test:     "",
			Size:     100,
			Expected: "https://www.gravatar.com/avatar/?s=100",
		},
		{
			Test:     "test",
			Size:     100,
			Expected: "https://www.gravatar.com/avatar/test?s=100",
		},
		{
			Test:     "test",
			Size:     0,
			Expected: "https://www.gravatar.com/avatar/test?s=0",
		},
	}

	for _, test := range tests {
		value := getAvatar(test.Test, test.Size)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function getAvatar, got '%s', wanted '%s' for '%s' and '%d'", value, test.Expected, test.Test, test.Size)
		}
	}
}

func TestFormatDateRFC(t *testing.T) {
	location, _ := time.LoadLocation("UTC")
	var tests = []struct {
		Test     time.Time
		Expected string
	}{
		{
			Test:     time.Date(2016, 5, 4, 3, 2, 1, 10, location),
			Expected: "2016-05-04T03:02:01Z",
		},
		{
			Test:     time.Now(),
			Expected: time.Now().Format(time.RFC3339),
		},
	}

	for _, test := range tests {
		value := formatDateRFC(test.Test)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function formatDateRFC, got '%s', wanted '%s' for '%s'", value, test.Expected, test.Test.String())
		}
	}
}

func TestGetCategory(t *testing.T) {
	var tests = []struct {
		TestCat    string
		TestParent bool
		Expected   categories.Categories
	}{
		{
			TestCat:    "",
			TestParent: false,
			Expected:   categories.Categories{},
		},
		{
			TestCat:    "",
			TestParent: true,
			Expected:   categories.Categories{},
		},
		{
			TestCat:    "3_12",
			TestParent: false,
			Expected:   categories.Categories{},
		},
		{
			TestCat:    "3",
			TestParent: false,
			Expected: categories.Categories{
				{"3_12", "anime_amv"},
				{"3_5", "anime_english_translated"},
				{"3_13", "anime_non_english_translated"},
				{"3_6", "anime_raw"},
			},
		},
		{
			TestCat:    "3",
			TestParent: true,
			Expected: categories.Categories{
				{"3_", "anime"},
				{"3_12", "anime_amv"},
				{"3_5", "anime_english_translated"},
				{"3_13", "anime_non_english_translated"},
				{"3_6", "anime_raw"},
			},
		},
	}
	for _, test := range tests {
		value := getCategory(test.TestCat, test.TestParent)
		if !reflect.DeepEqual(value, test.Expected) {
			t.Errorf("Unexpected value from the function getCategory, got '%v', wanted '%v' for '%s' and '%t'", value, test.Expected, test.TestCat, test.TestParent)
		}
	}
}

func TestCategoryName(t *testing.T) {
	var tests = []struct {
		TestCat    string
		TestSubCat string
		Expected   string
	}{
		{
			TestCat:    "",
			TestSubCat: "",
			Expected:   "",
		},
		{
			TestCat:    "d",
			TestSubCat: "s",
			Expected:   "",
		},
		{
			TestCat:    "3",
			TestSubCat: "",
			Expected:   "anime",
		},
		{
			TestCat:    "3",
			TestSubCat: "6",
			Expected:   "anime_raw",
		},
	}

	for _, test := range tests {
		value := categoryName(test.TestCat, test.TestSubCat)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function categoryName, got '%s', wanted '%s' for '%s' and '%s'", value, test.Expected, test.TestCat, test.TestSubCat)
		}
	}
}

func TestLanguageName(t *testing.T) {
	var tests = []struct {
		TestLang publicSettings.Language
		Expected string
	}{
		{
			TestLang: publicSettings.Language{"", "", ""},
			Expected: "",
		},
		{
			TestLang: publicSettings.Language{"", "fr", "fr-fr"},
			Expected: "French (France)",
		},
		{
			TestLang: publicSettings.Language{"", "fr", "fr"},
			Expected: "French",
		},
		{
			TestLang: publicSettings.Language{"something, something", "es", "es, es-mx"},
			Expected: "Spanish, Mexican Spanish",
		},
	}
	T := mockupTemplateT(t)
	for _, test := range tests {
		value := languageName(test.TestLang, T)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%s', wanted '%s' for '%v'", value, test.Expected, test.TestLang)
		}
	}
}

func TestLanguageNameFromCode(t *testing.T) {
	var tests = []struct {
		TestLang string
		Expected string
	}{
		{
			TestLang: "",
			Expected: "",
		},
		{
			TestLang: "fr-fr",
			Expected: "French (France)",
		},
		{
			TestLang: "ofjd",
			Expected: "",
		},
		{
			TestLang: "fr",
			Expected: "French",
		},
		{
			TestLang: "es, es-mx",
			Expected: "Spanish, Mexican Spanish",
		},
	}
	T := mockupTemplateT(t)
	for _, test := range tests {
		value := languageNameFromCode(test.TestLang, T)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%s', wanted '%s' for '%s'", value, test.Expected, test.TestLang)
		}
	}
}

func TestFileSize(t *testing.T) {
	var tests = []struct {
		TestSize int64
		Expected template.HTML
	}{
		{
			TestSize: 0,
			Expected: template.HTML("Unknown"),
		},
		{
			TestSize: 10,
			Expected: template.HTML("10.0 B"),
		},
	}
	T := mockupTemplateT(t)
	for _, test := range tests {
		value := fileSize(test.TestSize, T)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%s', wanted '%s' for '%d'", value, test.Expected, test.TestSize)
		}
	}
}

func TestLastID(t *testing.T) {
	var tests = []struct {
		TestTorrents []models.TorrentJSON
		TestURL      string
		Expected     int
	}{
		{
			TestTorrents: []models.TorrentJSON{{ID: 3}, {ID: 1}},
			TestURL:      "?sort=&order=",
			Expected:     3,
		},
		{
			TestTorrents: []models.TorrentJSON{{ID: 3}, {ID: 1}},
			TestURL:      "?sort=2&order=",
			Expected:     3,
		},
		{
			TestTorrents: []models.TorrentJSON{{ID: 1}, {ID: 3}},
			TestURL:      "?sort=2&order=true",
			Expected:     3,
		},
		{
			TestTorrents: []models.TorrentJSON{{ID: 1}, {ID: 3}},
			TestURL:      "?sort=3&order=true",
			Expected:     0,
		},
		{
			TestTorrents: []models.TorrentJSON{},
			TestURL:      "?sort=2&order=true",
			Expected:     0,
		},
		{
			TestTorrents: []models.TorrentJSON{},
			TestURL:      "?sort=2&order=false",
			Expected:     0,
		},
	}
	for _, test := range tests {
		url, _ := url.Parse(test.TestURL)
		value := lastID(url, test.TestTorrents)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%d', wanted '%d' for '%s' and '%v'", value, test.Expected, test.TestURL, test.TestTorrents)
		}
	}
}

func TestGetReportDescription(t *testing.T) {
	var tests = []struct {
		TestDesc string
		Expected string
	}{
		{
			TestDesc: "",
			Expected: "",
		},
		{
			TestDesc: "illegal",
			Expected: "Illegal content",
		},
		{
			TestDesc: "spam",
			Expected: "Spam / Garbage",
		},
		{
			TestDesc: "wrongcat",
			Expected: "Wrong category",
		},
		{
			TestDesc: "dup",
			Expected: "Duplicate / Deprecated",
		},
		{
			TestDesc: "illegal_content",
			Expected: "Illegal content",
		},
		{
			TestDesc: "spam_garbage",
			Expected: "Spam / Garbage",
		},
		{
			TestDesc: "wrong_category",
			Expected: "Wrong category",
		},
		{
			TestDesc: "duplicate_deprecated",
			Expected: "Duplicate / Deprecated",
		},
	}
	T := mockupTemplateT(t)
	for _, test := range tests {
		value := getReportDescription(test.TestDesc, T)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%s', wanted '%s' for '%s'", value, test.Expected, test.TestDesc)
		}
	}
}

func TestGenUploaderLink(t *testing.T) {
	var tests = []struct {
		TestID     uint
		TestName   template.HTML
		TestHidden bool
		Expected   template.HTML
	}{
		{
			TestID:     0,
			TestName:   template.HTML(""),
			TestHidden: false,
			Expected:   template.HTML("れんちょん"),
		},
		{
			TestID:     10,
			TestName:   template.HTML("dd"),
			TestHidden: true,
			Expected:   template.HTML("れんちょん"),
		},
		{
			TestID:     10,
			TestName:   template.HTML("dd"),
			TestHidden: false,
			Expected:   template.HTML("<a href=\"/user/10/dd\">dd</a>"),
		},
		{
			TestID:     0, // Old Uploader
			TestName:   template.HTML("dd"),
			TestHidden: false,
			Expected:   template.HTML("dd"),
		},
		{
			TestID:     10,
			TestName:   template.HTML(""),
			TestHidden: false,
			Expected:   template.HTML("れんちょん"),
		},
		{
			TestID:     10,
			TestName:   template.HTML(""),
			TestHidden: true,
			Expected:   template.HTML("れんちょん"),
		},
	}
	for _, test := range tests {
		value := genUploaderLink(test.TestID, test.TestName, test.TestHidden)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%s', wanted '%s' for '%d' and '%s' and '%t'", string(value), string(test.Expected), test.TestID, string(test.TestName), test.TestHidden)
		}
	}
}

func TestContains(t *testing.T) {
	var tests = []struct {
		TestArr  interface{}
		TestComp string
		Expected bool
	}{
		{
			TestArr:  "kilo",
			TestComp: "kilo",
			Expected: true,
		},
		{
			TestArr:  "kilo",
			TestComp: "loki", // Clearly not the same level
			Expected: false,
		},
		{
			TestArr:  "kilo",
			TestComp: "kiloo",
			Expected: false,
		},
		{
			TestArr:  publicSettings.Language{Code: "kilo"},
			TestComp: "kilo",
			Expected: true,
		},
		{
			TestArr:  publicSettings.Language{Code: "kilo"},
			TestComp: "loki", // Clearly not the same level
			Expected: false,
		},
		{
			TestArr:  publicSettings.Language{Code: "kilo"},
			TestComp: "kiloo",
			Expected: false,
		},
		{
			TestArr:  "kilo",
			TestComp: "",
			Expected: false,
		},
		{
			TestArr:  publicSettings.Language{Code: "kilo"},
			TestComp: "",
			Expected: false,
		},
	}
	for _, test := range tests {
		value := contains(test.TestArr, test.TestComp)
		if value != test.Expected {
			t.Errorf("Unexpected value from the function languageName, got '%t', wanted '%t' for '%v' and '%s'", value, test.Expected, test.TestArr, test.TestComp)
		}
	}
}

func testTorrentFileExists(t *testing.T) {
	var tests = []struct {
		hash 	     string
		Expected     bool
	}{
		{
			hash: "",
			Expected: false,
		},
	}
	for _, test := range tests {
		value := torrentFileExists(test.hash, "")
		if value != test.Expected {
			t.Errorf("Unexpected value from the function TorrentFileExists, got  '%t', wanted '%t' for '%s'", value, test.Expected, test.hash)
		}
	}	
}

func Testkilo_strcmp(t *testing.T) {
 	var tests = []struct {
 		TestString  string
 		TestString2 string
 		Expected bool
 	}{
 		{
 			TestString:  "kilo",
 			TestString2: "kilo",
			Expected: true,
 		},
 		{
 		TestString:  "kilo",
 			TestString2: "loki", // Clearly not the same level
 			Expected: false,
 		},
 	}
 	for _, test := range tests {
 		value := kilo_strcmp(test.TestString, test.TestString2, -1, 0)
 		if value != test.Expected {
 			t.Errorf("Unexpected value from the function languageName, got '%t', wanted '%t'", value, test.Expected, test.TestString, test.TestString)
 		}
	}
 }

 func TestToString(t *testing.T) {
 	var tests = []struct {
 		TestInt  int
 		Expected string
 	}{
 		{
 			TestInt:  0,
			Expected: "0",
 		},
 	}
 	for _, test := range tests {
		value := toString(test.TestInt)
 		if value != test.Expected {
 			t.Errorf("Unexpected value from the function languageName, got '%t', wanted '%t'", value, test.Expected)
 		}
	}
 }
 
 func Testkilo_strfind(t *testing.T) {
 	var tests = []struct {
 		TestString  string
 		TestString2 string
 		Expected bool
 	}{
 		{
 			TestString:  "kilo",
 			TestString2: "kilo",
			Expected: true,
 		},
 		{
 			TestString:  "kilo",
 			TestString2: "loki", // Clearly not the same level
 			Expected: false,
 		},
 	}
 	for _, test := range tests {
 		value := kilo_strfind(test.TestString, test.TestString2, 0)
 		if value != test.Expected {
 			t.Errorf("Unexpected value from the function languageName, got '%t', wanted '%t'", value, test.Expected, test.TestString, test.TestString)
 		}
	}
 }

func TestRand(t *testing.T) {
 	var tests = []struct {
 		TestInt  int
 		TestInt2 int
 		Expected int
 	}{
 		{
 			TestInt:  0,
 			TestInt2:  1,
			Expected: 1,
 		},
 	}
 	for _, test := range tests {
		value := kilo_rand(1)
 		if value != test.Expected {
 			//t.Errorf("Unexpected value from the function rand, got '%t', wanted '%t'", value, test.Expected)
 		}
	}
 }
 
 func TestGetDomain(t *testing.T) {
 	var tests = []struct {
 		domainName string
 	}{
 		{
 			domainName:  "wubwub",
 		},
 	}
 	for _, test := range tests {
		value := getDomainName()
 		if value != test.domainName {
 			//t.Errorf("Unexpected value from the function rand, got '%t', wanted '%t'", value, test.domainName)
 		}
	}
 }
 
 func TestGetTheme(t *testing.T) {
 	var tests = []struct {
 		domainName []string
 	}{
 		{
 			domainName:  []string{"test", "test", "test"},
 		},
 	}
 	for _, test := range tests {
		test.domainName = getThemeList()
	}
 }
 
  func testformatThemeName(t *testing.T) {
 	var tests = []struct {
 		domainName string
 	}{
 		{
 			domainName:  "test",
 		},
 	}
 	for _, test := range tests {
		value := formatThemeName("path")
 		if value != test.domainName {
 			
 		}
	}
 }

 
func mockupTemplateT(t *testing.T) publicSettings.TemplateTfunc {
	conf := config.Get().I18n
	conf.Directory = path.Join("..", conf.Directory)
	var retriever publicSettings.UserRetriever // not required during initialization

	err := publicSettings.InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}

	Ts, _, err := publicSettings.TfuncAndLanguageWithFallback("en-us")
	if err != nil {
		t.Error("Couldn't load language files!")
	}
	var T publicSettings.TemplateTfunc
	T = func(id string, args ...interface{}) template.HTML {	
		return template.HTML(fmt.Sprintf(Ts(id), args...))
	}
	return T
}
	
