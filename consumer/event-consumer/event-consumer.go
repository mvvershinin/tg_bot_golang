package event_consumer

import (
	"log"
	"tg_bot_golang/events"
	"time"
)

const // todo env
APP_DEBUG = true

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			//todo retry if need
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}

	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	debug := APP_DEBUG
	// todo sync.WaitGroup parallers run

	for _, event := range events {
		if debug {
			log.Printf("got new event '%s'", event.Text)
		}

		if err := c.processor.Process(event); err != nil {
			if debug {
				log.Printf("got new event '%s'", event.Text)
			}

			continue
		}
	}

	return nil
}
