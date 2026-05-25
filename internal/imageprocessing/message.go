package imageprocessing

type ProcessImageMessage struct {
	ImageID string `json:"image_id"`
	S3Key   string `json:"s3_key"`
}
