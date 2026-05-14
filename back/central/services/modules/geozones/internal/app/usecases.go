package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/errors"
)

func (u *UseCase) Create(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error) {
	if !validType(dto.Type) {
		return nil, domainerrors.ErrInvalidType
	}
	if len(dto.Geometry) == 0 {
		return nil, domainerrors.ErrInvalidGeometry
	}
	return u.repo.Create(ctx, dto)
}

func (u *UseCase) BulkImport(ctx context.Context, dto dtos.BulkImportDTO) (*dtos.BulkImportResult, error) {
	result := &dtos.BulkImportResult{}
	codeIndex := make(map[string]uint)

	for i, f := range dto.Features {
		if !validType(f.Type) {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("feature[%d]: tipo invalido %q", i, f.Type))
			continue
		}
		if len(f.Geometry) == 0 {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("feature[%d]: geometria vacia", i))
			continue
		}

		var parentID *uint
		if f.ParentCode != nil && *f.ParentCode != "" {
			if id, ok := codeIndex[*f.ParentCode]; ok {
				parentID = &id
			} else {
				existing, err := u.repo.GetByCode(ctx, dto.BusinessID, parentTypeFor(f.Type), *f.ParentCode)
				if err == nil && existing != nil {
					parentID = &existing.ID
					codeIndex[*f.ParentCode] = existing.ID
				}
			}
		}

		created, err := u.repo.Create(ctx, dtos.CreateGeozoneDTO{
			BusinessID: dto.BusinessID,
			ParentID:   parentID,
			Type:       f.Type,
			Code:       f.Code,
			Name:       f.Name,
			Geometry:   f.Geometry,
			Properties: f.Properties,
		})
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("feature[%d] %q: %v", i, f.Name, err))
			continue
		}
		if f.Code != nil {
			codeIndex[*f.Code] = created.ID
		}
		result.Created++
	}
	return result, nil
}

func (u *UseCase) Get(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error) {
	return u.repo.GetByID(ctx, id, includeGeom)
}

func (u *UseCase) List(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	return u.repo.List(ctx, params)
}

func (u *UseCase) Lookup(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error) {
	return u.repo.LookupContaining(ctx, params)
}

func (u *UseCase) Delete(ctx context.Context, id uint) error {
	return u.repo.Delete(ctx, id)
}

func validType(t string) bool {
	switch t {
	case dtos.TypeCountry, dtos.TypeState, dtos.TypeCity,
		dtos.TypeAdminDistrict, dtos.TypeLocality, dtos.TypeNeighborhood, dtos.TypeBarrio, dtos.TypeCustom:
		return true
	}
	return false
}

func parentTypeFor(t string) string {
	switch t {
	case dtos.TypeState:
		return dtos.TypeCountry
	case dtos.TypeCity:
		return dtos.TypeState
	case dtos.TypeLocality:
		return dtos.TypeCity
	case dtos.TypeNeighborhood:
		return dtos.TypeLocality
	}
	return ""
}
