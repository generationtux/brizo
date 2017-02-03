package resources

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

// PersonalAccessTokenLength is used to create and validate personal access tokens
const PersonalAccessTokenLength = 128

func generateRandomString(length int) string {
	const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[r.Intn(len(characters))]
	}

	return string(bytes)
}

// AccessToken as defined by Brizo.
type AccessToken struct {
	database.Model
	Token string `gorm:"not null;unique_index" sql:"type:varchar(128)" json:"token"`
}

// BeforeCreate is a hook that runs before inserting a new record into the database
func (a *AccessToken) BeforeCreate() (err error) {
	if a.Token == "" {
		a.Token = generateRandomString(PersonalAccessTokenLength)
	}

	return
}

// CreateAccessToken will add a new AccessToken to Brizo.
func CreateAccessToken(db *gorm.DB, token *AccessToken) (bool, error) {
	result := db.Create(&token)

	return result.RowsAffected == 1, result.Error
}

// CreateRandomAccessToken will create a new AccessToken without requiring a
// pepared access token instance.
func CreateRandomAccessToken(db *gorm.DB) (*AccessToken, error) {
	token := &AccessToken{}
	token.Token = generateRandomString(PersonalAccessTokenLength)
	successful, err := CreateAccessToken(db, token)
	if err != nil || !successful {
		log.Printf("Error when creating access token: '%s'\n", err)
		return &AccessToken{}, errors.New("there was an error when creating access token")
	}

	return token, nil
}

// GetAccessToken will get an existing access token by it's token string
func GetAccessToken(db *gorm.DB, token string) (*AccessToken, error) {
	accessToken := new(AccessToken)
	if err := db.Where("token = ?", token).First(&accessToken).Error; err != nil {
		return accessToken, err
	}

	if accessToken.ID == 0 {
		return new(AccessToken), errors.New("not-found")
	}

	return accessToken, nil
}

// HasAccessToken will return the existance of an access token in the database.
// Errors that occur on the passed db instance will not be returned, so the
// caller should check for any errors if desired.
func HasAccessToken(db *gorm.DB, token string) bool {
	if db.Where("token = ?", token).First(&AccessToken{}).RecordNotFound() {
		return false
	}

	return true
}
