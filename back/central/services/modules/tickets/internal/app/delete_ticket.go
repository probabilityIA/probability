package app

import "context"

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	attachments, err := uc.repo.ListAttachments(ctx, id)
	if err != nil {
		uc.log.Warn().Err(err).Uint("ticket_id", id).Msg("delete ticket: failed to list attachments for s3 cleanup")
	}
	for _, att := range attachments {
		if att.FileURL == "" {
			continue
		}
		if err := uc.storage.DeleteFile(ctx, att.FileURL); err != nil {
			uc.log.Warn().Err(err).Uint("attachment_id", att.ID).Str("file_url", att.FileURL).Msg("delete ticket: failed to delete s3 file")
		}
	}
	return uc.repo.Delete(ctx, id)
}
