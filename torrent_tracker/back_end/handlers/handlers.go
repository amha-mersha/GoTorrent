package handlers

import (
	"github.com/amha-mersha/GoTorrent/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

type announceRequest struct {
}

func HandleAnnouncePeers(c *gin.Context) {
}

func HandleGetTorrentList(c *gin.Context) {
}

func HandleGetTorrent(c *gin.Context) {

}

func HandleTorrentPost(c *gin.Context) {

}
