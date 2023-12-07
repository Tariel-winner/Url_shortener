package shortUrl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"errors" 
)

// MockURLStorage - мок для интерфейса URLStorage
type MockURLStorage struct {
	mock.Mock
}

func (m *MockURLStorage) GetOriginalUrl(shortUrl string) (string, error) {
	args := m.Called(shortUrl)
	return args.String(0), args.Error(1)
}

func (m *MockURLStorage) SaveUrlMapping(shortUrl, originalUrl string) error {
	args := m.Called(shortUrl, originalUrl)
	return args.Error(0)
}

func (m *MockURLStorage) IsOriginalUrlExists(originalUrl string) (string, bool, error) {
	args := m.Called(originalUrl)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockURLStorage) IsShortUrlUnique(shortUrl string) bool {
	args := m.Called(shortUrl)
	return args.Bool(0)
}

func (m *MockURLStorage) DeleteOriginalUrl(shortUrl string) {
	m.Called(shortUrl)
}

func (m *MockURLStorage) GetAllUrls() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}




	func TestGenerateShortUrl_OriginalUrlExists(t *testing.T) {
		mockStorage := new(MockURLStorage)
		originalUrl := "https://example.com"
		existingShortUrl := "abcde"
	
		mockStorage.On("IsOriginalUrlExists", originalUrl).Return(existingShortUrl, true, nil)
	
		result, err := GenerateShortUrl(originalUrl, mockStorage)
	
		assert.NoError(t, err)
		assert.Equal(t, existingShortUrl, result)
		mockStorage.AssertExpectations(t)
	}
	
	func TestGenerateShortUrl_ErrorInIsOriginalUrlExists(t *testing.T) {
		mockStorage := new(MockURLStorage)
		originalUrl := "https://example.com"
	
		mockStorage.On("IsOriginalUrlExists", originalUrl).Return("", false, errors.New("database error"))
	
		result, err := GenerateShortUrl(originalUrl, mockStorage)
	
		assert.Error(t, err)
		assert.Empty(t, result)
		mockStorage.AssertExpectations(t)
	}
	
	func TestGenerateShortUrl_ErrorInSaveUrlMapping(t *testing.T) {
		mockStorage := new(MockURLStorage)
		originalUrl := "https://example.com"
		shortUrl := "abcde"
	
		mockStorage.On("IsOriginalUrlExists", originalUrl).Return("", false, nil)
		mockStorage.On("IsShortUrlUnique", mock.Anything).Return(true)
		mockStorage.On("SaveUrlMapping", mock.Anything, mock.Anything).Return(errors.New("database error"))
	
		result, err := GenerateShortUrl(originalUrl, mockStorage)
	
		assert.Error(t, err)
		assert.NotEqual(t, shortUrl, result)
		mockStorage.AssertExpectations(t)
	}
	
	func TestGenerateShortUrl_FailedToGenerateUniqueShortUrl(t *testing.T) {
		mockStorage := new(MockURLStorage)
		originalUrl := "https://example.com"
	
		mockStorage.On("IsOriginalUrlExists", originalUrl).Return("", false, nil)
		mockStorage.On("IsShortUrlUnique", mock.Anything).Return(false)
	
		result, err := GenerateShortUrl(originalUrl, mockStorage)
	
		assert.Error(t, err)
		assert.Empty(t, result)
		mockStorage.AssertNumberOfCalls(t, "IsShortUrlUnique", 5)
	}