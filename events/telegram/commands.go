package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"tg_bot_golang/clients/telegram"
	"tg_bot_golang/lib/e"
	"tg_bot_golang/storage"
)

const (
	// todo env
	APP_DEBUG = true
	CmdRnd    = "/rnd"
	CmdStart  = "/start"
	CmdHelp   = "/help"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	debug := APP_DEBUG

	if debug {
		log.Printf("got new command '%s' from '%s'", text, username)
	}

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}
	// add page: http...
	switch text {
	case CmdRnd:
		return p.sendRandom(chatId, username)
	case CmdHelp:
		return p.sendHelp(chatId)
	case CmdStart:
		return p.sendHello(chatId)
	default:
		return p.sendHelp(chatId)
	}
}

func (p *Processor) savePage(chatId int, pageUrl string, username string) (err error) {
	defer func() { err = e.Wrap("cannot do command save page", err) }()

	sendMsg := NewMessageSender(chatId, p.tg)

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IfExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(MsgPageAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(MsgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = e.Wrap("cannot do command random page", err) }()

	sendMsg := NewMessageSender(chatId, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		err := sendMsg(MsgNoSavedPages)
		if err != nil {
			return err
		}
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatId int) error {
	sendMsg := NewMessageSender(chatId, p.tg)
	return sendMsg(MsgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	sendMsg := NewMessageSender(chatId, p.tg)
	return sendMsg(MsgHello)
}

func NewMessageSender(chatId int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatId, msg)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	// text starts from http...
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
