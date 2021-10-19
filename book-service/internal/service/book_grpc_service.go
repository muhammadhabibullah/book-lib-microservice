package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"book-service/internal/domain"
	"book-service/internal/repository"
	"book-service/pkg/proto"
)

type BookGRPCService struct {
	proto.UnimplementedBookServiceServer
	bookRepository domain.BookRepository
}

func NewBookGRPCService() *BookGRPCService {
	return &BookGRPCService{
		bookRepository: repository.NewBookMongoDBRepository(),
	}
}

func (s *BookGRPCService) CreateBook(ctx context.Context, request *proto.CreateBookRequest) (*proto.Book, error) {
	book := domain.Book{
		Title: request.Title,
	}

	if err := s.bookRepository.Create(ctx, &book); err != nil {
		return nil, err
	}

	return &proto.Book{
		Id:    book.ID.Hex(),
		Title: book.Title,
		Stock: int32(book.Stock),
	}, nil
}

func (s *BookGRPCService) FetchBook(ctx context.Context, request *proto.FetchBookRequest) (*proto.FetchBookResponse, error) {
	page, limit := request.Pagination.Page, request.Pagination.Limit
	fetchFilter := map[string]interface{}{
		"page":  page,
		"limit": limit,
	}
	if title := request.Title; title != "" {
		fetchFilter["title"] = title
	}

	books, err := s.bookRepository.Fetch(ctx, fetchFilter)
	if err != nil {
		return nil, err
	}

	protoBooks := make([]*proto.Book, 0)
	for _, book := range books {
		protoBooks = append(protoBooks, &proto.Book{
			Id:    book.ID.Hex(),
			Title: book.Title,
			Stock: int32(book.Stock),
		})
	}

	total, err := s.bookRepository.Count(ctx, fetchFilter)
	if err != nil {
		return nil, err
	}

	return &proto.FetchBookResponse{
		Pagination: &proto.BookPaginationResponse{
			Limit:    limit,
			Page:     page,
			LastPage: int32(math.Ceil(float64(total) / float64(limit))),
			Total:    int32(total),
		},
		Books: protoBooks,
	}, nil
}

func (s *BookGRPCService) FindByID(ctx context.Context, request *proto.FindBookByIDRequest) (*proto.Book, error) {
	book, err := s.bookRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &proto.Book{
		Id:    book.ID.Hex(),
		Title: book.Title,
		Stock: int32(book.Stock),
	}, nil
}

func (s *BookGRPCService) FindByTitle(ctx context.Context, request *proto.FindBookByTitleRequest) (*proto.Book, error) {
	book, err := s.bookRepository.FindByTitle(ctx, request.Title)
	if err != nil {
		return nil, err
	}

	return &proto.Book{
		Id:    book.ID.Hex(),
		Title: book.Title,
		Stock: int32(book.Stock),
	}, nil
}

func (s *BookGRPCService) UpdateBook(ctx context.Context, request *proto.UpdateBookRequest) (*proto.Book, error) {
	book, err := s.bookRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	book.Title = request.Title

	err = s.bookRepository.Update(ctx, &book)
	if err != nil {
		return nil, err
	}

	return &proto.Book{
		Id:    book.ID.Hex(),
		Title: book.Title,
		Stock: int32(book.Stock),
	}, nil
}

func (s *BookGRPCService) UpdateBookStock(ctx context.Context, request *proto.UpdateBookStockRequest) (*proto.Book, error) {
	if request.StockChange == 0 {
		return nil, errors.New("stock change requested is 0")
	}

	retry := 3

	for {
		book, err := s.bookRepository.FindByID(ctx, request.Id)
		if err != nil {
			return nil, err
		}

		newStock := book.Stock + int(request.StockChange)
		if newStock < 0 {
			return nil, errors.New("stock cannot be decreased to below 0")
		}

		err = s.bookRepository.UpdateStock(ctx, &book, newStock)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) && retry >= 0 {
				time.Sleep(200 * time.Millisecond)
				retry++
				continue
			}
			if retry == 0 {
				err = fmt.Errorf("failed update stock after retry 3 times: %w", err)
			}

			return nil, err
		}

		return &proto.Book{
			Id:    book.ID.Hex(),
			Title: book.Title,
			Stock: int32(book.Stock),
		}, nil
	}
}

func (s *BookGRPCService) DeleteBook(ctx context.Context, request *proto.DeleteBookRequest) (*proto.DeleteBookResponse, error) {
	book, err := s.bookRepository.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	err = s.bookRepository.Delete(ctx, &book)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteBookResponse{}, nil
}
