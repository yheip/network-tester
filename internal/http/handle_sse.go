package http

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"github.com/yheip/network-tester/internal/http/sse"
)

type session struct {
	eventChannel chan sse.Event
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]*session
}

var sessions = SessionStore{
	mu:       sync.Mutex{},
	sessions: map[string]*session{},
}

func (ss *SessionStore) addSession(sid string, s *session) {
	ss.mu.Lock()
	ss.sessions[sid] = s
	ss.mu.Unlock()
}

func (ss *SessionStore) removeSession(sid string) {
	ss.mu.Lock()
	delete(ss.sessions, sid)
	ss.mu.Unlock()
}

func HandleSSE(c fiber.Ctx) error {
	log := zerolog.Ctx(c.Context())
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	notify := c.Context().Done()
	sid := uuid.NewString()
	eventChan := make(chan sse.Event)

	s := session{
		eventChannel: eventChan,
	}

	log.Info().Str("sid", sid).Msg("Adding new SSE session")

	sessions.addSession(sid, &s)

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		keepaliveTicker := time.NewTicker(2 * time.Second)
		defer keepaliveTicker.Stop()

		go func() {
			eventChan <- sse.NewEvent("say_hello", fmt.Sprint("helloword"))
		}()

		for {
			select {
			case e := <-eventChan:
				sse.WriteEvent(w, e)
				w.Flush()
			case <-keepaliveTicker.C:
				sse.WriteKeepAlive(w)
				w.Flush()
			case <-notify:
				return
			}
		}
	}))

	c.Cookie(&fiber.Cookie{
		Name:  "sid",
		Value: sid,
	})

	return nil
}
