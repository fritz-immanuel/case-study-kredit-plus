package library

import (
	"math/rand"
	"time"
)

var (
	StrToDateFormat      = "2006-01-02"
	StrToTimestampFormat = "2006-01-02 15:04:05"
)

func UTCPlus7() time.Time {
	return time.Now().UTC().Add(time.Hour * time.Duration(7))
}

func Randomizer() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
