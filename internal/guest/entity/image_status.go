package entity

type ImageProcessingStatus string

const (
	ImageStatusPending    ImageProcessingStatus = "pending"
	ImageStatusProcessing ImageProcessingStatus = "processing"
	ImageStatusCompleted  ImageProcessingStatus = "completed"
	ImageStatusFailed     ImageProcessingStatus = "failed"
)
