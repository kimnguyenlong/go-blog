# Introductions
#### This is a backend project written in Go which provides Rest APIs of a blog app. Some features in this app are: 
- regiter/login
- follow/unfollow a use, get a user
- create/get topics
- create/get/update/delete posts
- create/get/update/delete post's comments
- create/get comment's replies
#### Used in this project:
- Gin Framework ([https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin))
- MongoDB 
# How to use
After cloning this repo, you can run the system via the following ways:
#### With Docker
Run cmd `docker-compose build` and then `docker-compose up`, then the services will be available on localhost8080.
You can use the exported Postman file `blog.postman_collection.json` to leverage my api call examples.