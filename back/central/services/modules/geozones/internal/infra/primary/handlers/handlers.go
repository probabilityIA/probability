package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/primary/handlers/response"
)

func (h *Handlers) List(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var parentID *uint
	if p := c.Query("parent_id"); p != "" {
		if id, err := strconv.ParseUint(p, 10, 64); err == nil {
			v := uint(id)
			parentID = &v
		}
	}

	includeGeom := c.Query("include_geometry") == "true"

	params := dtos.ListGeozonesParams{
		BusinessID:  businessID,
		Type:        c.Query("type"),
		ParentID:    parentID,
		Code:        c.Query("code"),
		Search:      c.Query("search"),
		IncludeGeom: includeGeom,
		Page:        page,
		PageSize:    pageSize,
	}

	items, total, err := h.uc.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.GeozoneResponse, len(items))
	for i := range items {
		data[i] = response.FromEntity(&items[i])
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.ListResponse{
		Data: data, Total: total, Page: params.Page, PageSize: params.PageSize, TotalPages: totalPages,
	})
}

func (h *Handlers) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	includeGeom := c.DefaultQuery("include_geometry", "true") == "true"

	g, err := h.uc.Get(c.Request.Context(), uint(id), includeGeom)
	if err != nil {
		if errors.Is(err, domainerrors.ErrGeozoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.FromEntity(g))
}

func (h *Handlers) Create(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	var req request.CreateGeozoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := h.uc.Create(c.Request.Context(), dtos.CreateGeozoneDTO{
		BusinessID: businessID,
		ParentID:   req.ParentID,
		Type:       req.Type,
		Code:       req.Code,
		Name:       req.Name,
		Geometry:   req.Geometry,
		Properties: req.Properties,
	})
	if err != nil {
		mapErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.FromEntity(g))
}

func (h *Handlers) Bulk(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	var req request.BulkImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	features := make([]dtos.BulkImportFeature, 0, len(req.Features))
	for _, f := range req.Features {
		features = append(features, dtos.BulkImportFeature{
			Type:       f.Properties.Type,
			Code:       f.Properties.Code,
			Name:       f.Properties.Name,
			ParentCode: f.Properties.ParentCode,
			Geometry:   f.Geometry,
		})
	}

	res, err := h.uc.BulkImport(c.Request.Context(), dtos.BulkImportDTO{
		BusinessID: businessID,
		Features:   features,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.BulkImportResponse{
		Created: res.Created, Skipped: res.Skipped, Errors: res.Errors,
	})
}

func (h *Handlers) Lookup(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat invalida"})
		return
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lng invalida"})
		return
	}

	items, err := h.uc.Lookup(c.Request.Context(), dtos.LookupParams{
		BusinessID: businessID, Lat: lat, Lng: lng, Type: c.Query("type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.GeozoneResponse, len(items))
	for i := range items {
		data[i] = response.FromEntity(&items[i])
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handlers) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		mapErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrGeozoneNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domainerrors.ErrInvalidGeometry),
		errors.Is(err, domainerrors.ErrInvalidType):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domainerrors.ErrDuplicateGeozone):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
