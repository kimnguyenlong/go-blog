package controller

import (
	customerror "blog/custom-error"
	"blog/entity"
	"blog/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserController interface {
	GetUsers() gin.HandlerFunc
	GetSingleUser() gin.HandlerFunc
	Update() gin.HandlerFunc
}

type userController struct {
	userModel *model.User
}

func NewUserController(db *mongo.Database) UserController {
	return &userController{
		userModel: model.NewUser(db),
	}
}

func (ctr userController) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (ctr userController) GetSingleUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("id")
		uOjbID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.D{
			{
				Key:   "_id",
				Value: uOjbID,
			},
		}
		excludes := bson.D{
			{
				Key:   "password",
				Value: 0,
			},
		}
		opts := options.FindOne().SetProjection(excludes)
		var user entity.User
		err = ctr.userModel.Base.FindOne(filter, opts).Decode(&user)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  user,
		})
	}
}

func (ctr userController) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("id")
		uObjID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var updateData struct {
			Follow string `json:"follow"`
		}
		ctx.BindJSON(&updateData)

		var user entity.User

		if updateData.Follow != "" {
			followUidString := updateData.Follow
			opr := "$addToSet"
			if updateData.Follow[0] == '-' {
				followUidString = updateData.Follow[1:]
				opr = "$pull"
			}
			followObjID, err := primitive.ObjectIDFromHex(followUidString)
			if err != nil {
				panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
			}
			user, err = ctr.userModel.Follow(uObjID, followObjID, opr)
			if err != nil {
				panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Update user successfully!",
			"data":    user,
		})
	}
}
