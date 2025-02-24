package link

import (
	"gorm.io/gorm"
	"math/rand"
	"short-link/internal/stat"
)

type Link struct {
	gorm.Model
	URL  string      `json:"url"`
	Hash string      `json:"hash" gorm:"uniqueIndex"`
	Stat []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		URL: url,
	}
	link.generateHash()
	return link
}

func (link *Link) generateHash() {
	link.Hash = RandStringRunes(6)
}

var letterRunes = []rune("QWERTYUIOPASDFGHJKLZXCVBNM1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
