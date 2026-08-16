package server

import (
	"regexp"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
)

type RouteMatch struct {
	Endpoint   *domain.Endpoint
	PathParams map[string]string
}

type Router struct {
	routes []*compiledRoute
}

type compiledRoute struct {
	endpoint   *domain.Endpoint
	regex      *regexp.Regexp
	paramNames []string
}

func NewRouter(endpoints []*domain.Endpoint) *Router {
	r := &Router{
		routes: make([]*compiledRoute, 0, len(endpoints)),
	}

	for _, ep := range endpoints {
		r.addRoute(ep)
	}

	return r
}

func (r *Router) addRoute(ep *domain.Endpoint) {
	pattern := ep.Path
	var paramNames []string

	// Remplacer {param} par regex ([^/]+)
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)

	for _, m := range matches {
		if len(m) >= 2 {
			paramNames = append(paramNames, m[1])
		}
	}

	regexPattern := "^" + paramRegex.ReplaceAllString(pattern, `([^/]+)`) + "$"
	compiled, err := regexp.Compile(regexPattern)
	if err != nil {
		return
	}

	r.routes = append(r.routes, &compiledRoute{
		endpoint:   ep,
		regex:      compiled,
		paramNames: paramNames,
	})
}

// Match trouve l'endpoint correspondant à la méthode et au chemin
func (r *Router) Match(method domain.HTTPMethod, path string) (*RouteMatch, bool) {
	cleanPath := path
	if cleanPath != "/" {
		cleanPath = strings.TrimSuffix(cleanPath, "/")
	}

	for _, cr := range r.routes {
		if cr.endpoint.Method != method {
			continue
		}

		submatches := cr.regex.FindStringSubmatch(cleanPath)
		if submatches != nil {
			params := make(map[string]string)
			if len(submatches) > 1 {
				for i, name := range cr.paramNames {
					if i+1 < len(submatches) {
						params[name] = submatches[i+1]
					}
				}
			}

			return &RouteMatch{
				Endpoint:   cr.endpoint,
				PathParams: params,
			}, true
		}
	}

	return nil, false
}
