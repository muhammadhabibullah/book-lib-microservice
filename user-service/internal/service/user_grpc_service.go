package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"user-service/internal/domain"
	"user-service/internal/domain/constant"
	"user-service/internal/repository"
	"user-service/pkg/jwt"
	"user-service/pkg/password"
	"user-service/pkg/proto"
)

type UserGRPCService struct {
	proto.UnimplementedUserServiceServer
	userRepository domain.UserRepository
	jwtService     jwt.Service
}

func NewUserGRPCService() *UserGRPCService {
	return &UserGRPCService{
		userRepository: repository.NewUserMongoDBRepository(),
		jwtService:     jwt.New(),
	}
}

func (s *UserGRPCService) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (
	*proto.User, error) {
	hashedPassword, _ := password.Hash(request.Password)
	user := domain.User{
		Email:          request.Email,
		HashedPassword: hashedPassword,
		Role:           request.Role,
	}

	if err := s.userRepository.Create(ctx, &user); err != nil {
		if strings.Contains(err.Error(), constant.UserEmailUniqueIndex) {
			return nil, fmt.Errorf("email %s already registered", user.Email)
		}
		return nil, err
	}

	return &proto.User{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserGRPCService) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("email not found")
		}
		return nil, err
	}

	if user.Meta.DeletedAt != nil {
		return nil, errors.New("your account is deleted")
	}
	if !password.Valid(request.Password, user.HashedPassword) {
		return nil, errors.New("wrong password")
	}

	token, err := s.jwtService.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{
		Token: token,
	}, nil
}

func (s *UserGRPCService) FetchUser(ctx context.Context, request *proto.FetchUserRequest) (*proto.FetchUserResponse, error) {
	page, limit := request.Pagination.Page, request.Pagination.Limit
	fetchFilter := map[string]interface{}{
		"page":  page,
		"limit": limit,
	}
	if email := request.Email; email != "" {
		fetchFilter["email"] = email
	}
	if role := request.Role; role != "" {
		fetchFilter["role"] = role
	}

	users, err := s.userRepository.Fetch(ctx, fetchFilter)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*proto.User, 0)
	for _, user := range users {
		protoUsers = append(protoUsers, &proto.User{
			Id:    user.ID.Hex(),
			Email: user.Email,
			Role:  user.Role,
		})
	}

	total, err := s.userRepository.Count(ctx, fetchFilter)
	if err != nil {
		return nil, err
	}

	return &proto.FetchUserResponse{
		Pagination: &proto.PaginationResponse{
			Limit:    limit,
			Page:     page,
			LastPage: int32(math.Ceil(float64(total) / float64(limit))),
			Total:    int32(total),
		},
		Users: protoUsers,
	}, nil
}

func (s *UserGRPCService) FindByID(ctx context.Context, request *proto.FindByIDRequest) (*proto.User, error) {
	user, err := s.userRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserGRPCService) FindByEmail(ctx context.Context, request *proto.FindByEmailRequest) (*proto.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserGRPCService) UpdateUser(ctx context.Context, request *proto.UpdateUserRequest) (*proto.User, error) {
	user, err := s.userRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	user.Email = request.Email

	err = s.userRepository.Update(ctx, &user)
	if err != nil {
		if strings.Contains(err.Error(), constant.UserEmailUniqueIndex) {
			return nil, fmt.Errorf("email %s already registered", user.Email)
		}
		return nil, err
	}

	return &proto.User{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserGRPCService) UpdateSelf(ctx context.Context, request *proto.UpdateSelfRequest) (*proto.User, error) {
	user, err := s.userRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if request.SelfEmail != user.Email {
		return nil, errors.New("unauthorized to change other user data")
	}

	user.Email = request.Email

	err = s.userRepository.Update(ctx, &user)
	if err != nil {
		if strings.Contains(err.Error(), constant.UserEmailUniqueIndex) {
			return nil, fmt.Errorf("email %s already registered", user.Email)
		}
		return nil, err
	}

	return &proto.User{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserGRPCService) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Delete(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteUserResponse{}, nil
}
