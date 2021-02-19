package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"user-service/internal/domain"
	"user-service/internal/domain/constant"
	"user-service/pkg/mongodb"
)

type userMongoDBRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserMongoDBRepository() domain.UserRepository {
	db := mongodb.GetDatabase()
	return &userMongoDBRepository{
		db:         db,
		collection: db.Collection(constant.UserCollection),
	}
}

func (r *userMongoDBRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	user.Meta.Create()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userMongoDBRepository) Fetch(ctx context.Context, param map[string]interface{}) ([]domain.User, error) {
	cursor, err := r.collection.Find(ctx, r.filterBy(param), r.pageBy(param))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]domain.User, 0)
	for cursor.Next(ctx) {
		var user domain.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userMongoDBRepository) Count(ctx context.Context, param map[string]interface{}) (int, error) {
	count, err := r.collection.CountDocuments(ctx, r.filterBy(param))
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (*userMongoDBRepository) filterBy(param map[string]interface{}) bson.D {
	filter := bson.D{}

	filter = append(filter, bson.E{"meta.deleted_at", nil})

	for key, value := range param {
		switch key {
		case "role", "email":
			filter = append(filter, bson.E{key, value})
		}
	}

	return filter
}

func (*userMongoDBRepository) pageBy(param map[string]interface{}) *options.FindOptions {
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

func (r *userMongoDBRepository) FindOne(ctx context.Context, filter bson.D) (user domain.User, err error) {
	var decoded bson.M
	err = r.collection.FindOne(ctx, filter).
		Decode(&decoded)
	if err != nil {
		return domain.User{}, err
	}

	bsonBytes, _ := bson.Marshal(decoded)
	err = bson.Unmarshal(bsonBytes, &user)
	return
}

func (r *userMongoDBRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.D{{"_id", objectID}}
	return r.FindOne(ctx, filter)
}

func (r *userMongoDBRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	filter := bson.D{{"email", email}}
	return r.FindOne(ctx, filter)
}

func (r *userMongoDBRepository) Update(ctx context.Context, user *domain.User) error {
	user.Meta.Update()

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", user}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var decode bson.M
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&decode)
	if err != nil {
		return err
	}

	bsonBytes, _ := bson.Marshal(decode)
	err = bson.Unmarshal(bsonBytes, user)
	return nil
}

func (r *userMongoDBRepository) Delete(ctx context.Context, user *domain.User) error {
	user.Meta.Delete()

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", user}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
