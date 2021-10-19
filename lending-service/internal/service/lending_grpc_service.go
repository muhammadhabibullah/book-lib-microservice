package service

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"lending-service/internal/domain"
	"lending-service/internal/domain/constant"
	"lending-service/internal/repository"
	"lending-service/pkg/proto"
)

const (
	defaultLendingDuration = 14 * 24 // 2 week
)

type LendingGRPCService struct {
	proto.UnimplementedLendingServiceServer
	lendingRepository domain.LendingRepository
	userServiceClient proto.UserServiceClient
	bookServiceClient proto.BookServiceClient
}

func NewLendingGRPCService(
	userServiceClient proto.UserServiceClient,
	bookServiceClient proto.BookServiceClient,
) *LendingGRPCService {
	return &LendingGRPCService{
		lendingRepository: repository.NewLendingMongoDBRepository(),
		userServiceClient: userServiceClient,
		bookServiceClient: bookServiceClient,
	}
}

func (s *LendingGRPCService) CreateLending(ctx context.Context, request *proto.CreateLendingRequest) (*proto.Lending, error) {
	userID, err := primitive.ObjectIDFromHex(request.UserId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	book, err := s.bookServiceClient.FindByID(ctx, &proto.FindBookByIDRequest{
		Id: request.BookId,
	})
	if err != nil {
		return nil, err
	}

	if book.Stock == 0 {
		return nil, errors.New("book stock is empty")
	}

	book, err = s.bookServiceClient.UpdateBookStock(ctx, &proto.UpdateBookStockRequest{
		Id:          book.Id,
		StockChange: -1,
	})
	if err != nil {
		return nil, err
	}

	bookID, _ := primitive.ObjectIDFromHex(book.Id)
	lending := domain.Lending{
		BookID:     bookID,
		UserID:     userID,
		Status:     constant.LendingActive,
		ReturnDate: time.Now().Add(defaultLendingDuration * time.Hour),
	}

	err = s.lendingRepository.Create(ctx, &lending)
	if err != nil {
		go func(id string) {
			_, err := s.bookServiceClient.UpdateBookStock(ctx, &proto.UpdateBookStockRequest{
				Id:          id,
				StockChange: 1,
			})
			if err != nil {
				log.Printf("error update book stock after failed create lending: %s\n", err)
			}
		}(book.Id)

		return nil, err
	}

	return &proto.Lending{
		Id:         lending.ID.Hex(),
		BookId:     lending.BookID.Hex(),
		UserId:     lending.UserID.Hex(),
		Status:     string(lending.Status),
		ReturnDate: timestamppb.New(lending.ReturnDate),
	}, nil
}

func (s *LendingGRPCService) FetchLending(ctx context.Context, request *proto.FetchLendingRequest) (*proto.FetchLendingResponse, error) {
	page, limit := request.Pagination.Page, request.Pagination.Limit
	fetchFilter := map[string]interface{}{
		"page":  page,
		"limit": limit,
	}
	if status := request.Status; status != "" {
		fetchFilter["status"] = status
	}
	if userID := request.UserId; userID != "" {
		fetchFilter["user_id"] = userID
	}
	if bookID := request.BookId; bookID != "" {
		fetchFilter["book_id"] = bookID
	}

	lendings, err := s.lendingRepository.Fetch(ctx, fetchFilter)
	if err != nil {
		return nil, err
	}

	protoLendings := make([]*proto.Lending, 0)
	for _, lending := range lendings {
		protoLendings = append(protoLendings, &proto.Lending{
			Id:         lending.ID.Hex(),
			BookId:     lending.BookID.Hex(),
			UserId:     lending.UserID.Hex(),
			Status:     string(lending.Status),
			ReturnDate: timestamppb.New(lending.ReturnDate),
		})
	}

	total, err := s.lendingRepository.Count(ctx, fetchFilter)
	if err != nil {
		return nil, err
	}

	return &proto.FetchLendingResponse{
		Pagination: &proto.LendingPaginationResponse{
			Limit:    limit,
			Page:     page,
			LastPage: int32(math.Ceil(float64(total) / float64(limit))),
			Total:    int32(total),
		},
		Lendings: protoLendings,
	}, nil
}

func (s *LendingGRPCService) RenewLending(ctx context.Context, request *proto.RenewLendingRequest) (*proto.Lending, error) {
	lending, err := s.lendingRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	lending.ReturnDate = time.Now().Add(defaultLendingDuration * time.Hour)

	err = s.lendingRepository.Update(ctx, &lending)
	if err != nil {
		return nil, err
	}

	return &proto.Lending{
		Id:         lending.ID.Hex(),
		BookId:     lending.BookID.Hex(),
		UserId:     lending.UserID.Hex(),
		Status:     string(lending.Status),
		ReturnDate: timestamppb.New(lending.ReturnDate),
	}, nil
}

func (s *LendingGRPCService) FinishLending(ctx context.Context, request *proto.FinishLendingRequest) (*proto.Lending, error) {
	lending, err := s.lendingRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	lending.Status = constant.LendingInactive

	err = s.lendingRepository.Update(ctx, &lending)
	if err != nil {
		return nil, err
	}

	_, err = s.bookServiceClient.UpdateBookStock(ctx, &proto.UpdateBookStockRequest{
		Id:          lending.BookID.Hex(),
		StockChange: 1,
	})
	if err != nil {
		return nil, err
	}

	return &proto.Lending{
		Id:         lending.ID.Hex(),
		BookId:     lending.BookID.Hex(),
		UserId:     lending.UserID.Hex(),
		Status:     string(lending.Status),
		ReturnDate: timestamppb.New(lending.ReturnDate),
	}, nil
}
