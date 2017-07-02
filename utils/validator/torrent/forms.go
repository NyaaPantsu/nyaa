package torrentValidator

// TorrentRequest struct
// Same json name as the constant!
type TorrentRequest struct {
	Name        string   `validate:"required" json:"name,omitempty"`
	Magnet      string   `json:"magnet,omitempty"`
	Category    string   `validate:"required" json:"c"`
	Remake      bool     `json:"remake,omitempty"`
	Description string   `json:"desc,omitempty"`
	Status      int      `json:"status,omitempty"`
	Hidden      bool     `json:"hidden,omitempty"`
	CaptchaID   string   `json:"-"`
	WebsiteLink string   `validate:"uri" json:"website_link,omitempty"`
	SubCategory int      `json:"sub_category,omitempty"`
	Languages   []string `json:"languages,omitempty"`

	Infohash      string         `json:"hash,omitempty"`
	CategoryID    int            `json:"-"`
	SubCategoryID int            `json:"-"`
	Filesize      int64          `json:"filesize,omitempty"`
	Filepath      string         `json:"-"`
	FileList      []uploadedFile `json:"filelist,omitempty"`
	Trackers      []string       `json:"trackers,omitempty"`
}

// UpdateRequest struct
type UpdateRequest struct {
	ID     uint           `json:"id"`
	Update TorrentRequest `json:"update"`
}

// Use this, because we seem to avoid using models, and we would need
// the torrent ID to create the File in the DB
type uploadedFile struct {
	Path     []string `json:"path"`
	Filesize int64    `json:"filesize"`
}
