package router

import (
	"Douyin/app/gateway/http"
	"Douyin/app/gateway/middleware"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	return InitRouter(r)
}

func InitRouter(r *gin.Engine) *gin.Engine {

	// 中间件: jwt 鉴权、跨域
	r.Use(
		middleware.JWT(),
		cors.Default(),
	)

	// 抖音接口
	DouyinGroup := r.Group("/douyin")

	DouyinGroup.GET("feed", http.NewPublishHandler().FeedHandler)

	{
		// 用户模块 ok
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

		// 社交模块 ok
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

		// 点赞模块
		FavoriteGroup := DouyinGroup.Group("/favorite")
		{
			FavoriteGroup.POST("/action/", http.NewFavoriteHandler().FavoriteActionHandler)
			FavoriteGroup.GET("/list/", http.NewFavoriteHandler().FavoriteListHandler)
		}

		// 评论模块
		CommentGroup := DouyinGroup.Group("/comment")
		{
			CommentGroup.POST("/action/", http.NewCommentHandler().CommentActionHandler)
			CommentGroup.GET("/list/", http.NewCommentHandler().CommentListHandler)
		}

		// 消息模块
		MessageGroup := DouyinGroup.Group("/message")
		{
			MessageGroup.POST("/action/", http.NewMessageHandler().ActionMessageHandler)
			MessageGroup.GET("/chat/", http.NewMessageHandler().ChatMessageHandler)
		}

	}

	return r
}
