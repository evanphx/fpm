package fpm

type MetaData struct {
	Name         string
	Description  string
	Vendor       string
	Version      string
	License      string
	Architecture string
	Maintainer   string
	Section      string
	URL          string
	Extra        map[string]interface{}
}

func NewMetaData() *MetaData {
	md := &MetaData{
		Name:         "unknown",
		Description:  "none",
		Vendor:       "fpm",
		License:      "unknown",
		Architecture: "all",
		Maintainer:   "noone",
		Section:      "unknown",
		URL:          "http://nowhere.com",
		Extra:        make(map[string]interface{}),
	}

	return md
}
