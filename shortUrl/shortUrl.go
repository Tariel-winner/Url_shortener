package shortUrl

import (
    "crypto/rand"
    "Service/database"
    "fmt"
)

const UrlLength = 10
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateShortUrl(originalUrl string, storage database.URLStorage) (string, error) {
    var (
        shortUrl string
        exist    bool = false
        err      error
    )

    if shortUrl, exist, err = storage.IsOriginalUrlExists(originalUrl); exist == false && err == nil {
        for i := 0; i < 5; i++ {
            randomBytes := make([]byte, UrlLength)
            _, err = rand.Read(randomBytes)
            if err != nil {
                return "", err
            }

            shortUrl = ""
            for i := 0; i < UrlLength; i++ {
                index := int(randomBytes[i]) % len(charset)
                shortUrl += string(charset[index])
            }

            if storage.IsShortUrlUnique(shortUrl) {
                err := storage.SaveUrlMapping(shortUrl, originalUrl)
                if err != nil {
                    return "", err
                }
                return shortUrl, nil
            }
        }

        // Измененная строка: возвращаем ошибку, если не удалось сгенерировать уникальный URL
        return "", fmt.Errorf("failed to generate a unique short URL after 5 attempts")
    } else if err != nil {
        return "", err
    }

    return shortUrl, nil
}
