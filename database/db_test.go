package database

import (
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"testing"
//	"github.com/alicebob/miniredis/v2"
	//"github.com/go-redis/redis"
)

func TestIsShortUrlUnique_UniqueShortUrl(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pgStorage := PostgresStorage{DB: db}
	// Создаем тестовые данные и ожидаемый результат
	shortUrl := "abc123"
	expectedCount := 0

	// Ожидаем, что запрос вернет ожидаемое количество записей
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortUrl).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	// Вызываем тестируемую функцию
	isUnique := pgStorage.IsShortUrlUnique(shortUrl)

	// Проверяем, что результат соответствует ожиданиям
	if isUnique != true {
		t.Errorf("Expected true, but got false")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestIsShortUrlUnique_NonUniqueShortUrl(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pgStorage := PostgresStorage{DB: db}

	// Создаем тестовые данные и ожидаемый результат
	shortUrl := "abc123"
	expectedCount := 1

	// Ожидаем, что запрос вернет ожидаемое количество записей
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortUrl).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	// Вызываем тестируемую функцию
	isUnique := pgStorage.IsShortUrlUnique(shortUrl)

	// Проверяем, что результат соответствует ожиданиям
	if isUnique != false {
		t.Errorf("Expected false, but got true")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestIsShortUrlUnique_QueryError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pgStorage := PostgresStorage{DB: db}

	// Создаем тестовые данные
	shortUrl := "abc123"

	// Ожидаем, что запрос вызовет ошибку
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortUrl).
		WillReturnError(errors.New("database error"))

	// Вызываем тестируемую функцию
	isUnique := pgStorage.IsShortUrlUnique(shortUrl)

	// Проверяем, что функция возвращает false (ошибка при выполнении запроса)
	if isUnique != false {
		t.Errorf("Expected false, but got true")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestIsShortUrlUnique_ScanError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pgStorage := PostgresStorage{DB: db}

	// Создаем тестовые данные
	shortUrl := "abc123"

	// Ожидаем, что запрос вернет некорректные данные
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortUrl).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("invalid"))

	// Вызываем тестируемую функцию
	isUnique := pgStorage.IsShortUrlUnique(shortUrl)

	// Проверяем, что функция возвращает false (ошибка при сканировании результата)
	if isUnique != false {
		t.Errorf("Expected false, but got true")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSaveUrlMapping_Success(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pgStorage := PostgresStorage{DB: db}

	// Создаем тестовые данные
	shortUrl := "abc123"
	originalUrl := "https://example.com"

	// Ожидаем, что запрос на вставку будет выполнен
	mock.ExpectQuery("INSERT INTO URLDB").WithArgs(shortUrl, originalUrl).
		WillReturnRows(sqlmock.NewRows([]string{"last_insert_id"}).AddRow(1))

	// Вызываем тестируемую функцию
	err = pgStorage.SaveUrlMapping(shortUrl, originalUrl)

	// Проверяем, что ошибка равна nil
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}



func TestSaveUrlMapping_ScanError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close() // Закрываем базу данных, чтобы вызвать ошибку при подключении
	pgStorage := PostgresStorage{DB: db}
	

	// Создаем тестовые данные
	shortUrl := "abc123"
	originalUrl := "https://example.com"

	// Ожидаем, что запрос вернет некорректные данные
	mock.ExpectQuery("INSERT INTO URLDB").WithArgs(shortUrl, originalUrl).
		WillReturnRows(sqlmock.NewRows([]string{"last_insert_id"}).AddRow("invalid"))

	// Вызываем тестируемую функцию
	err = pgStorage.SaveUrlMapping(shortUrl, originalUrl)

	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSaveUrlMapping_LastInsertIDError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close() 
	pgStorage := PostgresStorage{DB: db}
	
	// Создаем тестовые данные
	shortUrl := "abc123"
	originalUrl := "https://example.com"

	// Ожидаем, что запрос на вставку будет выполнен, но не вернет last_insert_id
	mock.ExpectQuery("INSERT INTO Url_mapping").WithArgs(shortUrl, originalUrl).
		WillReturnRows(sqlmock.NewRows([]string{"last_insert_id"}))

	// Вызываем тестируемую функцию
	err = pgStorage.SaveUrlMapping(shortUrl, originalUrl)

	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}




