package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString РіРµРЅРµСЂРёСЂСѓРµС‚ СЃР»СѓС‡Р°Р№РЅСѓСЋ hex-СЃС‚СЂРѕРєСѓ СѓРєР°Р·Р°РЅРЅРѕР№ РґР»РёРЅС‹
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
