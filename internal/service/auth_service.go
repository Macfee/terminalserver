package service

import (
	"audit-system/internal/auth"
	"audit-system/internal/model"
	"audit-system/pkg/database"
	"audit-system/pkg/logger"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	permissionCache map[uint][]string // 用户权限缓存
	cacheMutex      sync.RWMutex
}

var AuthService = &authService{
	permissionCache: make(map[uint][]string),
}

func (s *authService) ValidateUser(username, password string) (*model.User, error) {
	var user model.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

func (s *authService) CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return database.DB.Create(user).Error
}

func (s *authService) CanManageAlerts(userID uint) bool {
	return s.hasPermission(userID, "alert:manage")
}

func (s *authService) CanMonitorSession(userID uint, sessionID string) bool {
	// 检查用户是否有监控会话的权限
	return s.hasPermission(userID, "session:monitor")
}

func (s *authService) CanReplaySession(userID uint, sessionID string) bool {
	// 检查用户是否有回放会话的权限
	return s.hasPermission(userID, "session:replay")
}

func (s *authService) CanManageResource(userID uint, resourceType auth.ResourceType) bool {
	// 将ResourceType转换为字符串
	return s.hasPermission(userID, fmt.Sprintf("%s:manage", resourceType.String()))
}

func (s *authService) CanGrantAccess(userID uint, resourceID string) bool {
	// 检查用户是否有授权权限
	return s.hasPermission(userID, "resource:grant")
}

func (s *authService) hasPermission(userID uint, permission string) bool {
	// 先检查缓存
	s.cacheMutex.RLock()
	permissions, exists := s.permissionCache[userID]
	s.cacheMutex.RUnlock()

	if !exists {
		// 从数据库加载权限
		permissions = s.loadUserPermissions(userID)
		// 更新缓存
		s.cacheMutex.Lock()
		s.permissionCache[userID] = permissions
		s.cacheMutex.Unlock()
	}

	// 检查是否有超级管理员权限
	for _, p := range permissions {
		if p == "*" || p == "admin" {
			return true
		}
	}

	// 检查具体权限
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

func (s *authService) loadUserPermissions(userID uint) []string {
	var user model.User
	if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		logger.GetLogger().Error("加载用户权限失败", "error", err, "userID", userID)
		return nil
	}

	var permissions []string
	for _, perm := range user.Role.Permissions {
		permissions = append(permissions, perm.Name)
	}

	return permissions
}

func (s *authService) ClearPermissionCache(userID uint) {
	s.cacheMutex.Lock()
	delete(s.permissionCache, userID)
	s.cacheMutex.Unlock()
}
