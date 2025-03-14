package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	embed "embed"
	"github.com/go-redis/redis/v8"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DB struct {
	*sql.DB
	Redis *redis.Client
}

func NewDB() (*DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	for i := 0; i < 15; i++ {
		err = db.PingContext(context.Background())
		if err == nil {
			break
		}
		log.Printf("Ошибка подключения к БД. Повторная попытка через 2 секунды... (%d/15)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %v", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	return &DB{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func runMigrations(db *sql.DB) error {
	files, err := migrations.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("ошибка чтения миграций: %w", err)
	}

	for _, f := range files {
		if !f.IsDir() {
			content, err := migrations.ReadFile("migrations/" + f.Name())
			if err != nil {
				return fmt.Errorf("ошибка чтения файла %s: %w", f.Name(), err)
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return fmt.Errorf("ошибка выполнения миграции %s: %w", f.Name(), err)
			}
			log.Printf("Выполнена миграция: %s", f.Name())
		}
	}
	return nil
}
