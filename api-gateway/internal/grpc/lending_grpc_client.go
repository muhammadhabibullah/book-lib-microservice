package grpc

import (
	"context"
	"errors"
	"io"
	"log"

	"api-gateway/internal/domain/constant"
	"api-gateway/internal/graph/model"
	"api-gateway/pkg/grpc"
	"api-gateway/pkg/proto"
)

type LendingGRPCService struct {
	client proto.LendingServiceClient
}

func NewLendingGRPCService(
	client proto.LendingServiceClient,
) *LendingGRPCService {
	return &LendingGRPCService{
		client: client,
	}
}

func (c *LendingGRPCService) LendBook(ctx context.Context, input model.NewLending) (*model.Lending, error) {
	selfUserID, exist := ctx.Value(constant.UserIDGinCtxKey).(string)
	if !exist {
		return nil, errors.New("missing userID on authorization token")
	}

	stream, err := c.client.CreateLending(ctx, &proto.CreateLendingRequest{
		BookId: input.BookID,
		UserId: selfUserID,
	})
	if err != nil {
		log.Println(err)
		if err := grpc.ParseErrorStatus(err); err != nil {
			return nil, err
		}

		return nil, err
	}

	var lending *model.Lending

	for {
		lendingStream, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			if err := grpc.ParseErrorStatus(err); err != nil {
				return nil, err
			}

			return nil, err
		}

		lending = &model.Lending{
			ID:         lendingStream.GetId(),
			BookID:     lendingStream.GetBookId(),
			UserID:     lendingStream.GetUserId(),
			Status:     lendingStream.GetStatus(),
			ReturnDate: lendingStream.GetReturnDate().String(),
		}

		switch constant.LendingStatus(lending.Status) {
		case constant.LendingActive:
			return lending, nil
		case constant.LendingCanceled:
			return lending, errors.New("failed to lend book")
		default:
			log.Printf("received lending stream: %+v", lending)
		}
	}

	return lending, nil
}

func (c *LendingGRPCService) RenewLending(ctx context.Context, input model.RenewLendingRequest) (*model.Lending, error) {
	lending, err := c.client.RenewLending(ctx, &proto.RenewLendingRequest{
		Id: input.ID,
	})
	if err != nil {
		log.Println(err)
		if err := grpc.ParseErrorStatus(err); err != nil {
			return nil, err
		}

		return nil, err
	}

	return &model.Lending{
		ID:         lending.GetId(),
		BookID:     lending.GetBookId(),
		UserID:     lending.GetUserId(),
		Status:     lending.GetStatus(),
		ReturnDate: lending.GetReturnDate().String(),
	}, nil
}

func (c *LendingGRPCService) FinishLending(ctx context.Context, input model.FinishLendingRequest) (*model.Lending, error) {
	lending, err := c.client.FinishLending(ctx, &proto.FinishLendingRequest{
		Id: input.ID,
	})
	if err != nil {
		log.Println(err)
		if err := grpc.ParseErrorStatus(err); err != nil {
			return nil, err
		}

		return nil, err
	}

	return &model.Lending{
		ID:         lending.GetId(),
		BookID:     lending.GetBookId(),
		UserID:     lending.GetUserId(),
		Status:     lending.GetStatus(),
		ReturnDate: lending.GetReturnDate().String(),
	}, nil
}

func (c *LendingGRPCService) MyLending(ctx context.Context, input *model.MyLendingRequest) (*model.LendingPaged, error) {
	selfUserID, exist := ctx.Value(constant.UserIDGinCtxKey).(string)
	if !exist {
		return nil, errors.New("missing userID on authorization token")
	}

	fetchInput := model.FetchLendingRequest{
		Page:   input.Page,
		Limit:  input.Limit,
		UserID: &selfUserID,
		Status: input.Status,
	}

	return c.FetchLending(ctx, &fetchInput)
}

func (c *LendingGRPCService) FetchLending(ctx context.Context, input *model.FetchLendingRequest) (*model.LendingPaged, error) {
	var (
		limit  int32 = 10
		page   int32 = 1
		userID string
		status string
	)
	if input.Limit != nil {
		limit = int32(*input.Limit)
	}
	if input.Page != nil {
		page = int32(*input.Page)
	}
	if input.UserID != nil {
		userID = *input.UserID
	}
	if input.Status != nil {
		status = *input.Status
	}

	fetchedLending, err := c.client.FetchLending(ctx, &proto.FetchLendingRequest{
		Pagination: &proto.LendingPaginationRequest{
			Limit: limit,
			Page:  page,
		},
		UserId: userID,
		Status: status,
	})
	if err != nil {
		log.Println(err)
		if err := grpc.ParseErrorStatus(err); err != nil {
			return nil, err
		}

		return nil, err
	}

	lendings := make([]*model.Lending, 0)
	for _, lending := range fetchedLending.Lendings {
		lendings = append(lendings, &model.Lending{
			ID:         lending.GetId(),
			BookID:     lending.GetBookId(),
			UserID:     lending.GetUserId(),
			Status:     lending.GetStatus(),
			ReturnDate: lending.GetReturnDate().String(),
		})
	}

	return &model.LendingPaged{
		Lendings:     lendings,
		Page:         int(fetchedLending.GetPagination().GetPage()),
		Limit:        int(fetchedLending.GetPagination().GetLimit()),
		TotalLending: int(fetchedLending.GetPagination().GetTotal()),
		LastPage:     int(fetchedLending.GetPagination().GetLastPage()),
	}, nil
}
