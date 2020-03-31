package endpoints

import (
	"net/http"
)

const STATIC_DIR = "./wwwroot/"

func ServeStaticContent() {
	http.Handle("/", Index())

}

func Index() http.Handler {
	return http.FileServer(http.Dir(STATIC_DIR))
}

func ServeEndpoints() {

}
