package app

import (
	"net/http"
	"path"
	"strconv"

	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/router"
	"github.com/bendigiorgio/ikou/internal/app/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func StartServer(devMode bool) {
	intPort := utils.GlobalConfig.Port
	port := strconv.Itoa(intPort)

	serverUrl := "http://localhost:" + port
	utils.Logger.Sugar().Info("Starting server on port: ", serverUrl)

	basePath := utils.GlobalConfig.BasePath
	staticPath := utils.GlobalConfig.StaticPath
	useSrc := utils.GlobalConfig.UseSrc

	srcPath := basePath
	if useSrc {
		srcPath = path.Join(basePath, "src")
	}

	router.InitializeRouting(srcPath, devMode)

	r := mux.NewRouter()

	staticDir := http.FileServer(http.Dir(staticPath))
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
		apiRouteInfo, apiExists := router.ApiRouteMap[route]

		if exists {
			initialProps := react.PageProps{
				PageRoute: route,
			}

			entryInfo, entryExists := router.EntryRouteMap[route]
			if entryExists {
				initialProps.Data = entryInfo.HandlerFn(w, r, entryInfo.FilePath)
			}

			pageData, err := react.RenderPage(
				routeInfo.IsSSG,
				initialProps,
				routeInfo.PagePath,
			)
			if err != nil {
				utils.Logger.Error("Page not found", zap.Error(err))
				http.Error(w, "Page not found", http.StatusNotFound)
				return
			}
			err = pageData.Tmpl.Execute(w, pageData)
			if err != nil {
				utils.Logger.Error("Error executing template", zap.Error(err))
			}
			return
		}

		if apiExists {
			// check if the request method is allowed
			if r.Method != apiRouteInfo.Method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// if it exists call the api handler function
			apiRouteInfo.HandlerFn(w, r, apiRouteInfo.FilePath)
			return
		}

		utils.Logger.Error("Page not found", zap.String("route", route))
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	portString := ":" + port
	utils.Logger.Sugar().Fatal(http.ListenAndServe(portString, r))

}
