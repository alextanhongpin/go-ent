package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alextanhongpin/go-ent/ent"
	"github.com/alextanhongpin/go-ent/ent/user"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func main() {
	client := Open("postgresql://john:123456@127.0.0.1/ent")
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	ctx := context.Background()
	if _, err := CreateUser(ctx, client); err != nil {
		log.Fatal(err)
	}
	if _, err := QueryUser(ctx, client); err != nil {
		log.Fatal(err)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ("a8m")).
		Only(ctx) // `Only` fails if no user found, or more than 1 user returned.
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}

	log.Println("car was created: ", ford)

	// Create a new user, and add in the 2 cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}
