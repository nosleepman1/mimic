package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/teranga-host/mimic/internal/domain"
	"github.com/teranga-host/mimic/internal/store"
	"github.com/teranga-host/mimic/internal/ui"
)

type Config struct {
	Port       int
	Host       string
	Latency    time.Duration
	Jitter     time.Duration
	ChaosRate  float64
	EnableCORS bool
}

type Server struct {
	config    Config
	endpoints []*domain.Endpoint
	router    *Router
	store     *store.MemoryStore
	logger    *ui.Logger
	httpSrv   *http.Server
}

func NewServer(config Config, endpoints []*domain.Endpoint) *Server {
	if config.Port == 0 {
		config.Port = 8080
	}
	if config.Host == "" {
		config.Host = "0.0.0.0"
	}

	router := NewRouter(endpoints)
	memStore := store.NewMemoryStore()
	logger := ui.NewLogger()

	return &Server{
		config:    config,
		endpoints: endpoints,
		router:    router,
		store:     memStore,
		logger:    logger,
	}
}

func (s *Server) Start() error {
	handler := NewMockHandler(s.router, s.store, s.logger)

	middlewareCfg := MiddlewareConfig{
		Latency:    s.config.Latency,
		Jitter:     s.config.Jitter,
		ChaosRate:  s.config.ChaosRate,
		EnableCORS: s.config.EnableCORS,
	}

	wrappedHandler := ApplyMiddleware(handler, middlewareCfg)

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: wrappedHandler,
	}

	s.printServerInfo()

	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpSrv != nil {
		return s.httpSrv.Shutdown(ctx)
	}
	return nil
}

func (s *Server) printServerInfo() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	gold := color.New(color.FgHiYellow, color.Bold).SprintFunc()

	fmt.Println()
	fmt.Printf("%s Server is active and listening on %s\n\n", green("[RUNNING]"), cyan(fmt.Sprintf("http://localhost:%d", s.config.Port)))
	fmt.Printf("%s (%d loaded):\n", gold("Simulated Endpoints"), len(s.endpoints))

	for _, ep := range s.endpoints {
		methodColor := color.New(color.FgBlue, color.Bold).SprintfFunc()
		if ep.Method == "POST" {
			methodColor = color.New(color.FgGreen, color.Bold).SprintfFunc()
		} else if ep.Method == "DELETE" {
			methodColor = color.New(color.FgRed, color.Bold).SprintfFunc()
		} else if ep.Method == "PUT" || ep.Method == "PATCH" {
			methodColor = color.New(color.FgYellow, color.Bold).SprintfFunc()
		}

		crudNote := ""
		if ep.IsCRUD {
			crudNote = color.New(color.FgHiBlack).Sprint("[Stateful CRUD]")
		}

		fmt.Printf("  %-7s %-35s %s\n", methodColor(string(ep.Method)), ep.Path, crudNote)
	}

	fmt.Println()
	if s.config.Latency > 0 {
		fmt.Printf("  [CONFIG] Network Latency : %v (Jitter: %v)\n", s.config.Latency, s.config.Jitter)
	}
	if s.config.ChaosRate > 0 {
		fmt.Printf("  [CONFIG] Chaos Error Rate : %.1f%%\n", s.config.ChaosRate*100)
	}
	fmt.Println(color.New(color.FgHiBlack).Sprint("  ──────────────────────────────────────────────────────────────────────────"))
	fmt.Println()
}
