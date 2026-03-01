package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

// SaveLog СЃРѕС…СЂР°РЅСЏРµС‚ Р·Р°РїРёСЃСЊ Р»РѕРіР° Рё РѕР±РЅРѕРІР»СЏРµС‚ РїРѕРёСЃРєРѕРІС‹Рµ РёРЅРґРµРєСЃС‹ РІ Redis
func (r *redisRepository) SaveLog(ctx context.Context, entry models.LogEntry) error {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	score := float64(entry.Timestamp.UnixNano())
	allLogsKey := r.buildKey(constants.PrefixLogs)
	minScore := float64(time.Now().Add(-constants.LogTTL).UnixNano())

	err = r.client.ZAdd(ctx, allLogsKey, redis.Z{
		Score:  score,
		Member: jsonData,
	}).Err()
	if err != nil {
		return err
	}
	r.client.ZRemRangeByScore(ctx, allLogsKey, "-inf", fmt.Sprintf("%f", minScore))
	r.client.Expire(ctx, allLogsKey, constants.LogTTL)

	if entry.UserID != "" {
		userKey := r.buildKey(fmt.Sprintf("%s%s", constants.PrefixLogUser, entry.UserID))
		r.client.ZAdd(ctx, userKey, redis.Z{Score: score, Member: jsonData})
		r.client.ZRemRangeByScore(ctx, userKey, "-inf", fmt.Sprintf("%f", minScore))
		r.client.Expire(ctx, userKey, constants.LogTTL)
	}

	if entry.Service != "" {
		serviceKey := r.buildKey(fmt.Sprintf("%s%s", constants.PrefixLogService, entry.Service))
		r.client.ZAdd(ctx, serviceKey, redis.Z{Score: score, Member: jsonData})
		r.client.ZRemRangeByScore(ctx, serviceKey, "-inf", fmt.Sprintf("%f", minScore))
		r.client.Expire(ctx, serviceKey, constants.LogTTL)
	}

	if entry.Type != "" {
		typeKey := r.buildKey(fmt.Sprintf("%s%s", constants.PrefixLogType, entry.Type))
		r.client.ZAdd(ctx, typeKey, redis.Z{Score: score, Member: jsonData})
		r.client.ZRemRangeByScore(ctx, typeKey, "-inf", fmt.Sprintf("%f", minScore))
		r.client.Expire(ctx, typeKey, constants.LogTTL)
	}

	return nil
}

// GetLogs РІРѕР·РІСЂР°С‰Р°РµС‚ СЃРїРёСЃРѕРє РІСЃРµС… Р»РѕРіРѕРІ СЃ РїР°РіРёРЅР°С†РёРµР№
func (r *redisRepository) GetLogs(ctx context.Context, page, pageSize int) ([]models.LogEntry, int64, error) {
	key := r.buildKey(constants.PrefixLogs)
	minScore := float64(time.Now().Add(-constants.LogTTL).UnixNano())
	r.client.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", minScore))
	
	return r.getPaginatedLogs(ctx, key, page, pageSize)
}

// FilterLogs С„РёР»СЊС‚СЂСѓРµС‚ Р»РѕРіРё РїРѕ РїРѕР»СЊР·РѕРІР°С‚РµР»СЋ, СЃРµСЂРІРёСЃСѓ Рё С‚РёРїСѓ
func (r *redisRepository) FilterLogs(ctx context.Context, userID string, service, logType string, page, pageSize int) ([]models.LogEntry, int64, error) {
	var keys []string
	minScore := float64(time.Now().Add(-constants.LogTTL).UnixNano())
	minScoreStr := fmt.Sprintf("%f", minScore)
	
	if userID != "" && userID != "0" {
		key := r.buildKey(fmt.Sprintf("%s%s", constants.PrefixLogUser, userID))
		r.client.ZRemRangeByScore(ctx, key, "-inf", minScoreStr)
		keys = append(keys, key)
	}
	if service != "" {
		key := r.buildKey(fmt.Sprintf("%s%s", constants.PrefixLogService, service))
		r.client.ZRemRangeByScore(ctx, key, "-inf", minScoreStr)
		keys = append(keys, key)
	}
	if logType != "" {
		key := r.buildKey(fmt.Sprintf("%s%s", constants.PrefixLogType, logType))
		r.client.ZRemRangeByScore(ctx, key, "-inf", minScoreStr)
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return r.GetLogs(ctx, page, pageSize)
	}

	if len(keys) == 1 {
		return r.getPaginatedLogs(ctx, keys[0], page, pageSize)
	}

	tempKey := r.buildKey(fmt.Sprintf("temp:logs:%d", time.Now().UnixNano()))
	err := r.client.ZInterStore(ctx, tempKey, &redis.ZStore{
		Keys: keys,
	}).Err()
	if err != nil {
		return nil, 0, err
	}
	defer r.client.Del(ctx, tempKey)

	return r.getPaginatedLogs(ctx, tempKey, page, pageSize)
}

func (r *redisRepository) getPaginatedLogs(ctx context.Context, key string, page, pageSize int) ([]models.LogEntry, int64, error) {
	total, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}

	start := int64((page - 1) * pageSize)
	stop := start + int64(pageSize) - 1

	vals, err := r.client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, 0, err
	}

	var logs []models.LogEntry
	for _, val := range vals {
		var entry models.LogEntry
		if err := json.Unmarshal([]byte(val), &entry); err == nil {
			logs = append(logs, entry)
		}
	}

	return logs, total, nil
}
