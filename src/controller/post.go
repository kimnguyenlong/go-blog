package controller

import (
	customerror "blog/custom-error"
	"blog/entity"
	"blog/model"
	"blog/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostController interface {
	CreatePost() gin.HandlerFunc
	GetPosts() gin.HandlerFunc
	GetSinglePost() gin.HandlerFunc
	CreateComment() gin.HandlerFunc
	GetComments() gin.HandlerFunc
}

type postController struct {
	postModel    *model.Post
	commentModel *model.Comment
}

func NewPostController(db *mongo.Database) PostController {
	return &postController{
		postModel:    model.NewPost(db),
		commentModel: model.NewComment(db),
	}
}

func (postController postController) CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var post entity.Post
		err := ctx.BindJSON(&post)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		uid := ctx.GetString("uid")
		userObjId, _ := primitive.ObjectIDFromHex(uid)
		post.UserID = userObjId
		newPost, err := postController.postModel.Save(post)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"error":   false,
			"message": "Create new post successfully",
			"data":    newPost,
		})
	}
}

var sortCriterias = []string{
	"created",
	"updated",
}

func (ctr postController) GetPosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("uid")
		search := ctx.Query("search")
		sorts := ctx.Query("sort")

		filter := bson.D{}
		sortOpts := bson.D{}

		if uid != "" {
			uObjID, err := primitive.ObjectIDFromHex(uid)
			if err != nil {
				panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
			}
			filter = append(filter, bson.E{
				Key:   "user_id",
				Value: uObjID,
			})
		}

		if search != "" {
			pattern := fmt.Sprintf(".*%s.*", search)
			filter = append(filter, bson.E{
				Key: "$or",
				Value: bson.A{
					bson.D{
						{
							Key: "title",
							Value: bson.D{
								{
									Key: "$regex",
									Value: primitive.Regex{
										Pattern: pattern,
										Options: "i",
									},
								},
							},
						},
					},
					bson.D{
						{
							Key: "description",
							Value: bson.D{
								{
									Key: "$regex",
									Value: primitive.Regex{
										Pattern: pattern,
										Options: "i",
									},
								},
							},
						},
					},
					bson.D{
						{
							Key: "content",
							Value: bson.D{
								{
									Key: "$regex",
									Value: primitive.Regex{
										Pattern: pattern,
										Options: "i",
									},
								},
							},
						},
					},
				},
			})
		}

		if sorts != "" {
			for _, sort := range strings.Split(sorts, ",") {
				var order = 1
				var criteria = sort
				if sort[0] == '-' {
					order = -1
					criteria = sort[1:]
				}
				if !utils.ArrayContains(sortCriterias, criteria) {
					continue
				}
				sortOpts = append(sortOpts, bson.E{Key: criteria, Value: order})
			}
		}

		opts := options.Find().SetSort(sortOpts)
		cursor, err := ctr.postModel.Base.Find(filter, opts)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var posts []entity.Post
		err = cursor.All(context.Background(), &posts)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  posts,
		})
	}
}

func (ctr postController) GetSinglePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pid := ctx.Param("id")
		pObjID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.D{
			{
				Key:   "_id",
				Value: pObjID,
			},
		}
		var post entity.Post
		err = ctr.postModel.Base.FindOne(filter).Decode(&post)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  post,
		})
	}
}

func (ctr postController) CreateComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.GetString("uid")
		userObjId, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		pid := ctx.Param("id")
		pObjID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var comment entity.Comment
		err = ctx.BindJSON(&comment)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		comment.UserID = userObjId
		comment.PostID = pObjID
		newComment, err := ctr.commentModel.Save(comment)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"error":   false,
			"message": "Created new comment successfully",
			"data":    newComment,
		})
	}
}

func (ctr postController) GetComments() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pid := ctx.Param("id")
		pObjID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.D{
			{
				Key:   "post_id",
				Value: pObjID,
			},
		}
		var comments []entity.Comment
		cursor, err := ctr.commentModel.Base.Find(filter)
		defer cursor.Close(context.Background())
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		err = cursor.All(context.Background(), &comments)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  comments,
		})
	}
}
