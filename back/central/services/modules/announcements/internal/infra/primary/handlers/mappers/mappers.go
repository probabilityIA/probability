package mappers

import (
	"math"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/response"
)

func CreateRequestToDTO(req request.CreateAnnouncementRequest, createdByID uint) dtos.CreateAnnouncementDTO {
	dto := dtos.CreateAnnouncementDTO{
		BusinessID:    req.BusinessID,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Message:       req.Message,
		DisplayType:   entities.DisplayType(req.DisplayType),
		FrequencyType: entities.FrequencyType(req.FrequencyType),
		Priority:      req.Priority,
		IsGlobal:      req.IsGlobal,
		CreatedByID:   createdByID,
		TargetIDs:     req.TargetIDs,
	}

	if req.StartsAt != "" {
		if t, err := time.Parse(time.RFC3339, req.StartsAt); err == nil {
			dto.StartsAt = &t
		}
	}
	if req.EndsAt != "" {
		if t, err := time.Parse(time.RFC3339, req.EndsAt); err == nil {
			dto.EndsAt = &t
		}
	}

	for _, l := range req.Links {
		dto.Links = append(dto.Links, dtos.CreateLinkDTO{
			Label:     l.Label,
			URL:       l.URL,
			SortOrder: l.SortOrder,
		})
	}

	return dto
}

func UpdateRequestToDTO(id uint, req request.UpdateAnnouncementRequest) dtos.UpdateAnnouncementDTO {
	dto := dtos.UpdateAnnouncementDTO{
		ID:            id,
		BusinessID:    req.BusinessID,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Message:       req.Message,
		DisplayType:   entities.DisplayType(req.DisplayType),
		FrequencyType: entities.FrequencyType(req.FrequencyType),
		Priority:      req.Priority,
		IsGlobal:      req.IsGlobal,
		TargetIDs:     req.TargetIDs,
	}

	if req.StartsAt != "" {
		if t, err := time.Parse(time.RFC3339, req.StartsAt); err == nil {
			dto.StartsAt = &t
		}
	}
	if req.EndsAt != "" {
		if t, err := time.Parse(time.RFC3339, req.EndsAt); err == nil {
			dto.EndsAt = &t
		}
	}

	for _, l := range req.Links {
		dto.Links = append(dto.Links, dtos.CreateLinkDTO{
			Label:     l.Label,
			URL:       l.URL,
			SortOrder: l.SortOrder,
		})
	}

	return dto
}

func EntityToResponse(e *entities.Announcement) response.AnnouncementResponse {
	resp := response.AnnouncementResponse{
		ID:             e.ID,
		BusinessID:     e.BusinessID,
		CategoryID:     e.CategoryID,
		Title:          e.Title,
		Message:        e.Message,
		DisplayType:    string(e.DisplayType),
		FrequencyType:  string(e.FrequencyType),
		Priority:       e.Priority,
		IsGlobal:       e.IsGlobal,
		Status:         string(e.Status),
		StartsAt:       e.StartsAt,
		EndsAt:         e.EndsAt,
		ForceRedisplay: e.ForceRedisplay,
		CreatedByID:    e.CreatedByID,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		Images:         make([]response.ImageResponse, 0),
		Links:          make([]response.LinkResponse, 0),
		Targets:        make([]response.TargetResponse, 0),
	}

	if e.Category != nil {
		resp.Category = &response.CategoryResponse{
			ID:    e.Category.ID,
			Code:  e.Category.Code,
			Name:  e.Category.Name,
			Icon:  e.Category.Icon,
			Color: e.Category.Color,
		}
	}

	for _, img := range e.Images {
		resp.Images = append(resp.Images, response.ImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			SortOrder: img.SortOrder,
		})
	}

	for _, l := range e.Links {
		resp.Links = append(resp.Links, response.LinkResponse{
			ID:        l.ID,
			Label:     l.Label,
			URL:       l.URL,
			SortOrder: l.SortOrder,
		})
	}

	for _, t := range e.Targets {
		resp.Targets = append(resp.Targets, response.TargetResponse{
			ID:         t.ID,
			BusinessID: t.BusinessID,
		})
	}

	return resp
}

func EntityListToResponse(items []entities.Announcement, total int64, page, pageSize int) response.PaginatedAnnouncementsResponse {
	data := make([]response.AnnouncementResponse, 0, len(items))
	for _, item := range items {
		data = append(data, EntityToResponse(&item))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return response.PaginatedAnnouncementsResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func StatsToResponse(s *entities.AnnouncementStats) response.StatsResponse {
	return response.StatsResponse{
		TotalViews:       s.TotalViews,
		UniqueUsers:      s.UniqueUsers,
		TotalClicks:      s.TotalClicks,
		TotalAcceptances: s.TotalAcceptances,
		TotalClosed:      s.TotalClosed,
	}
}

func CategoriesToResponse(cats []entities.AnnouncementCategory) []response.CategoryResponse {
	result := make([]response.CategoryResponse, len(cats))
	for i, c := range cats {
		result[i] = response.CategoryResponse{
			ID:    c.ID,
			Code:  c.Code,
			Name:  c.Name,
			Icon:  c.Icon,
			Color: c.Color,
		}
	}
	return result
}
