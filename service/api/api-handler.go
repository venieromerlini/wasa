package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	/*
		WASA ROUTES
	*/
	rt.router.POST("/session", rt.doLogin)
	rt.router.GET("/home/:username", rt.getMyStream)
	rt.router.GET("/users", rt.findAllUser)
	rt.router.GET("/users/:username", rt.getUserProfile)
	rt.router.PUT("/users/:username", rt.setMyUserName)

	rt.router.POST("/bans", rt.banUser)
	rt.router.GET("/bans", rt.findAllBans)
	rt.router.DELETE("/bans/:banId", rt.unbanUser)

	rt.router.POST("/comments", rt.commentPhoto)
	rt.router.GET("/comments", rt.findAllComments)
	rt.router.DELETE("/comments/:commentId", rt.uncommentPhoto)

	rt.router.GET("/follows", rt.findAllFollows)
	rt.router.POST("/follows", rt.followUser)
	rt.router.DELETE("/follows/:followId", rt.unfollowUser)

	rt.router.GET("/likes", rt.findAllLikes)
	rt.router.POST("/likes", rt.likePhoto)
	rt.router.DELETE("/likes/:likeId", rt.unlikePhoto)

	rt.router.GET("/photos", rt.findUserPhotos)
	rt.router.POST("/photos", rt.uploadPhoto)
	rt.router.GET("/photos/:photoId", rt.findUserPhoto)
	rt.router.DELETE("/photos/:photoId", rt.deletePhoto)

	return rt.router
}
