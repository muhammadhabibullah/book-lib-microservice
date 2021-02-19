package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"book-service/internal/domain"
	"book-service/internal/domain/constant"
	"book-service/pkg/mongodb"
)

type bookMongoDBRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewBookMongoDBRepository() domain.BookRepository {
	db := mongodb.GetDatabase()
	return &bookMongoDBRepository{
		db:         db,
		collection: db.Collection(constant.BookCollection),
	}
}

func (r *bookMongoDBRepository) Create(ctx context.Context, book *domain.Book) error {
	book.ID = primitive.NewObjectID()
	book.Meta.Create()

	result, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		return err
	}

	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *bookMongoDBRepository) Fetch(ctx context.Context, param map[string]interface{}) ([]domain.Book, error) {
	cursor, err := r.collection.Find(ctx, r.filterBy(param), r.pageBy(param))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	books := make([]domain.Book, 0)
	for cursor.Next(ctx) {
		var book domain.Book
		if err = cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (r *bookMongoDBRepository) Count(ctx context.Context, param map[string]interface{}) (int, error) {
	count, err := r.collection.CountDocuments(ctx, r.filterBy(param))
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (*bookMongoDBRepository) filterBy(param map[string]interface{}) bson.D {
	filter := bson.D{}

	filter = append(filter, bson.E{"meta.deleted_at", nil})

	for key, value := range param {
		switch key {
		case "title":
			filter = append(filter, bson.E{key, primitive.Regex{Pattern: (value).(string), Options: "i"}})
		}
	}

	return filter
}

func (*bookMongoDBRepository) pageBy(param map[string]interface{}) *options.FindOptions {
	limit, ok := param["limit"].(int32)
	if !ok || limit <= 0 {
		limit = 10
	}
	page, ok := param["page"].(int32)
	if !ok || page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"meta.created_at": -1})

	return findOptions
}

func (r *bookMongoDBRepository) FindOne(ctx context.Context, filter bson.D) (book domain.Book, err error) {
	var decoded bson.M
	err = r.collection.FindOne(ctx, filter).
		Decode(&decoded)
	if err != nil {
		return domain.Book{}, err
	}

	bsonBytes, _ := bson.Marshal(decoded)
	err = bson.Unmarshal(bsonBytes, &book)
	return
}

func (r *bookMongoDBRepository) FindByID(ctx context.Context, id string) (domain.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Book{}, err
	}

	filter := bson.D{{"_id", objectID}}
	return r.FindOne(ctx, filter)
}

func (r *bookMongoDBRepository) FindByTitle(ctx context.Context, title string) (domain.Book, error) {
	filter := bson.D{{"title", primitive.Regex{Pattern: title, Options: "i"}}}
	return r.FindOne(ctx, filter)
}

func (r *bookMongoDBRepository) Update(ctx context.Context, book *domain.Book) error {
	book.Meta.Update()

	filter := bson.D{{"_id", book.ID}}
	update := bson.D{{"$set", book}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var decode bson.M
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&decode)
	if err != nil {
		return err
	}

	bsonBytes, _ := bson.Marshal(decode)
	err = bson.Unmarshal(bsonBytes, book)
	return nil
}

func (r *bookMongoDBRepository) UpdateStock(ctx context.Context, book *domain.Book, newStock int) error {
	book.Meta.Update()

	filter := bson.D{{"_id", book.ID}, {"stock", book.Stock}}
	update := bson.D{{"$set", bson.D{
		{"stock", newStock},
		{"meta.updated_at", book.UpdatedAt},
	}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var decode bson.M
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&decode)
	if err != nil {
		return err
	}

	bsonBytes, _ := bson.Marshal(decode)
	err = bson.Unmarshal(bsonBytes, book)
	return nil
}

func (r *bookMongoDBRepository) Delete(ctx context.Context, book *domain.Book) error {
	book.Meta.Delete()

	filter := bson.D{{"_id", book.ID}}
	update := bson.D{{"$set", book}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
