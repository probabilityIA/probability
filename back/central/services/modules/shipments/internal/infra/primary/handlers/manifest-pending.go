package handlers

import (
	"archive/zip"
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecasemanifest"
)

func (h *Handlers) ListManifestPending(c *gin.Context) {
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		businessIDStr = strconv.Itoa(int(c.GetUint("business_id")))
	}
	bid, err := strconv.ParseUint(businessIDStr, 10, 64)
	if err != nil || bid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "business_id requerido"})
		return
	}

	includeChildren := c.DefaultQuery("include_children", "true") == "true"
	carrier := strings.TrimSpace(c.Query("carrier"))

	groups, err := h.uc.Manifest.ListPending(c.Request.Context(), usecasemanifest.PendingFilter{
		BusinessID:      uint(bid),
		IncludeChildren: includeChildren,
		Carrier:         carrier,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	total := 0
	for _, g := range groups {
		total += g.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    groups,
		"total":   total,
	})
}

type generateManifestRequest struct {
	BusinessID  uint   `json:"business_id"`
	ShipmentIDs []uint `json:"shipment_ids"`
	Carrier     string `json:"carrier"`
}

func (h *Handlers) GenerateManifestPDF(c *gin.Context) {
	var req generateManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if req.BusinessID == 0 {
		if bid := c.GetUint("business_id"); bid > 0 {
			req.BusinessID = bid
		}
	}
	if req.BusinessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "business_id requerido"})
		return
	}
	if len(req.ShipmentIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "shipment_ids requerido"})
		return
	}

	userName, _ := c.Get("user_email")
	uname, _ := userName.(string)

	results, err := h.uc.Manifest.GeneratePDF(c.Request.Context(), usecasemanifest.GeneratePDFInput{
		BusinessID:  req.BusinessID,
		ShipmentIDs: req.ShipmentIDs,
		Carrier:     req.Carrier,
		UserName:    uname,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(results) == 1 {
		r := results[0]
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", `attachment; filename="`+r.Filename+`"`)
		c.Data(http.StatusOK, "application/pdf", r.PDF)
		return
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, r := range results {
		w, err := zw.Create(r.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		if _, err := w.Write(r.PDF); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}
	if err := zw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", `attachment; filename="manifiestos.zip"`)
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}
