package main

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/egoodhall/go-links/pkg/config"
	"github.com/egoodhall/go-links/pkg/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/time/rate"
)

func main() {
	ktx := kong.Parse(new(Cli), kong.Configuration(kongyaml.Loader, "go-links.yaml", "go-links.yml"))
	ktx.FatalIfErrorf(ktx.Run())
}

type Cli struct {
	ConfigFile kong.ConfigFlag `name:"config" short:"c"`
	Targets    config.Targets  `name:"targets"`
}

func (cli *Cli) Run() error {
	srv := echo.New()
	srv.Use(
		middleware.RequestID(),
		middleware.RemoveTrailingSlash(),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Every(100*time.Millisecond))),
	)

	srv.HTTPErrorHandler = func(err error, c echo.Context) {
		slog.Info("Handling error", "error", err)
		srv.DefaultHTTPErrorHandler(err, c)
	}

	srv.Renderer = ui.Render

	srv.OnAddRouteHandler = func(host string, route echo.Route, handler echo.HandlerFunc, middleware []echo.MiddlewareFunc) {
		slog.Info("Adding Endpoint", "method", route.Method, "path", route.Path)
	}

	for _, target := range cli.Targets {
		target.ArgUrls(func(alias string, url config.ArgUrl) {
			srv.GET(url.GetPath(alias), func(c echo.Context) error {
				return c.Redirect(http.StatusTemporaryRedirect, url.Render(c.Param))
			})
		})
	}

	srv.StaticFS("/_/static/", echo.MustSubFS(ui.StaticFiles, "static"))
	srv.GET("/_/clear", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	goLinks := ui.NewGoLinks(cli.Targets)
	aliases := cli.Targets.Aliases()

	// Search filtering
	srv.GET("/_/search", func(c echo.Context) error {
		ranks := fuzzy.RankFindNormalizedFold(c.QueryParam("q"), aliases)
		seen := make(map[int]struct{})
		links := make([]ui.GoLink, 0)
		for _, rank := range ranks {
			idx := slices.IndexFunc(goLinks, func(l ui.GoLink) bool {
				return slices.Contains(l.Aliases, rank.Target)
			})

			if idx < 0 {
				continue
			}

			if _, alreadySeen := seen[idx]; !alreadySeen {
				seen[idx] = struct{}{}
				links = append(links, goLinks[idx])
			}
		}

		return c.Render(http.StatusOK, "links", links)
	})

	// Arg URLs form
	srv.GET("/_/args-form", func(c echo.Context) error {
		return c.Render(http.StatusOK, "args-form", ui.Link{
			From: c.QueryParam("from"),
			To:   c.QueryParam("to"),
		})
	})

	// Render URL with args
	srv.POST("/_/render-url", func(c echo.Context) error {
		argUrl := config.ArgUrl(c.QueryParam("url"))
		url := argUrl.Render(c.FormValue)
		slog.Info("Rendered URL", "url", url)
		c.Response().Header().Add("HX-Redirect", url)
		return c.NoContent(http.StatusOK)
	})

	// Not found
	srv.RouteNotFound("/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", ui.Data{
			Query:   c.QueryParam("q"),
			GoLinks: goLinks,
		})
	})

	return srv.Start(":8080")
}
