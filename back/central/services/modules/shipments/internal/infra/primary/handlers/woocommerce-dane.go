package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type daneItemResponse struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	State string `json:"state,omitempty"`
}

func (h *Handlers) WooCommerceDaneStates(c *gin.Context) {
	if _, ok := h.authWooPublic(c); !ok {
		return
	}

	states, err := h.uc.Repo().ListDaneStates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"states": []daneItemResponse{}})
		return
	}

	out := make([]daneItemResponse, 0, len(states))
	for _, s := range states {
		out = append(out, daneItemResponse{Code: s.Code, Name: s.Name})
	}
	c.JSON(http.StatusOK, gin.H{"states": out})
}

func (h *Handlers) WooCommerceDaneCities(c *gin.Context) {
	if _, ok := h.authWooPublic(c); !ok {
		return
	}

	stateCode := strings.TrimSpace(c.Query("state"))
	term := strings.TrimSpace(c.Query("q"))

	if stateCode == "" && term == "" {
		c.JSON(http.StatusOK, gin.H{"cities": []daneItemResponse{}})
		return
	}

	limit, _ := strconv.Atoi(c.Query("limit"))

	cities, err := h.uc.Repo().SearchDaneCities(c.Request.Context(), stateCode, term, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"cities": []daneItemResponse{}})
		return
	}

	out := make([]daneItemResponse, 0, len(cities))
	for _, ct := range cities {
		out = append(out, daneItemResponse{Code: ct.Code, Name: ct.Name, State: ct.StateName})
	}
	c.JSON(http.StatusOK, gin.H{"cities": out})
}
