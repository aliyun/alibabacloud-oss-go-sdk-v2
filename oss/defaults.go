package oss

import "os"

const (
	MaxUploadParts int32 = 10000

	// Max part size, 5GB, For UploadPart
	MaxPartSize int64 = 5 * 1024 * 1024 * 1024

	// Min part size, 100KB, For UploadPart
	MinPartSize int64 = 100 * 1024

	// Default part size, 5M
	DefaultPartSize int64 = 5 * 1024 * 1024

	// Default part size for uploader uploads data
	DefaultUploadPartSize = DefaultPartSize

	// Default part size for downloader downloads object
	DefaultDownloadPartSize = DefaultPartSize

	// Default part size for copier copys object, 64M
	DefaultCopyPartSize int64 = 64 * 1024 * 1024

	// Default parallel
	DefaultParallel = 3

	// Default parallel for uploader uploads data
	DefaultUploadParallel = DefaultParallel

	// Default parallel for downloader downloads object
	DefaultDownloadParallel = DefaultParallel

	// Default parallel for copier copys object
	DefaultCopyParallel = DefaultParallel

	// Default prefetch threshold to swith to async read in ReadOnlyFile
	DefaultPrefetchThreshold int64 = 20 * 1024 * 1024

	// Default prefetch number for async read in ReadOnlyFile
	DefaultPrefetchNum = DefaultParallel

	// Default prefetch chunk size for async read in ReadOnlyFile
	DefaultPrefetchChunkSize = DefaultPartSize

	// Default threshold to use muitipart copy in Copier, 256M
	DefaultCopyThreshold int64 = 200 * 1024 * 1024

	// File permission
	FilePermMode = os.FileMode(0664)

	// Temp file suffix
	TempFileSuffix = ".temp"

	// Checkpoint file suffix for Downloader
	CheckpointFileSuffixDownloader = ".dcp"

	// Checkpoint file suffix for Uploader
	CheckpointFileSuffixUploader = ".ucp"

	// Checkpoint file Magic
	CheckpointMagic = "92611BED-89E2-46B6-89E5-72F273D4B0A3"

	// Product for signing
	DefaultProduct = "oss"

	// The URL's scheme, default is https
	DefaultEndpointScheme = "https"

	// Default signature version is v4
	DefaultSignatureVersion = SignatureVersionV4
)
