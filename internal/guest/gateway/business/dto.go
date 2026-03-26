package business

// businessVenueResponseDTO строго повторяет формат JSON, который отдает Business API.
// Эта структура НЕ экспортируется (начинается с маленькой буквы),
// чтобы никто за пределами пакета business не мог её использовать.
type businessVenueResponseDTO struct {
	VenueID string `json:"venue_id"`
	Status  string `json:"status"` // например: "active", "closed", "not_found"
}
