package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"read-adviser/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved page")

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (hash string, err error) {
	defer func() { err = e.WrapIfErr("can't calculate hash", err) }()

	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", err
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
