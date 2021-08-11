package playlist

type StoreSongRequest struct {
	SongId string `json:"songId" binding:"required"`
}
