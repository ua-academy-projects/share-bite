package dto

type CreateLocationInput struct {
	Name        string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
	TagIDs      []int
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
