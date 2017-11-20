package search

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/NyaaPantsu/nyaa/models/users"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/gin-gonic/gin"
)

// TorrentParam defines all parameters that can be provided when searching for a torrent
type TorrentParam struct {
	Full      bool // True means load all members
	Order     bool // True means ascending
	Hidden    bool // True means filter hidden torrents
	Locked    bool // False means filter locked torrents
	Deleted   bool // False means filter deleted torrents
	Status    Status
	Sort      SortMode
	Category  Categories
	Max       maxType
	Offset    uint32
	UserID    uint32
	TorrentID []uint32
	FromID    uint32
	FromDate  DateFilter
	ToDate    DateFilter
	NotNull   string // csv
	NameLike  string // csv
	NameSearch string //Contains what NameLike contains but without the excluded keywords, not used for search, just for page title
	Languages publicSettings.Languages
	MinSize   SizeBytes
	MaxSize   SizeBytes
	// Tags search
	AnidbID      uint32
	VndbID       uint32
	VgmdbID      uint32
	Dlsite       string
	VideoQuality string
	Tags         Tags
	Abort        bool
}

// Identifier returns a unique identifier for the struct
func (p *TorrentParam) Identifier() string {
	cats := ""
	for _, v := range p.Category {
		cats += fmt.Sprintf("%d%d", v.Main, v.Sub)
	}
	languages := ""
	for _, v := range p.Languages {
		languages += fmt.Sprintf("%s%s", v.Code, v.Name)
	}
	ids := ""
	for _, v := range p.TorrentID {
		ids += fmt.Sprintf("%d", v)
	}
	// Tags identifier
	tags := strings.Join(p.Tags, ",")
	tags += p.VideoQuality
	dbids := fmt.Sprintf("%d%d%d%s", p.AnidbID, p.VndbID, p.VgmdbID, p.Dlsite)

	identifier := fmt.Sprintf("%s%s%s%d%d%d%d%d%d%d%s%s%s%d%s%s%s%t%t%t%t%t", p.NameLike, p.NotNull, languages, p.Max, p.Offset, p.FromID, p.MinSize, p.MaxSize, p.Status, p.Sort, dbids, p.FromDate, p.ToDate, p.UserID, ids, cats, tags, p.Full, p.Order, p.Hidden, p.Locked, p.Deleted)
	return base64.URLEncoding.EncodeToString([]byte(identifier))
}

func parseUInt(c *gin.Context, key string) uint32 {
	// Get the user id from the url
	u64, err := strconv.ParseUint(c.Query(key), 10, 32)
	if err != nil {
		// if you can't convert it, you set it to 0
		u64 = 0
	}
	return uint32(u64)
}
func parseTorrentID(c *gin.Context) (uint32, []uint32) {
	// Get the torrent ID to limit the results to the ones after this torrent
	fromID, err := strconv.ParseUint(c.Query("fromID"), 10, 32)
	if err != nil {
		// if you can't convert it, you set it to 0
		fromID = 0
	}
	var torrentIDs []uint32
	ids := c.QueryArray("id")

	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err == nil {
			torrentIDs = append(torrentIDs, uint32(idInt))
		}
	}
	return uint32(fromID), torrentIDs
}

func parseOrder(c *gin.Context) bool {
	if c.Query("order") == "true" {
		return true
	}
	return false
}

// FromRequest : parse a request in torrent param
// TODO Should probably return an error ?
func (p *TorrentParam) FromRequest(c *gin.Context) {
	// Search by name
	// We take the search arguments from "q" in url
	p.NameLike = strings.TrimSpace(c.Query("q"))
	
	for _, word := range strings.Fields(p.NameLike) {
		if word[0] != '-' {
			p.NameSearch += word + " "
		}
	}

	// Maximum results returned
	// We take the maxximum results to display from "limit" in url
	p.Max.Parse(c.Query("limit"))

	// Limit search to one user
	// Get the user id from the url
	p.UserID = parseUInt(c, "userID")

	// if userID is not provided and username is, we try to find the user ID with the username
	if username := c.Query("user"); username != "" && p.UserID == 0 {
		log.Info(fmt.Sprint(username[0]))
		if username[0] == '#' {
			log.Info(username[1:])
			u64, err := strconv.ParseUint(username[1:], 10, 32)
			if err == nil {
				p.UserID = uint32(u64)
			}
		} else {
			user, _, _, err := users.FindByUsername(username)
			if err == nil {
				p.UserID = uint32(user.ID)
			} else {
				p.Abort = true
			}
		}
		// For other functions, we need to set userID in the request query
		q := c.Request.URL.Query()
		q.Set("userID", fmt.Sprintf("%d", p.UserID))
		c.Request.URL.RawQuery = q.Encode()
	}

	// Limit search to DbID
	// Get the id from the url
	p.AnidbID = parseUInt(c, "anidb")
	p.VndbID = parseUInt(c, "vndb")
	p.VgmdbID = parseUInt(c, "vgm")
	p.Dlsite = c.Query("dlsite")

	// Limit search to video quality
	// Get the video quality from url
	p.VideoQuality = c.Query("vq")

	// Limit search to some accepted tags
	// Get the tags from the url
	p.Tags.Parse(c.Query("tags"))

	// Order to return the results
	// Getting the order from the "order" argument in url, we default to descending order
	p.Order = parseOrder(c)

	// Limit to some status the results
	// helper to parse status from the "s" argument in url
	p.Status.Parse(c.Query("s"))

	// Sort the results
	// Parse the sorting mode of the result from the "sort" argument in url
	p.Sort.Parse(c.Query("sort"))

	// Set NoNull to improve pg query
	if p.Sort == Date {
		p.NotNull = p.Sort.ToDBField() + " IS NOT NULL"
	}
	// Category in which you have to search
	// Parse the categories from the "c" argument in url
	p.Category = ParseCategories(c.Query("c"))

	// Languages filter of the torrents
	// We get the languages filtering the results from the "lang" argument in url
	p.Languages = ParseLanguages(c.QueryArray("lang"))

	// From which date you need to search and  To which date you need to search
	// maxage is an int parameter limiting the results to the last "x" days (old nyaa behavior)
	p.FromDate, p.ToDate = backwardCompatibility(c.Query("maxage"), c.Query("fromDate"), c.Query("toDate"), c.Query("dateType"))

	// Parsing minimum and maximum size from the sizeType given (minSize & maxSize & sizeType in url)
	// Minimum size to search
	p.MinSize.Parse(c.Query("minSize"), c.Query("sizeType"))
	// Maximum size to search
	p.MaxSize.Parse(c.Query("maxSize"), c.Query("sizeType"))

	// Needed to display result after a certain torrentID or to limit results to some torrent IDs
	p.FromID, p.TorrentID = parseTorrentID(c)
}

// toESQuery : Builds a query string with for es query string query defined here
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html
func (p *TorrentParam) toESQuery(c *gin.Context) *Query {
	query := &Query{
		TorrentParam: p,
	}

	if len(p.Category) > 0 {
		query.Append(p.Category.ToESQuery())
	}

	if c.Query("userID") != "" {
		if !strings.Contains(c.Query("userID"), ",") {
			if p.UserID > 0 {
				query.Append("uploader_id:" + strconv.FormatInt(int64(p.UserID), 10))
				if p.Hidden {
					query.Append("hidden:false")
				}
			} else if p.UserID == 0 {
				query.Append(fmt.Sprintf("(uploader_id: %d OR hidden:%t)", p.UserID, true))
			}
		} else {
			userIDs := strings.Split(c.Query("userID"), ",")
			for _, str := range userIDs {
				if userID, err := strconv.Atoi(str); err == nil && userID >= 0 {
					if userID > 0 {
						query.Append(fmt.Sprintf("(uploader_id:%d AND hidden:false)", userID))
					} else {
						query.Append(fmt.Sprintf("(uploader_id:%d OR hidden:%t)", userID, true))
					}
				}
			}
		}
	}

	if c.Query("nuserID") != "" {
		nuserID := strings.Split(c.Query("nuserID"), ",")
		for _, str := range nuserID {
			if userID, err := strconv.Atoi(str); err == nil && userID >= 0 {
				if userID > 0 {
					query.Append(fmt.Sprintf("NOT(uploader_id:%d AND hidden:false)", userID))
				} else {
					query.Append(fmt.Sprintf("NOT(uploader_id:%d AND hidden:false) !hidden:true", userID))
				}
			}
		}
	}

	if p.Status != ShowAll {
		query.Append(p.Status.ToESQuery())
	}
	if !p.Locked {
		query.Append("!status:5")
	}


	if p.FromID != 0 {
		query.Append("id:>" + strconv.FormatInt(int64(p.FromID), 10))
	}

	if len(p.TorrentID) > 0 {
		for _, id := range p.TorrentID {
			query.Append(fmt.Sprintf("id:%d", id))
		}
	}

	if p.FromDate != "" || p.ToDate != "" {
		query.Append("date: [" + p.FromDate.ToESQuery() + " " + p.ToDate.ToESQuery() + "]")
	}

	if p.MinSize > 0 || p.MaxSize > 0 {
		query.Append("filesize: [" + p.MinSize.ToESQuery() + " " + p.MaxSize.ToESQuery() + "]")
	}

	if len(p.Languages) > 0 {
		langsToESQuery(query, p.Languages)
	}

	// Tags search
	// Anidb
	if p.AnidbID != 0 {
		query.Append("anidbid:" + strconv.FormatInt(int64(p.AnidbID), 10))
	}
	// Vndb
	if p.VndbID != 0 {
		query.Append("vndbid:" + strconv.FormatInt(int64(p.VndbID), 10))
	}
	// Vgmdb
	if p.VgmdbID != 0 {
		query.Append("vgmdbid:" + strconv.FormatInt(int64(p.VgmdbID), 10))
	}
	// Dlsite
	if p.Dlsite != "" {
		query.Append("dlsite:" + p.Dlsite)
	}
	// Video quality
	if p.VideoQuality != "" {
		query.Append("videoquality:" + p.VideoQuality)
	}
	// Other tags
	if len(p.Tags) > 0 {
		query.Append(p.Tags.ToESQuery())
	}

	return query
}

// FindES :
/* Uses elasticsearch to find the torrents based on TorrentParam
 */
func (p *TorrentParam) FindES(c *gin.Context, client *elastic.Client) ([]models.Torrent, int64, error) {
	search, err := p.toESQuery(c).ToESQuery(client)
	if err != nil {
		return nil, 0, err
	}

	result, err := search.Do(c)
	if err != nil {
		return nil, 0, err
	}

	log.Infof("Query '%s' took %d milliseconds.", p.NameLike, result.TookInMillis)
	log.Infof("Amount of results %d.", result.TotalHits())

	var torrents []models.Torrent
	var torrentCount int
	if len(result.Hits.Hits) <= 0 {
		return nil, 0, nil
	}
	for _, hit := range result.Hits.Hits {
		var tJSON models.TorrentJSON
		err := json.Unmarshal(*hit.Source, &tJSON)
		if err == nil {
			torrents = append(torrents, tJSON.ToTorrent())
			torrentCount++
		} else {
			log.Errorf("Cannot unmarshal elasticsearch torrent: %s", err)
		}
	}
	if torrentCount < len(result.Hits.Hits) {
		log.Errorf("Only %d / %d parsed correctly, see error above", torrentCount, len(result.Hits.Hits))
	}

	return torrents, result.TotalHits(), nil
}

func (p *TorrentParam) toDBQuery(c *gin.Context) *Query {
	query := &Query{}

	sql, cats := p.Category.ToDBQuery()
	query.Append(sql, cats...)

	if len(p.Languages) > 0 {
		query.Append("language "+searchOperator, "%"+langsToDBQuery(p.Languages)+"%")
	}

	if c.Query("userID") != "" {
		if !strings.Contains(c.Query("userID"), ",") {
			if p.UserID > 0 {
				query.Append("uploader", p.UserID)
				if p.Hidden {
					query.Append("hidden", false)
				}
			} else if p.UserID == 0 {
				query.Append("(uploader = ? OR hidden = ?)", p.UserID, true)
			}
		} else {
			userIDs := strings.Split(c.Query("userID"), ",")
			for _, str := range userIDs {
				if userID, err := strconv.Atoi(str); err == nil && userID >= 0 {
					if userID > 0 {
						query.Append("(uploader = ? AND hidden = ?)", userID, false)
					} else {
						query.Append("(uploader = ? OR hidden = ?)", userID, true)
					}
				}
			}
		}
	}

	if c.Query("nuserID") != "" {
		nuserID := strings.Split(c.Query("nuserID"), ",")
		for _, str := range nuserID {
			if userID, err := strconv.Atoi(str); err == nil && userID >= 0 {
				if userID > 0 {
					query.Append("NOT(uploader = ? AND hidden = ?)", userID, false)
				} else {
					query.Append("uploader <> ? AND hidden != ?", userID, true)
				}
			}
		}
	}

	if p.FromID != 0 {
		query.Append(config.Get().Models.TorrentsTableName + ".torrent_id > ?", p.FromID)
	}
	if len(p.TorrentID) > 0 {
		for _, id := range p.TorrentID {
			query.Append(config.Get().Models.TorrentsTableName + ".torrent_id = ?", id)
		}
	}
	if p.FromDate != "" {
		query.Append("date >= ?", p.FromDate.ToDBQuery())
	}
	if p.ToDate != "" {
		query.Append("date <= ?", p.ToDate.ToDBQuery())
	}
	if p.Status != 0 {
		query.Append(p.Status.ToDBQuery())
	}
	if !p.Locked {
		query.Append("status IS NOT ?", 5)
	}

	if len(p.NotNull) > 0 {
		query.Append(p.NotNull)
	}
	if p.MinSize > 0 {
		query.Append("filesize >= ?", p.MinSize.ToDBQuery())
	}
	if p.MaxSize > 0 {
		query.Append("filesize <= ?", p.MaxSize.ToDBQuery())
	}

	// Tags search
	// Anidb
	if p.AnidbID > 0 {
		query.Append("anidbid = ?", p.AnidbID)
	}
	// Vndb
	if p.VndbID > 0 {
		query.Append("vndbid = ?", p.VndbID)
	}
	// Vgmdb
	if p.VgmdbID > 0 {
		query.Append("vgmdbid = ?", p.VgmdbID)
	}
	// Dlsite
	if p.Dlsite != "" {
		query.Append("dlsite = ?", p.Dlsite)
	}
	// Video quality
	if p.VideoQuality != "" {
		query.Append("videoquality = ?", p.VideoQuality)
	}
	// Other tags
	if len(p.Tags) > 0 {
		query.Append(p.Tags.ToDBQuery())
	}

	querySplit := strings.Fields(p.NameLike)
	for _, word := range querySplit {
		if word[0] == '-' && len(word) > 1 {
			//Exclude words starting with -
			query.Append("torrent_name NOT "+searchOperator, "%"+word[1:]+"%")
			continue
		}
		
		firstRune, _ := utf8.DecodeRuneInString(word)
		if len(word) == 1 && unicode.IsPunct(firstRune) {
			// some queries have a single punctuation character
			// which causes a full scan instead of using the index
			// and yields no meaningful results.
			// due to len() == 1 we're just looking at 1-byte/ascii
			// punctuation characters.
			continue
		}

		if useTSQuery && stringIsASCII(word) {
			query.Append("torrent_name @@ plainto_tsquery(?)", word)
		} else {
			// TODO: possible to make this faster?
			query.Append("torrent_name "+searchOperator, "%"+word+"%")
		}
	}
	return query
}

// FindDB :
/* Uses SQL to find the torrents based on TorrentParam
 */
func (p *TorrentParam) FindDB(c *gin.Context) ([]models.Torrent, int64, error) {
	orderBy := p.Sort.ToDBField()
	query := p.toDBQuery(c)
	orderBy += " "

	switch p.Order {
	case true:
		orderBy += "asc"
		if models.ORM.Dialect().GetName() == "postgres" {
			orderBy += " NULLS FIRST"
		}
	case false:
		orderBy += "desc"
		if models.ORM.Dialect().GetName() == "postgres" {
			orderBy += " NULLS LAST"
		}
	}

	log.Infof("SQL query is :: %s\n", query.String())

	if p.Deleted {
		tor, count, err := torrents.FindDeleted(query, orderBy, int(p.Max), int(uint32(p.Max)*(p.Offset-1)))
		return tor, int64(count), err
	} else if p.Full {
		tor, count, err := torrents.FindWithUserOrderBy(query, orderBy, int(p.Max), int(uint32(p.Max)*(p.Offset-1)))
		return tor, int64(count), err
	}
	tor, count, err := torrents.FindOrderBy(query, orderBy, int(p.Max), int(uint32(p.Max)*(p.Offset-1)))
	return tor, int64(count), err
}

// Clone : To clone a torrent params
func (p *TorrentParam) Clone() TorrentParam {
	return TorrentParam{
		Order:     p.Order,
		Status:    p.Status,
		Sort:      p.Sort,
		Category:  p.Category,
		Max:       p.Max,
		Offset:    p.Offset,
		UserID:    p.UserID,
		TorrentID: p.TorrentID,
		FromID:    p.FromID,
		FromDate:  p.FromDate,
		ToDate:    p.ToDate,
		NotNull:   p.NotNull,
		NameLike:  p.NameLike,
		Languages: p.Languages,
		MinSize:   p.MinSize,
		MaxSize:   p.MaxSize,
		Locked:    p.Locked,
	}
}
