package handlers

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/announce", HandleAnnouncePeers)
	r.GET("/torrents", HandleGetTorrentList)
	r.POST("/torrents", HandleTorrentPost)
	r.GET("/torrents/:id", HandleGetTorrent)
	return r
}
