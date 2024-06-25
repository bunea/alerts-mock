package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sync"
	"text/template"
	"time"

	"github.com/bunea/csfi-alerts-mock/images"
	"github.com/bunea/csfi-alerts-mock/models"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Server struct {
	lock     sync.Mutex
	renderer *TemplateRenderer
	feeds    map[string]*models.Feed // key is feed id
}

func NewServer() *Server {
	var initialFeeds = make(map[string]*models.Feed)
	f, err := os.ReadFile("initial_feeds.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(f, &initialFeeds)
	if err != nil {
		panic(err)
	}

	return &Server{
		lock:  sync.Mutex{},
		feeds: initialFeeds,
		renderer: &TemplateRenderer{
			templates: template.Must(template.ParseFiles("feed.xml", "alert.html")),
		},
	}
}

func (s *Server) ListFeeds(c echo.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	feeds := make([]*models.Feed, 0, len(s.feeds))
	for _, feed := range s.feeds {
		feeds = append(feeds, &models.Feed{
			ID:          feed.ID,
			Title:       feed.Title,
			UpdatedAt:   feed.UpdatedAt,
			Entries:     []*models.Entry{},
			UpdateEvery: feed.UpdateEvery,
		})
	}

	return c.JSON(http.StatusOK, feeds)
}

func (s *Server) GetFeed(c echo.Context) error {
	id := c.Param("id")

	s.lock.Lock()
	feed, ok := s.feeds[id]
	if !ok {
		return c.JSON(http.StatusNotFound, "Feed not found")
	}
	s.lock.Unlock()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextXML)
	return c.Render(http.StatusOK, "feed.xml", feed)
}

func (s *Server) AddEntry(c echo.Context) error {
	feedID := c.Param("id")
	entry := new(models.Entry)
	if err := c.Bind(entry); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	if entry.Href == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request. Fill 'href'")
	}

	s.lock.Lock()
	feed, ok := s.feeds[feedID]
	if !ok {
		return c.JSON(http.StatusNotFound, "Feed not Found")
	}

	now := time.Now().UTC()
	feed.UpdatedAt = now
	entry.PublishedAt = now
	entry.UpdatedAt = now

	entries := feed.Entries
	entries = append(entries, entry)
	slices.SortFunc(entries, func(a, b *models.Entry) int {
		if a.UpdatedAt.After(b.UpdatedAt) {
			return -1
		}
		return 0
	})
	feed.Entries = entries

	s.lock.Unlock()

	return c.JSON(http.StatusCreated, feed)
}

func (s *Server) GetAlert(c echo.Context) error {
	id := c.Param("id")
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, feed := range s.feeds {
		for _, entry := range feed.Entries {
			if entry.ID == id {
				return c.Render(http.StatusOK, "alert.html", entry)
			}
		}
	}

	return c.JSON(http.StatusNotFound, "Alert not found")
}

func (s *Server) updateFeeds() {
	t := time.NewTicker(10 * time.Second)
	for range t.C {
		now := time.Now().UTC()

		s.lock.Lock()
		for _, feed := range s.feeds {
			if feed.UpdateEvery == 0 {
				continue
			}

			if now.Sub(feed.UpdatedAt.UTC()) < time.Duration(feed.UpdateEvery) {
				continue
			}

			var entry models.Entry
			err := faker.FakeData(&entry)
			if err != nil {
				fmt.Println(err)
			}
			entry.UpdatedAt = now
			entry.PublishedAt = now
			entry.Href = fmt.Sprintf("http://localhost:1323/alerts/%s", entry.ID)
			entry.ImageURL = images.GetImageURL(entry.ID)
			feed.UpdatedAt = now

			if len(feed.Entries) > 25 {
				feed.Entries = slices.Delete(feed.Entries, len(feed.Entries)-2, len(feed.Entries)-1)
			}
			feed.Entries = append(feed.Entries, &entry)
		}
		s.lock.Unlock()
	}
}

func main() {
	s := NewServer()

	e := echo.New()
	e.Debug = true
	e.Renderer = s.renderer
	e.GET("/alerts/list", s.ListFeeds)
	e.GET("/alerts/feed/:id", s.GetFeed)
	e.GET("/alerts/:id", s.GetAlert)
	e.POST("/alerts/feed/:id/entries", s.AddEntry)

	go s.updateFeeds()

	e.Logger.Fatal(e.Start(":1323"))
}
