package mime

const (
	ApplicationJSON           = "application/json"
	ApplicationFormURLEncoded = "application/x-www-form-urlencoded"
	ApplicationOctetStream    = "application/octet-stream"
	ImageJPEG                 = "image/jpeg"
	ImagePNG                  = "image/png"
	MultipartFormData         = "multipart/form-data"
)

var IsImage = map[string]bool{
	ImageJPEG: true,
	ImagePNG:  true,
}
