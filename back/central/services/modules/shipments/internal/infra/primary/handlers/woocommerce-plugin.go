package handlers

import (
	"archive/zip"
	"bytes"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/primary/plugin"
)

func (h *Handlers) WooCommercePluginDownload(c *gin.Context) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	err := fs.WalkDir(plugin.Files, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		data, readErr := plugin.Files.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		w, createErr := zw.Create(path)
		if createErr != nil {
			return createErr
		}
		_, writeErr := w.Write(data)
		return writeErr
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el plugin"})
		return
	}
	if err := zw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el plugin"})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="probability-shipping.zip"`)
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

func (h *Handlers) authorizeWooIntegration(c *gin.Context) (uint, bool) {
	idStr := c.Param("integration_id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id invalido"})
		return 0, false
	}
	integrationID := uint(id64)

	bizOfIntegration, err := h.uc.Repo().GetIntegrationBusinessID(c.Request.Context(), integrationID)
	if err != nil || bizOfIntegration == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "integracion no encontrada"})
		return 0, false
	}

	callerBiz, _ := middleware.GetBusinessID(c)
	if !middleware.IsSuperAdmin(c) && callerBiz != bizOfIntegration {
		c.JSON(http.StatusForbidden, gin.H{"error": "no autorizado para esta integracion"})
		return 0, false
	}

	return integrationID, true
}

func (h *Handlers) wooConnectionResponse(integrationID uint, salt string, revoked bool) gin.H {
	resp := gin.H{
		"integration_id":      integrationID,
		"backend_url":         h.pluginBaseURL,
		"plugin_download_url": strings.TrimRight(h.pluginBaseURL, "/") + "/api/v1/woocommerce/plugin-download",
		"plugin_version":      plugin.Version,
		"revoked":             revoked,
	}
	if revoked {
		resp["token"] = ""
		resp["connection_key"] = ""
		return resp
	}
	token := deriveWooToken(h.tokenSecret, integrationID, salt)
	resp["token"] = token
	resp["connection_key"] = buildWooConnectionKey(h.pluginBaseURL, integrationID, token)
	return resp
}

func (h *Handlers) WooCommerceConnectionInfo(c *gin.Context) {
	integrationID, ok := h.authorizeWooIntegration(c)
	if !ok {
		return
	}

	salt, revoked, err := h.uc.Repo().EnsureWooShippingToken(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo obtener la clave de conexion"})
		return
	}
	h.bustWooCache(c.Request.Context(), integrationID)

	c.JSON(http.StatusOK, h.wooConnectionResponse(integrationID, salt, revoked))
}

func (h *Handlers) WooCommerceRotateToken(c *gin.Context) {
	integrationID, ok := h.authorizeWooIntegration(c)
	if !ok {
		return
	}

	salt, err := h.uc.Repo().RotateWooShippingToken(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo rotar la clave"})
		return
	}
	h.bustWooCache(c.Request.Context(), integrationID)

	c.JSON(http.StatusOK, h.wooConnectionResponse(integrationID, salt, false))
}

func (h *Handlers) WooCommerceRevokeToken(c *gin.Context) {
	integrationID, ok := h.authorizeWooIntegration(c)
	if !ok {
		return
	}

	if err := h.uc.Repo().RevokeWooShippingToken(c.Request.Context(), integrationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo revocar la clave"})
		return
	}
	h.bustWooCache(c.Request.Context(), integrationID)

	c.JSON(http.StatusOK, h.wooConnectionResponse(integrationID, "", true))
}
