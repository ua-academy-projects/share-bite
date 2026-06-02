package aws

import (
	"github.com/uber/h3-go/v4"
)

type H3Service interface {
	GetH3Index(lat, lon float64, resolution int) string
	GetH3Neighbors(lat, lon float64, resolution int, k int) []string
}

type h3Service struct{}

func NewH3Service() H3Service {
	return &h3Service{}
}

func (s *h3Service) GetH3Index(lat, lon float64, resolution int) string {
	if lat == 0 && lon == 0 {
		return ""
	}
	latLng := h3.NewLatLng(lat, lon)
	cell, _ := h3.LatLngToCell(latLng, resolution)
	return cell.String()
}

func (s *h3Service) GetH3Neighbors(lat, lon float64, resolution int, k int) []string {
	if lat == 0 && lon == 0 {
		return nil
	}

	latLng := h3.NewLatLng(lat, lon)
	cell, _ := h3.LatLngToCell(latLng, resolution)

	neighbors, _ := h3.GridDisk(cell, k)

	var hashes []string
	for _, n := range neighbors {
		hashes = append(hashes, n.String())
	}

	return hashes
}
