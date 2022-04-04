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
	"time"

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
	DeletePost() gin.HandlerFunc
	UpdatePost() gin.HandlerFunc

	CreateComment() gin.HandlerFunc
	GetComments() gin.HandlerFunc
	DeleteComment() gin.HandlerFunc
	UpdateComment() gin.HandlerFunc
	CreateReply() gin.HandlerFunc
	GetReplies() gin.HandlerFunc
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

func (ctr postController) DeletePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := primitive.ObjectIDFromHex(ctx.GetString("uid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		pid, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.D{
			{
				Key:   "_id",
				Value: pid,
			},
			{
				Key:   "user_id",
				Value: uid,
			},
		}
		result, err := ctr.postModel.Base.DeleteOne(filter)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		if result.DeletedCount < 1 {
			panic(customerror.NewAPIError("No post found!", http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Delete a post successfully",
		})
	}
}

func (ctr postController) UpdatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := primitive.ObjectIDFromHex(ctx.GetString("uid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		pid, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var updateData entity.Post
		err = ctx.BindJSON(&updateData)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.D{
			{
				Key:   "_id",
				Value: pid,
			},
			{
				Key:   "user_id",
				Value: uid,
			},
		}
		update := bson.D{}

		if updateData.Title != "" {
			update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "title", Value: updateData.Title}}})
		}
		if updateData.Description != "" {
			update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "description", Value: updateData.Description}}})
		}
		if updateData.Content != "" {
			update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "content", Value: updateData.Content}}})
		}
		if len(updateData.Topics) > 0 {
			update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "topics", Value: updateData.Topics}}})
		}

		if len(update) > 0 {
			update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "updated", Value: time.Now().Unix()}}})
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		var updatedPost entity.Post
		err = ctr.postModel.Base.FindOneAndUpdateOne(filter, update, opts).Decode(&updatedPost)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Update a post successfully",
			"data":    updatedPost,
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
			{
				Key:   "parent_id",
				Value: nil,
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

func (ctr postController) DeleteComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := primitive.ObjectIDFromHex(ctx.GetString("uid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		pid, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		cid, err := primitive.ObjectIDFromHex(ctx.Param("cid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.M{
			"_id":     cid,
			"user_id": uid,
			"post_id": pid,
		}
		result, err := ctr.commentModel.Base.DeleteOne(filter)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		if result.DeletedCount < 1 {
			panic(customerror.NewAPIError("No comment found!", http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Delete a comment successfully",
		})
	}
}

func (ctr postController) UpdateComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := primitive.ObjectIDFromHex(ctx.GetString("uid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		pid, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		cid, err := primitive.ObjectIDFromHex(ctx.Param("cid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var updateData entity.Comment
		err = ctx.BindJSON(&updateData)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.M{
			"_id":     cid,
			"user_id": uid,
			"post_id": pid,
		}
		update := bson.M{}
		if updateData.Content != "" {
			update["$set"] = bson.M{"content": updateData.Content, "updated": time.Now().Unix()}
		}
		var updatedComment entity.Comment
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = ctr.commentModel.Base.FindOneAndUpdateOne(filter, update, opts).Decode(&updatedComment)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Update a comment successfully",
			"data":    updatedComment,
		})
	}
}

func (ctr postController) CreateReply() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := primitive.ObjectIDFromHex(ctx.GetString("uid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		pid, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		cid, err := primitive.ObjectIDFromHex(ctx.Param("cid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var reply entity.Comment
		err = ctx.BindJSON(&reply)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		now := time.Now().Unix()
		reply.ParentID = cid
		reply.UserID = uid
		reply.PostID = pid
		reply.Created = now
		reply.Updated = now
		newComment, err := ctr.commentModel.Save(reply)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Reply a comment successfully",
			"data":    newComment,
		})
	}
}

func (ctr postController) GetReplies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cid, err := primitive.ObjectIDFromHex(ctx.Param("cid"))
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		filter := bson.M{
			"parent_id": cid,
		}
		var replies []entity.Comment
		cursor, err := ctr.commentModel.Base.Find(filter)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		err = cursor.All(context.Background(), &replies)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  replies,
		})
	}
}
