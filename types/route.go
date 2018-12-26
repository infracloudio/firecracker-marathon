package types

type Upload struct {
	Language   string `json:language`
	UploadType string `json: uploadType`
	Code       string `json:code`
}
