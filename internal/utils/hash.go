package utils

import (
	"hash/fnv"

	"golang.org/x/crypto/bcrypt"
)

// HashStringToInt64 РґРµС‚РµСЂРјРёРЅРёСЂРѕРІР°РЅРЅРѕ С…РµС€РёСЂСѓРµС‚ СЃС‚СЂРѕРєСѓ РІ int64
func HashStringToInt64(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

// HashPassword С…РµС€РёСЂСѓРµС‚ РїР°СЂРѕР»СЊ СЃ РїРѕРјРѕС‰СЊСЋ bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash РїСЂРѕРІРµСЂСЏРµС‚ СЃРѕРѕС‚РІРµС‚СЃС‚РІРёРµ РїР°СЂРѕР»СЏ С…РµС€Сѓ
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
