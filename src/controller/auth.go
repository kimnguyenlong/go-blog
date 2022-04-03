package controller

import (
	customerror "blog/custom-error"
	"blog/entity"
	"blog/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
}

type authController struct {
	userModel *model.User
}

func NewAuthController(db *mongo.Database) AuthController {
	return &authController{
		userModel: model.NewUser(db),
	}
}

func (ctr authController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user entity.User
		err := ctx.BindJSON(&user)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		newUser, err := ctr.userModel.Save(user)
		if err != nil {
			{
				panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
			}
		}
		jwt, err := newUser.CreateJWT()
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"error":   false,
			"message": "Registered successfully!",
			"jwt":     jwt,
		})
	}
}

func (ctr authController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var credentials struct {
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		err := ctx.BindJSON(&credentials)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.D{
			{
				Key: "$or",
				Value: bson.A{
					bson.D{{Key: "email", Value: credentials.Account}},
					bson.D{{Key: "username", Value: credentials.Account}},
				},
			},
		}
		var user entity.User
		err = ctr.userModel.Base.FindOne(filter).Decode(&user)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		isCorrectPassword := user.CheckPassword(credentials.Password)
		if !isCorrectPassword {
			panic(customerror.NewAPIError("Incorrect password", http.StatusBadRequest))
		}
		jwtString, err := user.CreateJWT()
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Login successfully",
			"jwt":     jwtString,
		})
	}
}
