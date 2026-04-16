package api

const (
	// Maximum upload size of 32MB for blobs and files.
	MaxUploadSize = 32 << 20
	// Catch-all request body limit of 64MB enforced at middleware level.
	MaxBodySize = 64 << 20
)
