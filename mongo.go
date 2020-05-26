package transfer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoOptions mongodb options
type MongoOptions struct {
	Driver    Driver
	Database  string
	TableName string
	Fields    []Field
}

// MongoDB database transfer
type MongoDB struct {
	client     *mongo.Client
	database   string
	collection string
	fields     []Field
	ctx        context.Context
}

// NewMongoDB return MongoDB transfer
func NewMongoDB(args *MongoOptions) (*MongoDB, error) {

	var (
		opts = options.Client()
		host = fmt.Sprintf("%s:%s", args.Driver.Host, args.Driver.Port)
	)

	opts.Hosts = []string{host}

	if len(args.Driver.Username) > 0 && len(args.Driver.Password) > 0 {
		opts.Auth = &options.Credential{
			Username: args.Driver.Username,
			Password: args.Driver.Password,
		}
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	if len(args.Driver.Database) > 0 {
		client.Database(args.Driver.Database)
	}

	db := &MongoDB{
		client:     client,
		database:   args.Driver.Database,
		collection: args.TableName,
	}

	// set context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db.ctx = ctx

	// connect database
	if err := db.client.Connect(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// Reader database
// query: aggregation https://docs.mongodb.com/manual/aggregation/
func (mongo *MongoDB) Reader(query Query) (result []M, err error) {

	collection := mongo.client.Database(mongo.database).Collection(mongo.collection)

	pipeline := bson.D{}

	if err := query.UnmarshalQuery(&pipeline); err != nil {
		return nil, err
	}

	if query.Page != nil && query.Size != nil {
		pipeline = append(pipeline, bson.E{
			Key:   "$skip",
			Value: query.Page,
		})

		pipeline = append(pipeline, bson.E{
			Key:   "$limit",
			Value: query.Size,
		})
	}

	cur, err := collection.Aggregate(mongo.ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(mongo.ctx)

	for cur.Next(mongo.ctx) {
		var row M
		if err := cur.Decode(&row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, cur.Err()
}

// Writer data
func (mongo *MongoDB) Writer(data []M) error {
	return nil
}
