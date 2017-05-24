package common

// TorrentParam defines all parameters that can be provided when searching for a torrent
type TorrentParam struct {
	All       bool // True means ignore everything but Max and Offset
	Full      bool // True means load all members
	Order     bool // True means acsending
	Status    Status
	Sort      SortMode
	Category  Category
	Max       uint32
	Offset    uint32
	UserID    uint32
	TorrentID uint32
	NotNull   string // csv
	Null      string // csv
	NameLike  string // csv
}

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
		NotNull:   p.NotNull,
		Null:      p.Null,
		NameLike:  p.NameLike,
	}
}
