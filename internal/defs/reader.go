package defs

import "github.com/SoundarTech/mediamtx/internal/logger"

// Reader is an entity that can read a stream.
type Reader interface {
	logger.Writer
	Close()
	APIReaderDescribe() *APIPathReader
}
