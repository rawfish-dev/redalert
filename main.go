package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"redalert/alert"
	"redalert/common"
	// "redalert/pinger"
	"redalert/searcher"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

type DashboardInfo struct {
	Servers []common.Server
}

func dashboardHandler(c *appCtx, w http.ResponseWriter, r *http.Request) {

	templateBox, err := rice.FindBox("templates")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	templateString, err := templateBox.String("dash.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tmplMessage, err := template.New("dash").Parse(templateString)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	var transformedServers []common.Server
	for _, individualCheckableServer := range c.service.servers {
		transformedServers = append(transformedServers, common.Server{
			Name:      individualCheckableServer.GetServerDetails().Name,
			LastEvent: individualCheckableServer.GetServerWatcher().LastEvent,
		})
	}
	info := &DashboardInfo{Servers: transformedServers}

	if err := tmplMessage.Execute(w, info); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

type appCtx struct {
	service *Service
}

type appHandler struct {
	*appCtx
	h func(*appCtx, http.ResponseWriter, *http.Request)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.h(ah.appCtx, w, r)
}

func getPort() string {
	if os.Getenv("RA_PORT") == "" {
		return "8888"
	}
	return os.Getenv("RA_PORT")
}

func main() {
	service := &Service{
		alerts: alert.RegisteredAlerts,
	}
	service.initialize()

	// Pick up all pingers
	// for _, individualPinger := range pinger.Online {
	// 	service.AddServer(individualPinger)
	// }

	// Pick up all searchers
	for _, individualSearcher := range searcher.Online {
		service.AddServer(individualSearcher)
	}

	service.Start()
	context := &appCtx{
		service: service,
	}

	go func() {
		box := rice.MustFindBox("static")
		fs := http.FileServer(box.HTTPBox())
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.Handle("/", appHandler{context, dashboardHandler})

		port := getPort()
		fmt.Println("Listening on port ", port, " ...")
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic(err)
		}
	}()

	service.wg.Wait()

}
