package usecasetemplates

import (
	"context"
	"strings"
)

func (u *usecase) SyncCustomStatuses(ctx context.Context, submission CustomTemplateSubmission) error {
	wanted := make(map[string]bool, len(submission.Names))
	for _, name := range submission.Names {
		if name = strings.TrimSpace(name); name != "" {
			wanted[name] = true
		}
	}
	if len(wanted) == 0 {
		return nil
	}

	wabaID, token, baseURL, err := u.resolveSubmissionTarget(ctx, submission.BusinessID)
	if err != nil {
		return err
	}

	remote, err := u.apiFactory(baseURL).ListTemplates(ctx, wabaID, token)
	if err != nil {
		return err
	}

	found := 0
	for _, template := range remote {
		if !wanted[template.Name] {
			continue
		}

		reason := strings.TrimSpace(template.RejectedReason)
		if strings.EqualFold(reason, "NONE") {
			reason = ""
		}

		if err := u.HandleStatusUpdate(ctx, wabaID, template.Name, template.Language, template.Status, reason); err != nil {
			return err
		}
		found++
	}

	u.log.Info(ctx).
		Uint("business_id", submission.BusinessID).
		Str("waba_id", wabaID).
		Int("pedidas", len(wanted)).
		Int("encontradas", found).
		Msg("Estado de plantillas consultado en Meta")

	return nil
}
