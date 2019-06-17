module github.com/NyaaPantsu/nyaa

go 1.12

require (
	cloud.google.com/go v0.40.0 // indirect
	github.com/BurntSushi/toml v0.3.1
	github.com/CloudyKit/fastprinter v0.0.0-20170127035650-74b38d55f37a
	github.com/CloudyKit/jet v0.0.0-20170608194317-a1c4500d4bee
	github.com/RoaringBitmap/roaring v0.4.17
	github.com/Sirupsen/logrus v0.0.0-20170504071019-5b60b3d3ee01
	github.com/Stephen304/goscrape v0.0.0-20150207020750-2d9c5bd77935
	github.com/anacrolix/dht v1.0.1
	github.com/anacrolix/go-libutp v0.0.0-20180808010927-aebbeb60ea05
	github.com/anacrolix/missinggo v1.1.0
	github.com/anacrolix/sync v0.0.0-20180808010631-44578de4e778
	github.com/anacrolix/torrent v1.2.0
	github.com/anacrolix/utp v0.0.0-20180219060659-9e0e1d1d0572
	github.com/asaskevich/govalidator v0.0.0-20170529110029-aa5cce4a76ed
	github.com/boltdb/bolt v1.3.1
	github.com/bradfitz/iter v0.0.0-20190303215204-33e6a9893b0c
	github.com/bradfitz/slice v0.0.0-20180809154707-2b758aa73013
	github.com/davecgh/go-spew v1.1.1
	github.com/dchest/captcha v0.0.0-20150728125059-9e952142169c
	github.com/dgrijalva/jwt-go v0.0.0-20170608005149-a539ee1a749a
	github.com/dustin/go-humanize v1.0.0
	github.com/edsrzf/mmap-go v1.0.0
	github.com/fatih/structs v1.1.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/frustra/bbcode v0.0.0-20150429195712-e3d2906cb269
	github.com/gin-gonic/gin v0.0.0-20170428105923-d5b353c5d5a5
	github.com/glycerine/go-unsnap-stream v0.0.0-20181221182339-f9677308dec2
	github.com/go-playground/locales v0.0.0-20170327191450-1e5f1161c641
	github.com/go-playground/universal-translator v0.0.0-20170327191703-71201497bace
	github.com/go-playground/validator v9.4.0+incompatible
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/golang/mock v1.3.1 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/golang/snappy v0.0.1
	github.com/google/btree v1.0.0
	github.com/google/pprof v0.0.0-20190515194954-54271f7e092f // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/gorilla/feeds v0.0.0-20160207162205-441264de03a8
	github.com/gorilla/securecookie v1.1.1
	github.com/jinzhu/configor v0.0.0-20170522021620-ff2ac2b1ce3d
	github.com/jinzhu/gorm v1.9.9
	github.com/jinzhu/inflection v1.0.0
	github.com/justinas/nosurf v0.0.0-20161004085251-8e1568277264
	github.com/kr/pty v1.1.5 // indirect
	github.com/lib/pq v1.1.1
	github.com/mailru/easyjson v0.0.0-20190614124828-94de47d64c63 // indirect
	github.com/majestrate/i2p-tools v0.0.0-20170507194519-afc8e46afa95
	github.com/manucorporat/sse v0.0.0-20160126180136-ee05b128a739
	github.com/mattn/go-isatty v0.0.7
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/microcosm-cc/bluemonday v0.0.0-20161202143824-e79763773ab6
	github.com/mohae/deepcopy v0.0.0-20170603005431-491d3605edfb
	github.com/moul/http2curl v0.0.0-20161031194548-4e24498b31db
	github.com/mschoch/smat v0.0.0-20160514031455-90eadee771ae
	github.com/nicksnyder/go-i18n v1.10.0
	github.com/nicksnyder/go-i18n/v2 v2.0.2
	github.com/ory/fosite v0.0.0-20170709225030-c45a37d3bb9e
	github.com/parnurzeal/gorequest v0.0.0-20170429061244-5bf13be19878
	github.com/patrickmn/go-cache v0.0.0-20170418232947-7ac151875ffb
	github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/pelletier/go-buffruneio v0.2.0
	github.com/pelletier/go-toml v0.0.0-20170504040314-97253b98df84
	github.com/philhofer/fwd v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/russross/blackfriday v0.0.0-20170509060714-0ba0f2b6ed7c
	github.com/ryszard/goskiplist v0.0.0-20150312221310-2dfbae5fcf46
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tinylib/msgp v1.1.0
	github.com/willf/bitset v1.1.10
	github.com/willf/bloom v2.0.3+incompatible
	github.com/zeebo/bencode v1.0.0
	go.opencensus.io v0.22.0 // indirect
	go4.org v0.0.0-20170505145847-16ace784e4b1
	golang.org/x/crypto v0.0.0-20190617133340-57b3e21c3d56
	golang.org/x/exp v0.0.0-20190510132918-efd6b22b2522 // indirect
	golang.org/x/image v0.0.0-20190616094056-33659d3de4f5 // indirect
	golang.org/x/mobile v0.0.0-20190607214518-6fa95d984e88 // indirect
	golang.org/x/mod v0.1.0 // indirect
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools v0.0.0-20190614205625-5aca471b1d59 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20190611190212-a7e196e89fd3 // indirect
	google.golang.org/grpc v1.21.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20161104145732-dd45e6a67c53
	gopkg.in/olivere/elastic.v5 v5.0.81
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v2 v2.2.2
	honnef.co/go/tools v0.0.0-20190614002413-cb51c254f01b // indirect
)
