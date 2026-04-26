package app

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) UploadAttachment(ctx context.Context, ticketID uint, commentID *uint, uploaderID uint, file *multipart.FileHeader) (*entities.TicketAttachment, error) {
	if _, err := uc.repo.GetByID(ctx, ticketID); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	folder := fmt.Sprintf("tickets/%d", ticketID)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
	contentType := file.Header.Get("Content-Type")

	url, err := uc.storage.UploadFile(ctx, folder, filename, buf, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload: %w", err)
	}

	return uc.repo.AddAttachment(ctx, dtos.CreateAttachmentDTO{
		TicketID:     ticketID,
		CommentID:    commentID,
		UploadedByID: uploaderID,
		FileURL:      url,
		FileName:     file.Filename,
		MimeType:     contentType,
		Size:         file.Size,
	})
}

func (uc *UseCase) ListAttachments(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
	return uc.repo.ListAttachments(ctx, ticketID)
}

func (uc *UseCase) DeleteAttachment(ctx context.Context, attachmentID uint, requesterID uint, isSuperAdmin bool) error {
	att, err := uc.repo.GetAttachment(ctx, attachmentID)
	if err != nil {
		return err
	}
	if !isSuperAdmin && att.UploadedByID != requesterID {
		return dom.ErrForbidden
	}
	if err := uc.storage.DeleteFile(ctx, att.FileURL); err != nil {
		uc.log.Warn().Err(err).Uint("attachment_id", attachmentID).Msg("failed to delete s3 file")
	}
	return uc.repo.DeleteAttachment(ctx, attachmentID)
}

func (uc *UseCase) ListHistory(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
	return uc.repo.ListHistory(ctx, ticketID)
}
