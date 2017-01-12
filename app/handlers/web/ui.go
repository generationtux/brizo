package web

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/CloudyKit/jet"
	"github.com/generationtux/brizo/database"
)

var views = jet.NewHTMLSet("./ui")

// RootHandler redirects to Javascript app
func RootHandler(rw http.ResponseWriter, request *http.Request) {
	http.Redirect(rw, request, "/app", 301)
}

// UIHandler for requests to Javascript app
func UIHandler(rw http.ResponseWriter, request *http.Request) {
	view, err := views.GetTemplate("index.html")
	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}

	assets, err := getUIAssets()
	if err != nil {
		log.Println("Unable to get UI assets: ", err)
	}

	vars := make(jet.VarMap)
	vars.Set("assets", assets)
	view.Execute(rw, vars, nil)
}

// data structure for UI assets.json file
type uiAssets struct {
	App       map[string]string
	Vendor    map[string]string
	Polyfills map[string]string
}

// getUIAssets retrieves the current version of UI asset paths
func getUIAssets() (uiAssets, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return uiAssets{}, errors.New("Unable to find current file path.")
	}
	assetPath := path.Join(path.Dir(file), "../../../ui/dist/assets.json")

	var assets uiAssets
	assetContents, err := ioutil.ReadFile(assetPath)
	if err != nil {
		return uiAssets{}, err
	}

	err = json.Unmarshal(assetContents, &assets)
	if err != nil {
		return uiAssets{}, err
	}

	return assets, nil
}

// HealthzHandler for health check requests
func HealthzHandler(rw http.ResponseWriter, request *http.Request) {
	err := database.Health()

	if err != nil {
		log.Println(err)
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
		return
	}

	rw.Write([]byte("ok"))
}
