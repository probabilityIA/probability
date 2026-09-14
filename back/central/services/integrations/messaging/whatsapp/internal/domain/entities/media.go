package entities

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

const (
	MediaTypeImage    = "image"
	MediaTypeDocument = "document"
	MediaTypeAudio    = "audio"
	MediaTypeVideo    = "video"
	MediaTypeSticker  = "sticker"
)

const (
	MaxImageBytes      = 5 << 20
	MaxAttachmentBytes = 10 << 20
	MaxCaptionLength   = 1024
)

var ErrMediaNotAllowed = errors.New("archivo no permitido")

var outboundMediaTypes = map[string]string{
	"image/jpeg":         MediaTypeImage,
	"image/png":          MediaTypeImage,
	"video/mp4":          MediaTypeVideo,
	"video/3gpp":         MediaTypeVideo,
	"audio/aac":          MediaTypeAudio,
	"audio/mp4":          MediaTypeAudio,
	"audio/mpeg":         MediaTypeAudio,
	"audio/amr":          MediaTypeAudio,
	"audio/ogg":          MediaTypeAudio,
	"application/pdf":    MediaTypeDocument,
	"application/msword": MediaTypeDocument,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": MediaTypeDocument,
	"application/vnd.ms-excel": MediaTypeDocument,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         MediaTypeDocument,
	"application/vnd.ms-powerpoint":                                             MediaTypeDocument,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": MediaTypeDocument,
	"text/plain": MediaTypeDocument,
	"text/csv":   MediaTypeDocument,
}

var mediaExtensions = map[string]string{
	"image/jpeg":         ".jpg",
	"image/png":          ".png",
	"image/webp":         ".webp",
	"video/mp4":          ".mp4",
	"video/3gpp":         ".3gp",
	"audio/aac":          ".aac",
	"audio/mp4":          ".m4a",
	"audio/mpeg":         ".mp3",
	"audio/amr":          ".amr",
	"audio/ogg":          ".ogg",
	"application/pdf":    ".pdf",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/vnd.ms-excel": ".xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
	"application/vnd.ms-powerpoint":                                             ".ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	"text/plain": ".txt",
	"text/csv":   ".csv",
}

type OutboundMedia struct {
	Filename string
	MimeType string
	Data     []byte
}

func NormalizeMime(mimeType string) string {
	base := strings.SplitN(mimeType, ";", 2)[0]
	return strings.ToLower(strings.TrimSpace(base))
}

func ClassifyOutboundMedia(mimeType string, size int) (string, error) {
	mime := NormalizeMime(mimeType)
	mediaType, ok := outboundMediaTypes[mime]
	if !ok {
		return "", fmt.Errorf("%w: WhatsApp no acepta el formato %q", ErrMediaNotAllowed, mime)
	}
	if size <= 0 {
		return "", fmt.Errorf("%w: el archivo esta vacio", ErrMediaNotAllowed)
	}
	limit := MaxAttachmentBytes
	if mediaType == MediaTypeImage {
		limit = MaxImageBytes
	}
	if size > limit {
		return "", fmt.Errorf("%w: supera el limite de %d MB", ErrMediaNotAllowed, limit>>20)
	}
	return mediaType, nil
}

func SafeMediaFilename(filename string) string {
	name := path.Base(strings.ReplaceAll(strings.TrimSpace(filename), "\\", "/"))
	if name == "." || name == "/" {
		return ""
	}
	if len(name) > 200 {
		name = name[len(name)-200:]
	}
	return name
}

func MediaExtension(mimeType, filename string) string {
	if ext, ok := mediaExtensions[NormalizeMime(mimeType)]; ok {
		return ext
	}
	if ext := strings.ToLower(path.Ext(SafeMediaFilename(filename))); ext != "" && len(ext) <= 6 {
		return ext
	}
	return ".bin"
}

func ChatMediaKey(businessID uint, mimeType, filename, id string, now time.Time) string {
	return fmt.Sprintf("whatsapp/%d/%04d/%02d/%s%s", businessID, now.Year(), int(now.Month()), id, MediaExtension(mimeType, filename))
}

func MediaPreviewLabel(mediaType string) string {
	switch mediaType {
	case MediaTypeImage:
		return "[Imagen]"
	case MediaTypeDocument:
		return "[Documento]"
	case MediaTypeAudio:
		return "[Audio]"
	case MediaTypeVideo:
		return "[Video]"
	case MediaTypeSticker:
		return "[Sticker]"
	default:
		return ""
	}
}
