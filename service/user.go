package service

import (
	"deouy/wechatbot/config"
	"deouy/wechatbot/gpt"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

// UserServiceInterface 用户业务接口
type UserServiceInterface interface {
	GetUserSessionContext(userId string) []gpt.Message
	SetUserSessionContext(userId string, reply string)
	ClearUserSessionContext(userId string, msg string) bool
}

var _ UserServiceInterface = (*UserService)(nil)

// UserService 用戶业务
type UserService struct {
	// 缓存
	cache *cache.Cache
}

// NewUserService 创建新的业务层
func NewUserService() UserServiceInterface {
	return &UserService{cache: cache.New(time.Second*config.LoadConfig().SessionTimeout, time.Minute*10)}
}

// ClearUserSessionContext 清空GTP上下文，接收文本中包含`我要问下一个问题`，并且Unicode 字符数量不超过20就清空
func (s *UserService) ClearUserSessionContext(userId string, msg string) bool {
	messages := s.GetUserSessionContext(userId)

	if strings.Contains(msg, "我要问下一个问题") && len(messages) < 20 {
		s.cache.Delete(userId)
		return true
	}
	return false
}

// GetUserSessionContext 获取用户会话上下文文本
func (s *UserService) GetUserSessionContext(userId string) []gpt.Message {
	sessionContext, ok := s.cache.Get(userId)
	if !ok {
		return []gpt.Message{}
	}
	return sessionContext.([]gpt.Message)
}

// SetUserSessionContext 设置用户会话上下文文本，question用户提问内容，GTP回复内容
func (s *UserService) SetUserSessionContext(userId string, reply string) {
	messages := s.GetUserSessionContext(userId)
	message := gpt.CreateMessage(reply)
	messages = append(messages, message)
	s.cache.Set(userId, messages, time.Second*config.LoadConfig().SessionTimeout)
}
