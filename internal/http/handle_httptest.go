package http

import (
	"context"
	"fmt"
	"html"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/yheip/network-tester/internal/http/sse"
)

func HandleHttpTest(run curlcmd) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		url := c.FormValue("url")
		log := zerolog.Ctx(c.Context())
		log.Info().Msgf("url: %s", url)

		sid := c.Cookies("sid")

		session := sessions.sessions[sid]

		if run != nil {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				response, info, err := run(ctx, url)
				if err != nil {
					log.Error().Err(err).Msg("curl error")
					return
				}
				var wg sync.WaitGroup
				wg.Add(2)

				go func() {
					defer wg.Done()
					for l := range response {
						data := fmt.Sprintf(`<div>%s</div>`, html.EscapeString(l))
						session.eventChannel <- sse.NewEvent("httptest-response", data)
					}
				}()

				go func() {
					defer wg.Done()
					for l := range info {
						data := fmt.Sprintf(`<div>%s</div>`, l)
						session.eventChannel <- sse.NewEvent("httptest-info", data)
					}
				}()

				wg.Wait()
			}()
		}

		return c.Render("views/httptest_output", fiber.Map{})
	}
}

type curlcmd func(ctx context.Context, url string) (<-chan string, <-chan string, error)
