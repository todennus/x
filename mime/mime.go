package mime

const (
	ApplicationJSON           = "application/json"
	ApplicationFormURLEncoded = "application/x-www-form-urlencoded"
	ApplicationOctetStream    = "application/octet-stream"
	ImageJPEG                 = "image/jpeg"
	ImagePNG                  = "image/png"
	MultipartFormData         = "multipart/form-data"
)

var imageMimeType = map[string]bool{
	ImageJPEG: true,
	ImagePNG:  true,
}

func IsImage(t ...string) bool {
	for i := range t {
		if !imageMimeType[t[i]] {
			return false
		}
	}

	return true
}
