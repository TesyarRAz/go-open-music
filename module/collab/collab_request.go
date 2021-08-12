package collab

type storeRequest struct {
	PlaylistID string `json:"playlistId" binding:"required"`
	UserID     string `json:"userId" binding:"required"`
}

type destroyRequest struct {
	PlaylistID string `json:"playlistId" binding:"required"`
	UserID     string `json:"userId" binding:"required"`
}
