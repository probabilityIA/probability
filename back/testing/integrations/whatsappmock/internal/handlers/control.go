package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListSent(c *gin.Context) {
	sends := h.store.Sends()
	c.JSON(http.StatusOK, gin.H{"total": len(sends), "data": sends})
}

func (h *Handler) ListTemplates(c *gin.Context) {
	templates := h.store.Templates()
	c.JSON(http.StatusOK, gin.H{"total": len(templates), "data": templates})
}

func (h *Handler) GetScript(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"auto_reply": h.store.AutoReply(), "script": h.store.Script()})
}

type setScriptRequest struct {
	Script    map[string]string `json:"script"`
	AutoReply *bool             `json:"auto_reply"`
}

func (h *Handler) SetScript(c *gin.Context) {
	var body setScriptRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.Script) > 0 {
		h.store.SetScript(body.Script)
	}
	if body.AutoReply != nil {
		h.store.SetAutoReply(*body.AutoReply)
	}

	c.JSON(http.StatusOK, gin.H{"auto_reply": h.store.AutoReply(), "script": h.store.Script()})
}

type manualReplyRequest struct {
	Phone  string `json:"phone"`
	Button string `json:"button"`
}

func (h *Handler) ManualReply(c *gin.Context) {
	var body manualReplyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Phone == "" || body.Button == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone y button son obligatorios"})
		return
	}

	send := h.store.LastSendTo(body.Phone)
	if send == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no hay ningun mensaje enviado a ese numero"})
		return
	}

	h.sendButtonReply(send, body.Button)

	c.JSON(http.StatusOK, gin.H{
		"replied_to":   send.TemplateName,
		"context_id":   send.MessageID,
		"button":       body.Button,
		"phone_number": body.Phone,
	})
}

type approveRequest struct {
	Name   string `json:"name"`
	Event  string `json:"event"`
	Reason string `json:"reason"`
	Silent bool   `json:"silent"`
}

func (h *Handler) ApproveByName(c *gin.Context) {
	var body approveRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name es obligatorio"})
		return
	}

	template := h.store.TemplateByName(body.Name)
	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "el mock no tiene esa plantilla"})
		return
	}

	event := body.Event
	if event == "" {
		event = "APPROVED"
	}

	if body.Silent {
		template.Status = event
		c.JSON(http.StatusOK, gin.H{"name": template.Name, "event": event, "silent": true})
		return
	}

	h.emitTemplateStatus(template, event, body.Reason)

	c.JSON(http.StatusOK, gin.H{"name": template.Name, "event": event})
}

func (h *Handler) Reset(c *gin.Context) {
	h.store.Reset()
	c.JSON(http.StatusOK, gin.H{"status": "reset"})
}
