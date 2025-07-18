package mysse

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	URL         string
	LastEventID string
	Retry       time.Duration
	OnEvent     func(e *Event)
	OnError     func(err error)
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewClient(url string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		URL:    url,
		Retry:  3 * time.Second,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Client) Start() {
	go c.listen()
}

func (c *Client) Stop() {
	c.cancel()
}

func (c *Client) listen() {
	for {
		if err := c.connect(); err != nil {
			if c.OnError != nil {
				c.OnError(err)
			}
		}
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(c.Retry):
		}
	}
}

func (c *Client) connect() error {
	req, err := http.NewRequestWithContext(c.ctx, "GET", c.URL, nil)
	if err != nil {
		return err
	}
	if c.LastEventID != "" {
		req.Header.Set("Last-Event-ID", c.LastEventID)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("non-200 response")
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var event Event
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if event.Data != "" && c.OnEvent != nil {
				c.OnEvent(&event)
				event = Event{}
			}
			continue
		}
		if strings.HasPrefix(line, "id:") {
			event.ID = strings.TrimSpace(line[3:])
			c.LastEventID = event.ID
		} else if strings.HasPrefix(line, "event:") {
			event.Event = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(line, "data:") {
			if event.Data != "" {
				event.Data += "\n"
			}
			event.Data += strings.TrimSpace(line[5:])
		}
	}
}
