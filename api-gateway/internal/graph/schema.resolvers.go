package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"api-gateway/internal/graph/generated"
	"api-gateway/internal/graph/model"
)

func (r *mutationResolver) RegisterLibrarian(ctx context.Context, input model.NewUser) (*model.User, error) {
	return r.UserGRPCService.CreateUser(ctx, input, model.RoleLibrarian)
}

func (r *mutationResolver) RegisterMember(ctx context.Context, input model.NewUser) (*model.User, error) {
	return r.UserGRPCService.CreateUser(ctx, input, model.RoleMember)
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	return r.UserGRPCService.Login(ctx, input)
}

func (r *mutationResolver) FetchUser(ctx context.Context, input model.FetchUserFilter) (*model.UserPaged, error) {
	return r.UserGRPCService.FetchUser(ctx, input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	return r.UserGRPCService.UpdateUser(ctx, input)
}

func (r *mutationResolver) UpdateSelf(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	return r.UserGRPCService.UpdateSelf(ctx, input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input model.DeleteUser) (*model.User, error) {
	return nil, r.UserGRPCService.DeleteUser(ctx, input)
}

func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	return r.BookGRPCService.CreateBook(ctx, input)
}

func (r *mutationResolver) FetchBook(ctx context.Context, input model.FetchBookFilter) (*model.BookPaged, error) {
	return r.BookGRPCService.FetchBook(ctx, input)
}

func (r *mutationResolver) UpdateBook(ctx context.Context, input model.UpdateBook) (*model.Book, error) {
	return r.BookGRPCService.UpdateBook(ctx, input)
}

func (r *mutationResolver) UpdateBookStock(ctx context.Context, input model.UpdateBookStock) (*model.Book, error) {
	return r.BookGRPCService.UpdateBookStock(ctx, input)
}

func (r *mutationResolver) DeleteBook(ctx context.Context, input model.DeleteBook) (*model.Book, error) {
	return nil, r.BookGRPCService.DeleteBook(ctx, input)
}

func (r *mutationResolver) LendBook(ctx context.Context, input model.NewLending) (*model.Lending, error) {
	return r.LendingGRPCService.LendBook(ctx, input)
}

func (r *mutationResolver) RenewLending(ctx context.Context, input model.RenewLendingRequest) (*model.Lending, error) {
	return r.LendingGRPCService.RenewLending(ctx, input)
}

func (r *mutationResolver) FinishLending(ctx context.Context, input model.FinishLendingRequest) (*model.Lending, error) {
	return r.LendingGRPCService.FinishLending(ctx, input)
}

func (r *mutationResolver) MyLending(ctx context.Context, input *model.MyLendingRequest) (*model.LendingPaged, error) {
	return r.LendingGRPCService.MyLending(ctx, input)
}

func (r *mutationResolver) FetchLending(ctx context.Context, input *model.FetchLendingRequest) (*model.LendingPaged, error) {
	return r.LendingGRPCService.FetchLending(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
