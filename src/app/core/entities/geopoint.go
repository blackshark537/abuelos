package entities

type Geopoint struct {
	Lat float32 `json:"lat" xml:"lat" form:"lat"`
	Lon float32 `json:"lon" xml:"lon" form:"lon"`
}
