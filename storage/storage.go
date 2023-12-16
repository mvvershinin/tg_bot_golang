package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	"tg_bot_golang/lib/e"
)

type Storage interface {
	Save(p *Page) error
	Remove(p *Page) error
	IfExists(p *Page) (bool, error)
	PickRandom(userName string) (*Page, error)
}

type Page struct {
	URL      string
	UserName string
	// created_at time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("unable calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("unable calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
