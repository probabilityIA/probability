package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func resolveBusinessID(c *gin.Context) (uint, error) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, errors.New("no se pudo identificar la empresa")
	}
	if !middleware.IsSuperAdmin(c) {
		return businessID, nil
	}
	param := c.Query("business_id")
	if param == "" {
		return 0, errors.New("super admin: business_id es requerido como query param")
	}
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("super admin: business_id invalido")
	}
	return uint(id), nil
}

func isAdminUser(c *gin.Context) bool {
	if middleware.IsSuperAdmin(c) {
		return true
	}
	roles, ok := middleware.GetUserRoles(c)
	if !ok {
		return false
	}
	for _, r := range roles {
		if strings.Contains(strings.ToLower(r), "admin") {
			return true
		}
	}
	return false
}

func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	loc := time.UTC
	now := time.Now().In(loc)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	sd := strings.TrimSpace(c.Query("start_date"))
	ed := strings.TrimSpace(c.Query("end_date"))
	if sd != "" && ed != "" {
		s, err1 := time.Parse("2006-01-02", sd)
		e, err2 := time.Parse("2006-01-02", ed)
		if err1 == nil && err2 == nil {
			e = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, loc)
			return s, e
		}
	}

	switch strings.TrimSpace(c.Query("range")) {
	case "today":
		return startOfToday, endOfToday
	case "week":
		return startOfToday.AddDate(0, 0, -6), endOfToday
	case "3months":
		return startOfToday.AddDate(0, 0, -89), endOfToday
	case "all":
		return time.Time{}, time.Time{}
	default:
		return startOfToday.AddDate(0, 0, -29), endOfToday
	}
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
