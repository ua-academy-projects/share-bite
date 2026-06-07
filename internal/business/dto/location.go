package dto

type CreateLocationInput struct {
	Name        string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
	TagIDs      []int
	H3Hash      *string
}

type UpdateLocationInput struct {
	Name        *string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
	TagIDs      *[]int
}

type ListNearbyVenuesInput struct {
	Lat   float64 `form:"lat" binding:"required,latitude"`
	Lon   float64 `form:"lon" binding:"required,longitude"`
	Limit int     `form:"limit" binding:"max=100"`
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

type VenueHoursDayInput struct {
	Weekday   int     `json:"weekday" binding:"required,min=1,max=7"`
	OpenTime  *string `json:"openTime"`
	CloseTime *string `json:"closeTime"`
}

type UpdateVenueHoursInput struct {
	Days []VenueHoursDayInput `json:"days" binding:"required,min=1,max=7"`
}

type UpdateVenueHoursOutput struct {
	VenueID int                  `json:"venueId"`
	Days    []VenueHoursDayInput `json:"days"`
}
