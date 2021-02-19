package grpc

import (
	"context"
	"log"

	"api-gateway/internal/graph/model"
	proto "api-gateway/pkg/proto/gen"
)

type BookGRPCService struct {
	client proto.BookServiceClient
}

func NewBookGRPCService(
	client proto.BookServiceClient,
) *BookGRPCService {
	return &BookGRPCService{
		client: client,
	}
}

func (c *BookGRPCService) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	book, err := c.client.CreateBook(ctx, &proto.CreateBookRequest{
		Title: input.Title,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.Book{
		ID:    book.GetId(),
		Title: book.GetTitle(),
		Stock: int(book.GetStock()),
	}, nil
}

func (c *BookGRPCService) FetchBook(ctx context.Context, filter model.FetchBookFilter) (*model.BookPaged, error) {
	var (
		limit int32 = 10
		page  int32 = 1
		title string
	)
	if filter.Limit != nil {
		limit = int32(*filter.Limit)
	}
	if filter.Page != nil {
		page = int32(*filter.Page)
	}
	if filter.Title != nil {
		title = *filter.Title
	}

	fetchBookResponse, err := c.client.FetchBook(ctx, &proto.FetchBookRequest{
		Pagination: &proto.BookPaginationRequest{
			Limit: limit,
			Page:  page,
		},
		Title: title,
	})
	if err != nil {
		return nil, err
	}

	books := make([]*model.Book, 0)
	for _, fetchedBook := range fetchBookResponse.Books {
		books = append(books, &model.Book{
			ID:    fetchedBook.GetId(),
			Title: fetchedBook.GetTitle(),
			Stock: int(fetchedBook.GetStock()),
		})
	}

	return &model.BookPaged{
		Books:     books,
		Page:      int(fetchBookResponse.Pagination.GetPage()),
		Limit:     int(fetchBookResponse.Pagination.GetLimit()),
		TotalBook: int(fetchBookResponse.Pagination.GetTotal()),
		LastPage:  int(fetchBookResponse.Pagination.GetLastPage()),
	}, nil
}

func (c *BookGRPCService) FindByID(ctx context.Context, id string) (*model.Book, error) {
	book, err := c.client.FindByID(ctx, &proto.FindBookByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &model.Book{
		ID:    book.GetId(),
		Title: book.GetTitle(),
		Stock: int(book.GetStock()),
	}, nil
}

func (c *BookGRPCService) FindByTitle(ctx context.Context, title string) (*model.Book, error) {
	book, err := c.client.FindByTitle(ctx, &proto.FindBookByTitleRequest{
		Title: title,
	})
	if err != nil {
		return nil, err
	}

	return &model.Book{
		ID:    book.GetId(),
		Title: book.GetTitle(),
		Stock: int(book.GetStock()),
	}, nil
}

func (c *BookGRPCService) UpdateBook(ctx context.Context, input model.UpdateBook) (*model.Book, error) {
	book, err := c.client.UpdateBook(ctx, &proto.UpdateBookRequest{
		Id:    input.ID,
		Title: input.Title,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.Book{
		ID:    book.GetId(),
		Title: book.GetTitle(),
		Stock: int(book.GetStock()),
	}, nil
}

func (c *BookGRPCService) UpdateBookStock(ctx context.Context, input model.UpdateBookStock) (*model.Book, error) {
	book, err := c.client.UpdateBookStock(ctx, &proto.UpdateBookStockRequest{
		Id:          input.ID,
		StockChange: int32(input.StockChange),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.Book{
		ID:    book.GetId(),
		Title: book.GetTitle(),
		Stock: int(book.GetStock()),
	}, nil
}

func (c *BookGRPCService) DeleteBook(ctx context.Context, input model.DeleteBook) error {
	_, err := c.client.DeleteBook(ctx, &proto.DeleteBookRequest{
		Id: input.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
