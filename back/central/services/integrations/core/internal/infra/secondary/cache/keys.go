package cache

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ttlMetadata            time.Duration = 0
	ttlCredentials         time.Duration = 0
	ttlPlatformCredentials time.Duration = 0
)

func integrationKey(id uint) string {
	return "integration:meta:" + strconv.Itoa(int(id))
}

func credentialsKey(id uint) string {
	return "integration:creds:" + strconv.Itoa(int(id))
}

func codeKey(code string) string {
	return "integration:code:" + code
}

func businessTypeIndexKey(businessID, typeID uint) string {
	return fmt.Sprintf("integration:idx:biz:%d:type:%d", businessID, typeID)
}

func storeTypeIndexKey(storeID string, typeID uint) string {
	return fmt.Sprintf("integration:idx:store:%s:type:%d", storeID, typeID)
}

func platformCredentialsKey(integrationTypeID uint) string {
	return fmt.Sprintf("integration:platform_creds:%d", integrationTypeID)
}

func configValueIndexKey(typeID uint, field, value string) string {
	return fmt.Sprintf("integration:idx:cfg:%d:%s:%s", typeID, field, value)
}
