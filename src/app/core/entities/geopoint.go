package entities

type Geopoint struct {
	Lat float32 `json: latitude`
	Lon float32 `json: longitude`
}
