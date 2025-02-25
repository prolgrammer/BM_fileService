package mongo

import (
	mongoConfig "app/config/mongo"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Client struct {
	Database *mongoDriver.Database
	client   *mongoDriver.Client
	cfg      mongoConfig.Config
}

var (
	ErrNoChange = errors.New("no changes applied")
)

func NewClient(cfg mongoConfig.Config) (*Client, error) {
	client := &Client{
		cfg: cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(client.connectionString()).
		SetAuth(options.Credential{
			Username:      cfg.User,
			Password:      cfg.Password,
			AuthMechanism: "SCRAM-SHA-256",
		})

	mongoClient, err := mongoDriver.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo client: %w", err)
	}

	client.client = mongoClient
	client.Database = mongoClient.Database(cfg.Database)

	fmt.Printf("Successful connect to MongoDB at %s\n", client.cfg.Host)

	return client, nil
}

func (c *Client) connectionString() string {
	return fmt.Sprintf("mongodb://%s:%s", c.cfg.Host, c.cfg.Port)
}

func (c *Client) Close(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect client: %w", err)
	}

	fmt.Printf("Successful disconnect client\n")

	return nil
}

func (c *Client) MigrateUp() error {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		c.cfg.User,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.Database,
	)

	m, err := migrate.New(
		c.cfg.MigrationsPath,
		mongoURI)

	if err != nil {
		return fmt.Errorf("failed to create migration handler: %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("No changes to migrations\n")
			return migrate.ErrNoChange
		}
		return fmt.Errorf("failed to migrate up: %w", err)
	}

	fmt.Printf("Migrations applied successfully\n")

	return nil
}
