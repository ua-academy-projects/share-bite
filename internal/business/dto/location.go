package dto

type CreateLocationInput struct {
	Name        string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
}

type UpdateLocationInput struct {
	Name        *string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
}

type ListNearbyVenuesInput struct {
	Lat   float64 `form:"lat" binding:"required,latitude"`
	Lon   float64 `form:"lon" binding:"required,longitude"`
	Limit int     `form:"limit" binding:"min=1,max=100"`
	Skip  int     `form:"skip" binding:"min=0"`
}

type NearbyVenueItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Avatar   *string `json:"avatar"`
	Distance float64 `json:"distance"`
}

type ListNearbyVenuesOutput struct {
	Items []NearbyVenueItem `json:"items"`
	Total int               `json:"total"`
}
