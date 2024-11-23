package handlers

import (
	"net/http"

	"github.com/amha-mersha/GoTorrent/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

type announceRequest struct {
	InfoHash   string `form:"info_hash" binding:"required"`
	PeerID     string `form:"peer_id" binding:"required"`
	IP         string `form:"ip"`
	Port       int    `form:"port" binding:"required"`
	Uploaded   int    `form:"uploaded"`
	Downloaded int    `form:"downloaded"`
	Left       int    `form:"left"`
}

func NewHandlers(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleAnnouncePeers(c *gin.Context) {
	var peerRequest announceRequest
	if err := c.ShouldBindQuery(&peerRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	peers, err := h.service.Announce(
		peerRequest.InfoHash,
		peerRequest.PeerID,
		peerRequest.IP,
		peerRequest.Port,
		peerRequest.Uploaded,
		peerRequest.Downloaded,
		peerRequest.Left,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"peers": peers})
}

func (h *Handler) HandleGetTorrentList(c *gin.Context) {
}

func (h *Handler) HandleGetTorrent(c *gin.Context) {

}

func (h *Handler) HandleTorrentPost(c *gin.Context) {

}
