package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"tg_bot_golang/lib/e"
	"tg_bot_golang/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerms = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("unable to save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err = os.MkdirAll(fPath, defaultPerms); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.Wrap("cannot pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) Remove(p *storage.Page) error {
	filename, err := fileName(p)
	if err != nil {
		return e.Wrap("unable get filename", err)
	}

	path := filepath.Join(s.basePath, p.UserName, filename)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("unable remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IfExists(p *storage.Page) (bool, error) {
	filename, err := fileName(p)
	if err != nil {
		return false, e.Wrap("unable get filename", err)
	}

	path := filepath.Join(s.basePath, p.UserName, filename)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("unable check if file exists %s", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("cannot open file", err)
	}

	defer func() { _ = file.Close() }()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, e.Wrap("cannot decode page", err)
	}

	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
