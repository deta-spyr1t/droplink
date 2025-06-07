package model

import "time"

type FileMeta struct {
	ID           string
	Filename     string
	Path         string
	ExpiresAt    time.Time
	MaxDownloads int
	Downloads    int
}

var Files = make(map[string]*FileMeta)
