package cache

import (
	"crypto/md5"
	"crypto/sha256"
	"hash/crc32"
	"moria/config"
)

func hash(imageID string) string {
	hashType := config.GetCacheHash()
	switch hashType {
	case config.CacheHashSHA256:
		return sha256Hash(imageID)
	case config.CacheHashMD5:
		return md5Hash(imageID)
	case config.CacheHashCRC32:
		return crc32Hash(imageID)
	}

	return imageID
}

func sha256Hash(data string) string {
	s := sha256.New()
	s.Write([]byte(data))
	return string(s.Sum(nil))
}

func md5Hash(data string) string {
	s := md5.New()
	s.Write([]byte(data))
	return string(s.Sum(nil))
}

func crc32Hash(data string) string {
	s := crc32.NewIEEE()
	s.Write([]byte(data))
	return string(s.Sum(nil))
}
