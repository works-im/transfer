package transfer

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB database transfer
type MongoDB struct {
	client     *mongo.Client
	database   string
	collection string

	Mapping Mapping
}

// NewMongoDB return MongoDB transfer
func NewMongoDB(args *DatabaseOptions) (*MongoDB, error) {

	var (
		opts = options.Client()
		host = fmt.Sprintf("%s:%d", args.Driver.Host, args.Driver.Port)
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

	// connect database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	db := &MongoDB{
		client:     client,
		database:   args.Driver.Database,
		collection: args.Driver.Table,
		Mapping:    args.Mapping,
	}

	return db, nil
}

// Reader database
// query: aggregation https://docs.mongodb.com/manual/aggregation/
func (mongo *MongoDB) Reader(query Query) (packet Packet, err error) {

	collection := mongo.client.Database(mongo.database).Collection(mongo.collection)

	pipeline := []M{}

	if err := query.UnmarshalQuery(&pipeline); err != nil {
		return nil, err
	}

	// select fields
	if len(mongo.Mapping) > 0 {
		fields := M{}
		for _, field := range mongo.Mapping {
			fields[field.Target] = "$" + field.Source
		}
		pipeline = append(pipeline, M{"$project": fields})
	}

	// offset
	pipeline = append(pipeline, M{"$skip": (query.Page - 1) * query.Size})

	// page limit
	if query.Size > 0 {
		pipeline = append(pipeline, M{"$limit": query.Size})
	}

	// set context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var row M
		if err := cur.Decode(&row); err != nil {
			return nil, err
		}

		for field, value := range row {
			// ObjectID to string
			if id, ok := value.(primitive.ObjectID); ok {
				row[field] = id.Hex()
			}
		}

		packet = append(packet, row)
	}

	return packet, cur.Err()
}

// Writer data
func (mongo *MongoDB) Writer(packet Packet) error {
	return nil
}
