package service

import (
	"audit-system/internal/auth"
)

type resourceService struct {
	manager *auth.ResourceManager
}

var ResourceService = &resourceService{
	manager: auth.NewResourceManager(),
}

func (s *resourceService) CreateResource(resource auth.Resource) error {
	return s.manager.CreateResource(resource)
}

func (s *resourceService) GrantAccess(resourceID string, userID uint, permissions []auth.Permission) error {
	return s.manager.GrantAccess(resourceID, userID, permissions)
}
