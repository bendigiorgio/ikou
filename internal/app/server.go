package app

import (
	"log"
	"net/http"
	"path"

	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/router"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func StartServer() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting server with mux router")
	router.InitializeRouting("frontend/src")

	r := mux.NewRouter()

	staticDir := http.FileServer(http.Dir("frontend/public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", staticDir))

	r.HandleFunc("/{route:.*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		route := path.Clean(vars["route"])

		if route == "/" || route == "." {
			route = "/"
		} else {
			route = "/" + route
		}

		routeInfo, exists := router.RouteMap[route]
		if !exists {
			logger.Error("Page not found", zap.String("route", route))
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		initialProps := react.PageProps{
			PageRoute: route,
		}

		pageData, err := react.RenderPage(
			routeInfo.IsSSG,
			"frontend/src/clientEntry.tsx",
			initialProps,
			routeInfo.PagePath,
		)

		if err != nil {
			logger.Error("Page not found", zap.Error(err))
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		err = pageData.Tmpl.Execute(w, pageData)
		if err != nil {
			logger.Error("Error executing template", zap.Error(err))
		}
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
