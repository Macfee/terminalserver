package auth

import (
	"errors"
	"sync"
)

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

const (
	ResourceTypeHost    ResourceType = "host"
	ResourceTypeSession ResourceType = "session"
	ResourceTypeUser    ResourceType = "user"
)

type Permission struct {
	ResourceType ResourceType
	Action       string
}

// 预定义权限常量
var (
	// 会话权限
	PermSessionView   = Permission{ResourceTypeSession, "view"}
	PermSessionCreate = Permission{ResourceTypeSession, "create"}
	PermSessionEdit   = Permission{ResourceTypeSession, "edit"}
	PermSessionDelete = Permission{ResourceTypeSession, "delete"}

	// 主机权限
	PermHostView   = Permission{ResourceTypeHost, "view"}
	PermHostCreate = Permission{ResourceTypeHost, "create"}
	PermHostEdit   = Permission{ResourceTypeHost, "edit"}
	PermHostDelete = Permission{ResourceTypeHost, "delete"}

	// 用户权限
	PermUserView   = Permission{ResourceTypeUser, "view"}
	PermUserCreate = Permission{ResourceTypeUser, "create"}
	PermUserEdit   = Permission{ResourceTypeUser, "edit"}
	PermUserDelete = Permission{ResourceTypeUser, "delete"}
)

type Resource struct {
	ID          string
	Type        ResourceType
	Name        string
	Description string
	Permissions []Permission
}

type Role struct {
	ID          uint
	Name        string
	Permissions []Permission
}

type ResourceManager struct {
	resources sync.Map
	roles     map[string]Role
	userRoles map[uint][]string
	acl       sync.Map
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		roles:     make(map[string]Role),
		userRoles: make(map[uint][]string),
	}
}

func (m *ResourceManager) AddRole(role Role) error {
	if _, exists := m.roles[role.Name]; exists {
		return errors.New("role already exists")
	}
	m.roles[role.Name] = role
	return nil
}

func (m *ResourceManager) AssignRole(userID uint, roleName string) error {
	if _, exists := m.roles[roleName]; !exists {
		return errors.New("role not found")
	}

	if roles, exists := m.userRoles[userID]; exists {
		m.userRoles[userID] = append(roles, roleName)
	} else {
		m.userRoles[userID] = []string{roleName}
	}
	return nil
}

func (m *ResourceManager) GetUserPermissions(userID uint) []Permission {
	var permissions []Permission
	if roles, exists := m.userRoles[userID]; exists {
		for _, roleName := range roles {
			if role, ok := m.roles[roleName]; ok {
				permissions = append(permissions, role.Permissions...)
			}
		}
	}
	return permissions
}

func (r *ResourceManager) CheckPermission(userID uint, perm Permission) bool {
	roles := r.userRoles[userID]
	for _, roleName := range roles {
		role, exists := r.roles[roleName]
		if !exists {
			continue
		}
		for _, p := range role.Permissions {
			if p == perm {
				return true
			}
		}
	}
	return false
}

func (m *ResourceManager) CreateResource(resource Resource) error {
	if _, exists := m.resources.Load(resource.ID); exists {
		return errors.New("resource already exists")
	}
	m.resources.Store(resource.ID, resource)
	return nil
}

func (m *ResourceManager) GrantAccess(resourceID string, userID uint, permissions []Permission) error {
	var aclMap map[uint][]Permission
	if existing, ok := m.acl.Load(resourceID); ok {
		aclMap = existing.(map[uint][]Permission)
	} else {
		aclMap = make(map[uint][]Permission)
	}
	aclMap[userID] = permissions
	m.acl.Store(resourceID, aclMap)
	return nil
}

func (m *ResourceManager) CheckAccess(resourceID string, userID uint, requiredPerm Permission) bool {
	if aclMap, ok := m.acl.Load(resourceID); ok {
		if perms, exists := aclMap.(map[uint][]Permission)[userID]; exists {
			for _, p := range perms {
				if p == requiredPerm {
					return true
				}
			}
		}
	}
	return false
}
