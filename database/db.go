package database

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
)
type URLStorage interface {
	GetOriginalUrl(shortUrl string) (string, error)
	SaveUrlMapping(shortUrl, originalUrl string) error
	IsOriginalUrlExists(originalUrl string) (string, bool, error)
	IsShortUrlUnique(shortUrl string) bool
	DeleteOriginalUrl(shortUrl string)
	GetAllUrls() ([]string, error)
}

type PostgresStorage struct {
	DB *sql.DB
}

type RedisStorage struct {
	Client *redis.Client
}

func ConnectToDB() (*sql.DB, error) {
	// Инициализация базы данных PostgresSQL
	db, err := sql.Open("postgres", "host=Url.postgres port=5432 user=admin password=admin dbname=URLDB ")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *PostgresStorage) IsShortUrlUnique(shortUrl string) bool {
    var count int
    query := "SELECT COUNT(*) FROM URLDB WHERE short_Url = $1"
    err := p.DB.QueryRow(query, shortUrl).Scan(&count)
    if err != nil {
        fmt.Println("Failure of query:", err)
        return false
    }
    return count == 0
}

func (p *PostgresStorage) SaveUrlMapping(shortUrl, originalUrl string) error {
    query := "INSERT INTO URLDB (short_Url, original_Url) VALUES ($1, $2)"
    _, err := p.DB.Exec(query, shortUrl, originalUrl)
    if err != nil {
        fmt.Println("Failed to execute query:", err)
        return err
    }
    fmt.Printf("Saved URLDB: Original: %s, Short: %s\n", originalUrl, shortUrl)
    return nil
}

func (p *PostgresStorage) IsOriginalUrlExists(originalUrl string) (string, bool, error) {
    var shortUrl string
    query := "SELECT short_Url FROM URLDB WHERE original_Url = $1"
    err := p.DB.QueryRow(query, originalUrl).Scan(&shortUrl)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", false, nil
        }
        return "", false, err
    }
    return shortUrl, true, nil
}

func (p *PostgresStorage) GetOriginalUrl(shortUrl string) (string, error) {
    var originalUrl string
    query := "SELECT original_Url FROM URLDB WHERE short_Url = $1"
    err := p.DB.QueryRow(query, shortUrl).Scan(&originalUrl)
    return originalUrl, err
}

func (p *PostgresStorage) GetAllUrls() ([]string, error) {
    var urls []string
    query := "SELECT original_Url FROM URLDB"
    rows, err := p.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var url string
        if err := rows.Scan(&url); err != nil {
            return nil, err
        }
        urls = append(urls, url)
    }
    return urls, rows.Err()
}

func (p *PostgresStorage) DeleteOriginalUrl(Url string) {
    query := "DELETE FROM URLDB WHERE short_Url = $1"
    _, err := p.DB.Exec(query, Url)
    if err != nil {
        fmt.Println(err)
    }
}




func ConnectToRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *RedisStorage) IsOriginalUrlExists(originalUrl string) (string, bool, error) {
    var shortUrl string
    keys, err := r.Client.Keys("*").Result() // Используйте r.Client
    if err != nil {
        return "", false, err
    }
    for _, key := range keys {
        value, err := r.Client.Get(key).Result() // Используйте r.Client
        if err != nil {
            continue
        }
        if value == originalUrl {
            shortUrl = key
            return shortUrl, true, nil
        }
    }
    return "", false, nil
}

func (r *RedisStorage) SaveUrlMapping(shortUrl, originalUrl string) error {
    err := r.Client.Set(shortUrl, originalUrl, 0).Err() // Используйте r.Client
    if err != nil {
        return err
    }
    return nil
}

func (r *RedisStorage) GetOriginalUrl(shortUrl string) (string, error) {
    originalUrl, err := r.Client.Get(shortUrl).Result() // Используйте r.Client
    if err != nil {
        return "", err
    }
    return originalUrl, nil
}


func (r *RedisStorage) IsShortUrlUnique(shortUrl string) bool {
    // Проверяем, существует ли ключ shortUrl в Redis
    exists, err := r.Client.Exists(shortUrl).Result()
    if err != nil {
        fmt.Println("Error checking URL existence in Redis:", err)
        return false
    }
    // Если ключа нет, то URL уникален
    return exists == 0
}


func (r *RedisStorage) DeleteOriginalUrl(shortUrl string) {
    err := r.Client.Del(shortUrl).Err() // Используйте r.Client
    if err != nil {
        fmt.Println("Error deleting URL in Redis:", err)
    }
}

func (r *RedisStorage) GetAllUrls() ([]string, error) {
    var urls []string
    iter := r.Client.Scan(0, "*", 0).Iterator() // Используйте r.Client
    for iter.Next() {
        urls = append(urls, iter.Val())
    }
    if err := iter.Err(); err != nil {
        return nil, err
    }
    return urls, nil
}


