package router

import (
	"Douyin/app/gateway/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	return InitRouter(r)
}

func InitRouter(r *gin.Engine) *gin.Engine {
	DouyinGroup := r.Group("/douyin")

	{
		// 用户模块
		UserGroup := DouyinGroup.Group("/user")
		{
			UserGroup.POST("/register", http.NewUserHandler().RegisterHandler)
			UserGroup.POST("/login", http.NewUserHandler().LoginHandler)
			UserGroup.GET("/", http.NewUserHandler().UserInfoHandler)
		}

		// 视频发布模块
		PublishGroup := DouyinGroup.Group("/publish")
		{
			// TODO 这个可以完成基本功能，但是我估计并不是真正的视频上传。毕竟目前的代码向mq上传视频文件太大，不可能
			PublishGroup.POST("/action/", http.NewPublishHandler().PublishHandler)
			// 发布列表： 返回发布所有视频的所有信息
			PublishGroup.GET("/list/", http.NewPublishHandler().PublishListHandler)
		}

		// 社交模块
		RelationGroup := DouyinGroup.Group("/relation")
		{
			// 关注操作
			RelationGroup.POST("/action/", http.NewRelationHandler().ActionRelationHandler)
			// 关注列表
			RelationGroup.GET("/follow/list/", http.NewRelationHandler().ListFollowRelationHandler)
			// 粉丝列表
			RelationGroup.GET("/follower/list/", http.NewRelationHandler().ListFollowerRelationHandler)
			// 好友列表
			RelationGroup.GET("/friend/list/", http.NewRelationHandler().ListFriendRelationHandler)
		}

	}

	return r
}
