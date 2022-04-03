package entity

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Username  string               `bson:"username,omitempty" json:"username"`
	Email     string               `bson:"email,omitempty" json:"email"`
	Password  string               `bson:"password,omitempty" json:"password"`
	Following []primitive.ObjectID `bson:"following,omitempty" json:"following"`
	Followers []primitive.ObjectID `bson:"followers,omitempty" json:"followers"`
	Created   int64                `bson:"created,omitempty" json:"created"`
	Updated   int64                `bson:"updated,omitempty" json:"updated"`
}

func (user User) CreateJWT() (string, error) {
	jwtLifeTime, err := strconv.Atoi(os.Getenv("JWT_LIFE_TIME"))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID.Hex(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * time.Duration(jwtLifeTime)).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (user User) CheckPassword(candidatePassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(candidatePassword)) == nil
}
