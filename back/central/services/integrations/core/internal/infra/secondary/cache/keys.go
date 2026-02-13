package cache

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// TTLs
	ttlMetadata    = 24 * time.Hour
	ttlCredentials = 1 * time.Hour
)

// integrationKey retorna la key de metadata de integración
// Formato: integration:meta:{id}
func integrationKey(id uint) string {
	return "integration:meta:" + strconv.Itoa(int(id))
}

// credentialsKey retorna la key de credentials
// Formato: integration:creds:{id}
func credentialsKey(id uint) string {
	return "integration:creds:" + strconv.Itoa(int(id))
}

// codeKey retorna la key del índice por código
// Formato: integration:code:{code}
func codeKey(code string) string {
	return "integration:code:" + code
}

// businessTypeIndexKey retorna la key del índice por business+type
// Formato: integration:idx:biz:{biz}:type:{type}
func businessTypeIndexKey(businessID, typeID uint) string {
	return fmt.Sprintf("integration:idx:biz:%d:type:%d", businessID, typeID)
}
