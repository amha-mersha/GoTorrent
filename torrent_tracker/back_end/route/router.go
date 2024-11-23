package route

import (
	"github.com/amha-mersha/GoTorrent/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *handlers.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/announce", handler.HandleAnnouncePeers)
	r.GET("/torrents", handler.HandleGetTorrentList)
	r.POST("/torrents", handler.HandleTorrentPost)
	r.GET("/torrents/:id", handler.HandleGetTorrent)
	return r
}
