package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ctxKey represents the type of value for the context key.===================================
//WEB HANDLER
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

var registered = make(map[string]bool)

type App struct {
	mux      *httptreemux.ContextMux
	otmux    http.Handler
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {

	mux := httptreemux.NewContextMux()

	return &App{
		mux:      mux,
		otmux:    otelhttp.NewHandler(mux, "request"),
		shutdown: shutdown,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, path string, handler Handler) {
	a.handle(false, method, path, handler)
}

// handle performs the real work of applying boilerplate and framework code
// for a handler.
func (a *App) handle(debug bool, method string, path string, handler Handler) {
	if debug {
		// Track all the handlers that are being registered so we don't have
		// the same handlers registered twice to this singleton.
		if _, exists := registered[method+path]; exists {
			return
		}
		registered[method+path] = true
	}

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Start or expand a distributed trace.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		v := Values{
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	// Add this handler for the specified verb and route.
	if debug {
		f := func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == method:
				h(w, r)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		http.DefaultServeMux.HandleFunc("/debug"+path, f)
		return
	}
	a.mux.Handle(method, path, h)
}

//==============================================================
//ROUTE
func Routes(build string, shutdown chan os.Signal, log *log.Logger) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := NewApp(shutdown)

	sg := speedyGroup{
		speedy: New(log),
	}
	app.Handle(http.MethodGet, "/v1/getspeed", sg.Query)

	return app
}