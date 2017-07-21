package torrentValidator

// TorrentRequest struct
// Same json name as the constant!
type TorrentRequest struct {
	Name        string   `validate:"required" form:"name" json:"name,omitempty"`
	Magnet      string   `json:"magnet,omitempty" form:"magnet"`
	Category    string   `validate:"required" form:"c" json:"c"`
	Remake      bool     `json:"remake,omitempty" form:"remake"`
	Description string   `json:"desc,omitempty" form:"desc"`
	Status      int      `json:"status,omitempty" form:"status"`
	Hidden      bool     `json:"hidden,omitempty" form:"hidden"`
	CaptchaID   string   `json:"-" form:"captchaID"`
	WebsiteLink string   `validate:"uri" json:"website_link,omitempty" form:"website_link"`
	SubCategory int      `json:"sub_category,omitempty" form:"sub_category"`
	Languages   []string `json:"languages,omitempty" form:"languages"`

	Infohash      string         `json:"hash,omitempty" form:"hash"`
	CategoryID    int            `json:"-" form:"category_id"`
	SubCategoryID int            `json:"-" form:"subcategory_id"`
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

// ReassignForm : Structure for reassign Form used by the reassign page
type ReassignForm struct {
	AssignTo uint
	By       string
	Data     string

	Torrents []uint
}
