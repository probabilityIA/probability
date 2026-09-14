package ports

import "context"

type MediaInfo struct {
	URL      string
	MimeType string
	FileSize int64
}

type IMediaAPI interface {
	UploadMedia(ctx context.Context, phoneNumberID uint, accessToken, filename, mimeType string, data []byte) (string, error)
	SendMediaMessage(ctx context.Context, phoneNumberID uint, accessToken, to, mediaType, mediaID, caption, filename string) (string, error)
	GetMediaInfo(ctx context.Context, mediaID, accessToken string) (*MediaInfo, error)
	DownloadMedia(ctx context.Context, url, accessToken string) ([]byte, error)
}

type IChatMediaStorage interface {
	Put(ctx context.Context, key string, body []byte, contentType string) error
	Delete(ctx context.Context, key string) error
}
