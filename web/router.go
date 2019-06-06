package web

import (
//	"github.com/lengsh/findme/user"
	"net/http"
	"fmt"
)

func Router() {
	//	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(utils.ServerConfig.WebDir+"doc"))))
	lfs_router_register()

	http.HandleFunc("/", lfsStat)
	//http.HandleFunc("/hello/", hello)
	http.HandleFunc("/lfs/login", login)
//	http.HandleFunc("/html/", hellohtml)
//	http.HandleFunc("/api/userinfo/v1/", user.JwtAuth(userinfo))
}

func getDefaultData(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	fmt.Println(r.Header)
		var data  = make(map[string] interface{}, 1)
		if agent, ok := r.Header["User-Agent"]; ok {
			data["agent"] = agent[0]
		}
		if refer, ok := r.Header["Referer"]; ok {
			data["referer"] = refer[0]
		}
	return data
}

