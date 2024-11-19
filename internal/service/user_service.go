package service

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

type userService struct{}

var UserService = new(userService)

func (s *userService) GetUserProfile(userID uint) (*User, error) {
	// TODO: 从数据库获取用户信息
	return &User{
		ID:       userID,
		Username: "test",
		Role:     "admin",
		Email:    "test@example.com",
	}, nil
}
