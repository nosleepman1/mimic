package server

import (
	"math/rand"
	"net/http"
	"time"
)

type MiddlewareConfig struct {
	Latency   time.Duration
	Jitter    time.Duration
	ChaosRate float64 // 0.0 to 1.0 (ex: 0.05 = 5% 500 errors)
	EnableCORS bool
}

func ApplyMiddleware(next http.Handler, config MiddlewareConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. CORS
		if config.EnableCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// 2. Simulation de Latence & Jitter
		if config.Latency > 0 {
			sleepDuration := config.Latency
			if config.Jitter > 0 {
				jitterMax := int64(config.Jitter)
				randomJitter := time.Duration(rand.Int63n(jitterMax*2) - jitterMax)
				sleepDuration += randomJitter
				if sleepDuration < 0 {
					sleepDuration = 0
				}
			}
			time.Sleep(sleepDuration)
		}

		// 3. Simulation de Chaos (Erreur 500 aléatoire)
		if config.ChaosRate > 0 && rand.Float64() < config.ChaosRate {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"Internal Server Error","message":"Injected by Mimic Chaos Simulation Engine"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
