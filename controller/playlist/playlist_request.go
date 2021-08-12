package playlist

type storeSongRequest struct {
	SongId string `json:"songId" binding:"required"`
}

type destroySongRequest struct {
	SongId string `json:"songId" binding:"required"`
}
