package config

import "log"

type CacheHashType string
type SourceType string
type SourceRemoteUploadAuthType string

const (
	CacheHashNone   CacheHashType = "none"
	CacheHashSHA256 CacheHashType = "sha256"
	CacheHashMD5    CacheHashType = "md5"
	CacheHashCRC32  CacheHashType = "crc32"

	SourceTypeLocal        SourceType = "local"
	SourceTypeRemoteFTP    SourceType = "remote_ftp"
	SourceTypeRemoteSFTP   SourceType = "remote_sftp"
	SourceTypeRemoteSSH    SourceType = "remote_ssh"
	SourceTypeRemoteUpload SourceType = "remote_upload"

	SourceRemoteUploadAuthNone   SourceRemoteUploadAuthType = "none"
	SourceRemoteUploadAuthCustom SourceRemoteUploadAuthType = "custom"
	SourceRemoteUploadAuthBasic  SourceRemoteUploadAuthType = "basic"
	SourceRemoteUploadAuthBearer SourceRemoteUploadAuthType = "bearer"
)

func IsCacheEnabled() bool {
	value, _ := getBoolEnv("CACHE_ENABLED")

	return value
}

func GetCacheExpiration() int {
	value, _ := getIntEnv("CACHE_EXPIRATION")

	return value
}

func GetCacheHash() CacheHashType {
	value, _ := getStringEnv("CACHE_HASH", nil)

	return CacheHashType(value)
}

func GetSourceType() SourceType {
	value, _ := getStringEnv("SOURCE_TYPE", nil)

	return SourceType(value)
}

func GetSourceLocalPath() string {
	value, _ := getStringEnv("SOURCE_LOCAL_PATH", nil)

	return value
}

func GetSourceRemoteHost() string {
	value, _ := getStringEnv("SOURCE_REMOTE_HOST", nil)

	return value
}

func GetSourceRemotePort() int {
	value, _ := getIntEnv("SOURCE_REMOTE_PORT")

	return value
}

func GetSourceRemoteUploadAuthType() SourceRemoteUploadAuthType {
	value, _ := getStringEnv("SOURCE_REMOTE_UPLOAD_AUTH_TYPE", nil)

	return SourceRemoteUploadAuthType(value)
}

func GetSourceRemoteUploadAuthValue() string {
	value, _ := getStringEnv("SOURCE_REMOTE_UPLOAD_AUTH_VALUE", nil)

	return value
}

func GetSourceRemoteUploadAuthUsername() string {
	value, _ := getStringEnv("SOURCE_REMOTE_UPLOAD_AUTH_USERNAME", nil)

	return value
}

func GetSourceRemoteUploadAuthPassword() string {
	value, _ := getStringEnv("SOURCE_REMOTE_UPLOAD_AUTH_PASSWORD", nil)

	return value
}

func GetSourceRemoteUploadAuthBearerToken() string {
	value, _ := getStringEnv("SOURCE_REMOTE_UPLOAD_AUTH_BEARER_TOKEN", nil)

	return value
}

func GetSourceRemoteFTPPath() string {
	value, _ := getStringEnv("SOURCE_REMOTE_FTP_PATH", nil)

	return value
}

func GetSourceRemoteFTPUser() string {
	value, _ := getStringEnv("SOURCE_REMOTE_FTP_USER", nil)

	return value
}

func GetSourceRemoteFTPPass() string {
	value, _ := getStringEnv("SOURCE_REMOTE_FTP_PASS", nil)

	return value
}

func IsSourceRemoteFTPCloseOnEnd() bool {
	value, _ := getBoolEnv("SOURCE_REMOTE_FTP_CLOSE_ON_END")

	return value
}

func GetSourceRemoteSSHAuth() string {
	value, err := getStringEnv("SOURCE_REMOTE_SSH_AUTH", []string{"password", "key"})
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func GetSourceRemoteSSHUser() string {
	value, _ := getStringEnv("SOURCE_REMOTE_SSH_USER", nil)

	return value
}

func GetSourceRemoteSSHPass() string {
	value, _ := getStringEnv("SOURCE_REMOTE_SSH_PASS", nil)

	return value
}

func GetSourceRemoteSSHKey() string {
	value, _ := getStringEnv("SOURCE_REMOTE_SSH_KEY", nil)

	return value
}

func GetSourceRemoteSSHKeyPass() string {
	value, _ := getStringEnv("SOURCE_REMOTE_SSH_KEY_PASS", nil)

	return value
}

func GetSourceRemoteSSHPath() string {
	value, _ := getStringEnv("SOURCE_REMOTE_SSH_PATH", nil)

	return value
}

func IsSourceRemoteSSHCloseOnEnd() bool {
	value, _ := getBoolEnv("SOURCE_REMOTE_SSH_CLOSE_ON_END")

	return value
}

func GetSourceRemoteUploadURL() string {
	value, _ := getStringEnv("SOURCE_REMOTE_UPLOAD_URL", nil)

	return value
}

func IsProcessingEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_ENABLED")

	return value
}

func IsProcessingCompressionEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_COMPRESSION_ENABLED")

	return value
}

func IsProcessingResizeEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_RESIZE_ENABLED")

	return value
}

func IsProcessingCropEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_CROP_ENABLED")

	return value
}

func IsProcessingRotateEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_ROTATE_ENABLED")

	return value
}

func IsProcessingToGrayscaleEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_TO_GRAYSCALE_ENABLED")

	return value
}

func IsProcessingBlurEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_BLUR_ENABLED")

	return value
}

func IsProcessingWatermarkEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_WATERMARK_ENABLED")

	return value
}

func IsProcessingFlipEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_FLIP_ENABLED")

	return value
}

func IsProcessingFlopEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_FLOP_ENABLED")

	return value
}

func IsProcessingZoomEnabled() bool {
	value, _ := getBoolEnv("PROCESSING_ZOOM_ENABLED")

	return value
}

func IsAPIEnabled() bool {
	value, _ := getBoolEnv("API_ENABLED")

	return value
}

func GetAPIUploadEndpoint() string {
	value, _ := getStringEnv("API_UPLOAD_ENDPOINT", nil)

	return value
}

func IsSecurityCORSEnabled() bool {
	value, _ := getBoolEnv("SECURITY_CORS_ENABLED")

	return value
}

func GetSecurityCORSOrigin() string {
	value, _ := getStringEnv("SECURITY_CORS_ORIGIN", nil)

	return value
}

func GetSecurityAPIAuthHeader() string {
	value, _ := getStringEnv("SECURITY_API_AUTH_HEADER", nil)

	return value
}

func GetSecurityAPIAuthToken() string {
	value, _ := getStringEnv("SECURITY_API_AUTH_TOKEN", nil)

	return value
}

func GetSecurityAllowedMimeTypes() []string {
	value, _ := getStringArrayEnv("SECURITY_ALLOWED_MIME_TYPES", ",")

	return value
}

func IsFrontendEnabled() bool {
	value, _ := getBoolEnv("FRONTEND_ENABLED")

	return value
}

func GetFrontendEndpoint() string {
	value, _ := getStringEnv("FRONTEND_ENDPOINT", nil)

	return value
}

func GetFrontendUsername() string {
	value, _ := getStringEnv("FRONTEND_USERNAME", nil)

	return value
}

func GetFrontendPassword() string {
	value, _ := getStringEnv("FRONTEND_PASSWORD", nil)

	return value
}
