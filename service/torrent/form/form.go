package torrentform

type PanelPost struct {
	Name  string `form:"name" needed:"true" len_min:"3" len_max:"20"`
	Hash     string `form:"hash" needed:"true"`
	Category  int `form:"cat" needed:"true"`
	Sub_Category int `form:"subcat"`
	Status string `form:"status" needed:"true"`
	Description   string   `form:"desc"`
	WebsiteLink   string   `form:"website"`
}