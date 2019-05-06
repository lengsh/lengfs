package web

import (
	"net/http"
)

func Router() {
	//	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(utils.ServerConfig.WebDir+"doc"))))
	lfs_router_register()

	http.HandleFunc("/", lfsStat)
	http.HandleFunc("/login", login)

}
