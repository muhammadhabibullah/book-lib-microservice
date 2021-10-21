package service

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *LendingGRPCService) CreateLending(request *proto.CreateLendingRequest, stream proto.LendingService_CreateLendingServer) error {
	ctx := context.Background()

	userID, err := primitive.ObjectIDFromHex(request.UserId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid user ID: %s", request.UserId)
	}

	book, err := s.bookServiceClient.FindByID(ctx, &proto.FindBookByIDRequest{
		Id: request.BookId,
	})
	if err != nil {
		return err
	}

	if book.Stock == 0 {
		return status.Error(codes.Aborted, "book stock is empty")
	}

	bookID, _ := primitive.ObjectIDFromHex(book.Id)
	lending := domain.Lending{
		BookID:     bookID,
		UserID:     userID,
		Status:     constant.LendingDraft,
		ReturnDate: time.Now().Add(defaultLendingDuration * time.Hour),
	}

	err = s.lendingRepository.Create(ctx, &lending)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = stream.Send(&proto.Lending{
		Id:         lending.ID.Hex(),
		BookId:     lending.BookID.Hex(),
		UserId:     lending.UserID.Hex(),
		Status:     string(lending.Status),
		ReturnDate: timestamppb.New(lending.ReturnDate),
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	_, err = s.bookServiceClient.UpdateBookStock(ctx, &proto.UpdateBookStockRequest{
		Id:          book.Id,
		StockChange: -1,
	})
	if err != nil {
		lending.Status = constant.LendingCanceled
		cancelErr := s.cancelCreateLending(ctx, lending, stream)
		if cancelErr != nil {
			return status.Error(codes.Internal, cancelErr.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	lending.Status = constant.LendingActive
	err = s.lendingRepository.Update(ctx, &lending)
	if err != nil {
		_, cancelBookErr := s.bookServiceClient.UpdateBookStock(ctx, &proto.UpdateBookStockRequest{
			Id:          book.Id,
			StockChange: 1,
		})
		if cancelBookErr != nil {
			return status.Error(codes.Internal, cancelBookErr.Error())
		}

		cancelErr := s.cancelCreateLending(ctx, lending, stream)
		if cancelErr != nil {
			return status.Error(codes.Internal, cancelErr.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	err = stream.Send(&proto.Lending{
		Id:         lending.ID.Hex(),
		BookId:     lending.BookID.Hex(),
		UserId:     lending.UserID.Hex(),
		Status:     string(lending.Status),
		ReturnDate: timestamppb.New(lending.ReturnDate),
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *LendingGRPCService) cancelCreateLending(ctx context.Context, lending domain.Lending, stream proto.LendingService_CreateLendingServer) error {
	err := s.lendingRepository.Update(ctx, &lending)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = stream.Send(&proto.Lending{
		Id:         lending.ID.Hex(),
		BookId:     lending.BookID.Hex(),
		UserId:     lending.UserID.Hex(),
		Status:     string(lending.Status),
		ReturnDate: timestamppb.New(lending.ReturnDate),
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *LendingGRPCService) FetchLending(ctx context.Context, request *proto.FetchLendingRequest) (*proto.FetchLendingResponse, error) {
	page, limit := request.Pagination.Page, request.Pagination.Limit
	fetchFilter := map[string]interface{}{
		"page":  page,
		"limit": limit,
	}
	if requestStatus := request.Status; requestStatus != "" {
		fetchFilter["status"] = requestStatus
	}
	if userID := request.UserId; userID != "" {
		fetchFilter["user_id"] = userID
	}
	if bookID := request.BookId; bookID != "" {
		fetchFilter["book_id"] = bookID
	}

	lendings, err := s.lendingRepository.Fetch(ctx, fetchFilter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "lending with %s ID is not found", request.Id)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	lending.ReturnDate = time.Now().Add(defaultLendingDuration * time.Hour)

	err = s.lendingRepository.Update(ctx, &lending)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "lending with %s ID is not found", request.Id)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	lending.Status = constant.LendingInactive

	err = s.lendingRepository.Update(ctx, &lending)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
