package service

import (
	"audit-system/pkg/database"
	"errors"

	"gorm.io/gorm"
)

type Host struct {
	ID         string `json:"id" gorm:"primarykey"`
	Name       string `json:"name" gorm:"not null"`
	IP         string `json:"ip" gorm:"not null;unique"`
	Port       int    `json:"port" gorm:"not null"`
	Username   string `json:"username" gorm:"not null"`
	Password   string `json:"password" gorm:"not null"`
	Status     string `json:"status" gorm:"default:'offline'"`
	LastOnline int64  `json:"last_online"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime int64  `json:"update_time" gorm:"autoUpdateTime"`
}

type hostService struct{}

var HostService = new(hostService)

func (s *hostService) ListHosts(page, pageSize int) ([]Host, int64, error) {
	var hosts []Host
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := database.DB.Model(&Host{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := database.DB.Offset(offset).Limit(pageSize).Find(&hosts).Error; err != nil {
		return nil, 0, err
	}

	return hosts, total, nil
}

func (s *hostService) CreateHost(host *Host) error {
	if host.IP == "" || host.Username == "" {
		return errors.New("IP和用户名不能为空")
	}

	// 检查IP是否已存在
	var count int64
	if err := database.DB.Model(&Host{}).Where("ip = ?", host.IP).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("IP已存在")
	}

	return database.DB.Create(host).Error
}

func (s *hostService) GetHost(hostID string) (*Host, error) {
	var host Host
	if err := database.DB.First(&host, "id = ?", hostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("主机不存在")
		}
		return nil, err
	}
	return &host, nil
}

func (s *hostService) UpdateHost(hostID string, host *Host) error {
	// 检查主机是否存在
	var existingHost Host
	if err := database.DB.First(&existingHost, "id = ?", hostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("主机不存在")
		}
		return err
	}

	// 如果IP变更，检查新IP是否已存在
	if host.IP != existingHost.IP {
		var count int64
		if err := database.DB.Model(&Host{}).Where("ip = ? AND id != ?", host.IP, hostID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("IP已存在")
		}
	}

	return database.DB.Model(&existingHost).Updates(host).Error
}

func (s *hostService) DeleteHost(hostID string) error {
	result := database.DB.Delete(&Host{}, "id = ?", hostID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("主机不存在")
	}
	return nil
}
