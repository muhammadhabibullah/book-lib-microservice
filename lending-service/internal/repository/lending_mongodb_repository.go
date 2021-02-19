package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"lending-service/internal/domain"
	"lending-service/internal/domain/constant"
	"lending-service/pkg/mongodb"
)

type lendingMongoDBRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewLendingMongoDBRepository() domain.LendingRepository {
	db := mongodb.GetDatabase()
	return &lendingMongoDBRepository{
		db:         db,
		collection: db.Collection(constant.LendingCollection),
	}
}

func (r *lendingMongoDBRepository) Create(ctx context.Context, lending *domain.Lending) error {
	lending.ID = primitive.NewObjectID()
	lending.Meta.Create()

	result, err := r.collection.InsertOne(ctx, lending)
	if err != nil {
		return err
	}

	lending.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *lendingMongoDBRepository) Fetch(ctx context.Context, param map[string]interface{}) ([]domain.Lending, error) {
	cursor, err := r.collection.Find(ctx, r.filterBy(param), r.pageBy(param))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	lendings := make([]domain.Lending, 0)
	for cursor.Next(ctx) {
		var lending domain.Lending
		if err = cursor.Decode(&lending); err != nil {
			return nil, err
		}
		lendings = append(lendings, lending)
	}

	return lendings, nil
}

func (r *lendingMongoDBRepository) Count(ctx context.Context, param map[string]interface{}) (int, error) {
	count, err := r.collection.CountDocuments(ctx, r.filterBy(param))
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (*lendingMongoDBRepository) filterBy(param map[string]interface{}) bson.D {
	filter := bson.D{}

	filter = append(filter, bson.E{"meta.deleted_at", nil})

	for key, value := range param {
		switch key {
		case "status":
			filter = append(filter, bson.E{key, value})
		case "user_id", "book_id":
			objectID, _ := primitive.ObjectIDFromHex(value.(string))
			filter = append(filter, bson.E{key, objectID})
		}
	}

	return filter
}

func (*lendingMongoDBRepository) pageBy(param map[string]interface{}) *options.FindOptions {
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

func (r *lendingMongoDBRepository) FindOne(ctx context.Context, filter bson.D) (lending domain.Lending, err error) {
	var decoded bson.M
	err = r.collection.FindOne(ctx, filter).
		Decode(&decoded)
	if err != nil {
		return domain.Lending{}, err
	}

	bsonBytes, _ := bson.Marshal(decoded)
	err = bson.Unmarshal(bsonBytes, &lending)
	return
}

func (r *lendingMongoDBRepository) FindByID(ctx context.Context, id string) (domain.Lending, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Lending{}, err
	}

	filter := bson.D{{"_id", objectID}}
	return r.FindOne(ctx, filter)
}

func (r *lendingMongoDBRepository) Update(ctx context.Context, lending *domain.Lending) error {
	lending.Meta.Update()

	filter := bson.D{{"_id", lending.ID}}
	update := bson.D{{"$set", lending}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var decode bson.M
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&decode)
	if err != nil {
		return err
	}

	bsonBytes, _ := bson.Marshal(decode)
	err = bson.Unmarshal(bsonBytes, lending)
	return nil
}
