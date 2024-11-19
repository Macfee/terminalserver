package model

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	RoleID   uint   `gorm:"not null"`
	Role     Role   `gorm:"foreignKey:RoleID"`
}

type Role struct {
	ID          uint         `gorm:"primarykey"`
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"unique;not null"`
}
