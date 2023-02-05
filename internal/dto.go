package internal

type ShortenURLReq struct {
	URL string `json:"url" form:"url"`
}

type ShortenURLResp struct {
	Key string `json:"key" form:"key"`
}
