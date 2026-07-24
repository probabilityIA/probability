package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

const syncBatchSize = 1000

func (uc *UseCase) Sync(ctx context.Context) (*dtos.SyncResult, error) {
	concepts, err := uc.repo.ListConcepts(ctx)
	if err != nil {
		return nil, err
	}
	result := &dtos.SyncResult{BySource: map[string]int{}}
	for i := range concepts {
		concept := concepts[i]
		if !concept.IsAutomatic || !concept.IsActive || concept.SourceType == entities.SourceManual {
			continue
		}
		created, err := uc.syncConcept(ctx, &concept)
		if err != nil {
			uc.log.Error(ctx).Err(err).Str("source_type", concept.SourceType).Msg("Accounting sync: source failed")
			continue
		}
		if created > 0 {
			result.BySource[concept.SourceType] = created
			result.Created += created
			uc.log.Info(ctx).Str("source_type", concept.SourceType).Int("created", created).Msg("Accounting sync: entries created")
		}
	}
	return result, nil
}

func (uc *UseCase) syncConcept(ctx context.Context, concept *entities.Concept) (int, error) {
	created := 0
	for {
		candidates, err := uc.repo.FindSyncCandidates(ctx, concept.SourceType, syncBatchSize)
		if err != nil {
			return created, err
		}
		if len(candidates) == 0 {
			return created, nil
		}
		roundInserted := 0
		for _, cand := range candidates {
			entry := &entities.Entry{
				ConceptID:   concept.ID,
				BusinessID:  cand.BusinessID,
				EntryDate:   cand.EntryDate,
				Amount:      round2(cand.Amount),
				Kind:        concept.Kind,
				SourceType:  cand.SourceType,
				SourceID:    cand.SourceID,
				Description: cand.Description,
				IsAutomatic: true,
			}
			applyTaxes(entry, concept)
			inserted, err := uc.repo.CreateEntry(ctx, entry)
			if err != nil {
				uc.log.Error(ctx).Err(err).
					Str("source_type", cand.SourceType).
					Str("source_id", cand.SourceID).
					Msg("Accounting sync: failed to insert entry")
				continue
			}
			if inserted {
				created++
				roundInserted++
			}
		}
		if len(candidates) < syncBatchSize || roundInserted == 0 {
			return created, nil
		}
	}
}
