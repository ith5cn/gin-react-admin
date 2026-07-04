package system

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"server/config"
	commonResponse "server/model/common/response"
	loggerInit "server/setup/logger"
	redisInit "server/setup/redis"
	"server/utils"

	"go.uber.org/zap"
)

// 在线用户基于 Redis 实现：登录/刷新时按 access token 的 jti 写一条会话记录，
// TTL 与 access token 一致，token 过期后记录自动消失，无需清理任务。
// 踢下线 = 删掉该会话的 access/refresh jti（ValidateToken 会立即失败）。

const onlineKeyPrefix = "online:token:"

// OnlineSession 是 Redis 中的在线会话记录。
// RefreshJTI 只用于踢下线时撤销 refresh token，列表响应里会清空，不暴露给前端。
type OnlineSession struct {
	UserID     uint      `json:"userId"`
	Username   string    `json:"username"`
	IP         string    `json:"ip"`
	OS         string    `json:"os"`
	Browser    string    `json:"browser"`
	LoginTime  time.Time `json:"loginTime"`
	AccessJTI  string    `json:"accessJti"`
	RefreshJTI string    `json:"refreshJti,omitempty"`
}

// RecordOnlineSession 记录一条在线会话，登录和刷新 token 成功后调用。
// 与日志一样尽力而为：Redis 写失败只记服务端日志，不影响登录主流程。
func RecordOnlineSession(pair *utils.TokenPair, userID uint, username, ip, userAgent string) {
	if pair == nil || pair.AccessJTI == "" {
		return
	}
	client := redisInit.Redis.Get()
	if client == nil {
		return
	}

	osName, browser := parseUserAgent(userAgent)
	session := OnlineSession{
		UserID:     userID,
		Username:   username,
		IP:         ip,
		OS:         osName,
		Browser:    browser,
		LoginTime:  time.Now(),
		AccessJTI:  pair.AccessJTI,
		RefreshJTI: pair.RefreshJTI,
	}
	payload, err := json.Marshal(session)
	if err != nil {
		loggerInit.Logger.Get().Error("record online session failed", zap.Error(err))
		return
	}

	ttl := time.Duration(config.JwtConfig().AccessExpiresMinute) * time.Minute
	if err := client.Set(context.Background(), onlineKeyPrefix+pair.AccessJTI, payload, ttl).Err(); err != nil {
		loggerInit.Logger.Get().Error("record online session failed", zap.Error(err))
	}
}

// RemoveOnlineSession 删除在线会话记录（登出时调用），尽力而为。
func RemoveOnlineSession(accessJTI string) {
	if accessJTI == "" {
		return
	}
	client := redisInit.Redis.Get()
	if client == nil {
		return
	}
	if err := client.Del(context.Background(), onlineKeyPrefix+accessJTI).Err(); err != nil {
		loggerInit.Logger.Get().Error("remove online session failed", zap.Error(err))
	}
}

// OnlineUserList 分页查询在线用户，支持按用户名/IP 过滤。
// 用 SCAN 遍历（不是 KEYS，避免大 key 空间阻塞 Redis），
// 在线会话量级小，过滤、排序、分页都在内存里做。
func OnlineUserList(query map[string]string) (*commonResponse.PageResult, error) {
	client := redisInit.Redis.Get()
	if client == nil {
		return nil, utils.ErrRedisNotInitialized
	}
	ctx := context.Background()

	sessions := make([]OnlineSession, 0)
	iter := client.Scan(ctx, 0, onlineKeyPrefix+"*", 100).Iterator()
	for iter.Next(ctx) {
		payload, err := client.Get(ctx, iter.Val()).Result()
		if err != nil {
			// SCAN 和 GET 之间 key 可能刚好过期，跳过即可。
			continue
		}
		var session OnlineSession
		if err := json.Unmarshal([]byte(payload), &session); err != nil {
			continue
		}
		session.RefreshJTI = ""
		sessions = append(sessions, session)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	username := query["username"]
	ip := query["ip"]
	filtered := sessions[:0]
	for _, session := range sessions {
		if username != "" && !strings.Contains(session.Username, username) {
			continue
		}
		if ip != "" && !strings.Contains(session.IP, ip) {
			continue
		}
		filtered = append(filtered, session)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].LoginTime.After(filtered[j].LoginTime)
	})

	page := parsePage(query)
	total := int64(len(filtered))
	start := (page.Page - 1) * page.Size
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + page.Size
	if end > len(filtered) {
		end = len(filtered)
	}
	return &commonResponse.PageResult{List: filtered[start:end], Total: total}, nil
}

// KickOnlineUser 把指定会话踢下线：撤销其 access/refresh token 并删除会话记录。
// 被踢的 token 下一次请求就会 401（ValidateToken 查不到 jti）。
func KickOnlineUser(accessJTI string) error {
	if accessJTI == "" {
		return ErrOnlineUserNotFound
	}
	client := redisInit.Redis.Get()
	if client == nil {
		return utils.ErrRedisNotInitialized
	}
	ctx := context.Background()

	payload, err := client.Get(ctx, onlineKeyPrefix+accessJTI).Result()
	if err != nil {
		return ErrOnlineUserNotFound
	}
	var session OnlineSession
	if err := json.Unmarshal([]byte(payload), &session); err != nil {
		return ErrOnlineUserNotFound
	}

	keys := []string{
		utils.TokenStateKey(utils.TokenTypeAccess, session.AccessJTI),
		onlineKeyPrefix + accessJTI,
	}
	if session.RefreshJTI != "" {
		keys = append(keys, utils.TokenStateKey(utils.TokenTypeRefresh, session.RefreshJTI))
	}
	return client.Del(ctx, keys...).Err()
}
