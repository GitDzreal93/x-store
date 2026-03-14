package service

import (
	"fmt"
	"strconv"

	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/pkg/crypto"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repo.UserRepo
	orderRepo *repo.OrderRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:  repo.NewUserRepo(),
		orderRepo: repo.NewOrderRepo(),
	}
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResp struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterReq) (*AuthResp, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败")
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: nickname,
		Role:     "buyer",
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败")
	}

	// 生成 JWT Token
	token, err := crypto.GenerateToken(config.Global.JWT.Secret, config.Global.JWT.ExpireHours, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败")
	}

	user.Password = "" // 不返回密码
	return &AuthResp{Token: token, User: user}, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginReq) (*AuthResp, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, fmt.Errorf("账号已被禁用")
	}

	// 生成 JWT Token
	token, err := crypto.GenerateToken(config.Global.JWT.Secret, config.Global.JWT.ExpireHours, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败")
	}

	user.Password = "" // 不返回密码
	return &AuthResp{Token: token, User: user}, nil
}

// GetByID 根据 ID 获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.Password = "" // 不返回密码
	return user, nil
}

// GetUserOrders 获取用户订单列表
func (s *UserService) GetUserOrders(userID uint, page, size string) ([]model.Order, int64, error) {
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	sz, _ := strconv.Atoi(size)
	if sz < 1 || sz > 100 {
		sz = 20
	}
	return s.orderRepo.ListByUser(userID, p, sz)
}
