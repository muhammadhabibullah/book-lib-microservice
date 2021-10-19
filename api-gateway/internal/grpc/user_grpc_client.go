package grpc

import (
	"context"
	"errors"
	"log"

	"api-gateway/internal/domain/constant"
	"api-gateway/internal/graph/model"
	"api-gateway/pkg/proto"
)

type UserGRPCService struct {
	client proto.UserServiceClient
}

func NewUserGRPCService(
	client proto.UserServiceClient,
) *UserGRPCService {
	return &UserGRPCService{
		client: client,
	}
}

func (c *UserGRPCService) CreateUser(ctx context.Context, input model.NewUser, role model.Role) (*model.User, error) {
	user, err := c.client.CreateUser(ctx, &proto.CreateUserRequest{
		Email:    input.Email,
		Password: input.Password,
		Role:     string(role),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
		Role:  model.Role(user.GetRole()),
	}, nil
}

func (c *UserGRPCService) Login(ctx context.Context, input model.Login) (string, error) {
	loginResponse, err := c.client.Login(ctx, &proto.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	return loginResponse.Token, nil
}

func (c *UserGRPCService) FetchUser(ctx context.Context, filter model.FetchUserFilter) (*model.UserPaged, error) {
	var (
		limit int32 = 10
		page  int32 = 1
		email string
		role  string
	)
	if filter.Limit != nil {
		limit = int32(*filter.Limit)
	}
	if filter.Page != nil {
		page = int32(*filter.Page)
	}
	if filter.Email != nil {
		email = *filter.Email
	}
	if filter.Role != nil {
		role = string(*filter.Role)
	}

	fetchUserResponse, err := c.client.FetchUser(ctx, &proto.FetchUserRequest{
		Pagination: &proto.PaginationRequest{
			Limit: limit,
			Page:  page,
		},
		Email: email,
		Role:  role,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	users := make([]*model.User, 0)
	for _, fetchedUser := range fetchUserResponse.Users {
		users = append(users, &model.User{
			ID:    fetchedUser.GetId(),
			Email: fetchedUser.GetEmail(),
			Role:  model.Role(fetchedUser.GetRole()),
		})
	}

	return &model.UserPaged{
		Users:     users,
		Page:      int(fetchUserResponse.Pagination.GetPage()),
		Limit:     int(fetchUserResponse.Pagination.GetLimit()),
		TotalUser: int(fetchUserResponse.Pagination.GetTotal()),
		LastPage:  int(fetchUserResponse.Pagination.GetLastPage()),
	}, nil
}

func (c *UserGRPCService) FindByID(ctx context.Context, id string) (*model.User, error) {
	user, err := c.client.FindByID(ctx, &proto.FindByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
		Role:  model.Role(user.GetRole()),
	}, nil
}

func (c *UserGRPCService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := c.client.FindByEmail(ctx, &proto.FindByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
		Role:  model.Role(user.GetRole()),
	}, nil
}

func (c *UserGRPCService) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	user, err := c.client.UpdateUser(ctx, &proto.UpdateUserRequest{
		Id:    input.ID,
		Email: input.Email,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
		Role:  model.Role(user.GetRole()),
	}, nil
}

func (c *UserGRPCService) UpdateSelf(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	selfEmail, exist := ctx.Value(constant.EmailGinCtxKey).(string)
	if !exist {
		return nil, errors.New("missing email on authorization token")
	}

	user, err := c.client.UpdateSelf(ctx, &proto.UpdateSelfRequest{
		SelfEmail: selfEmail,
		Id:        input.ID,
		Email:     input.Email,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
		Role:  model.Role(user.GetRole()),
	}, nil
}

func (c *UserGRPCService) DeleteUser(ctx context.Context, input model.DeleteUser) error {
	_, err := c.client.DeleteUser(ctx, &proto.DeleteUserRequest{
		Email: input.Email,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
