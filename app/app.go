package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/azzamt11/todoapp_backend/app/handler"
	"github.com/azzamt11/todoapp_backend/app/model"
	"github.com/azzamt11/todoapp_backend/config"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/projects", a.handleRequest(handler.GetAllProjects))
	a.Post("/projects", a.handleRequest(handler.CreateProject))
	a.Get("/projects/{id}", a.handleRequest(handler.GetProject))
	a.Put("/projects/{id}", a.handleRequest(handler.UpdateProject))
	a.Delete("/projects/{id}", a.handleRequest(handler.DeleteProject))
	a.Put("/projects/{id}/archive", a.handleRequest(handler.ArchiveProject))
	a.Delete("/projects/{id}/archive", a.handleRequest(handler.RestoreProject))

	// Routing for handling the tasks
	a.Get("/projects/{projectID:[0-9]+}/tasks", a.handleRequest(handler.GetAllTasks))
	a.Post("/projects/{projectID:[0-9]+}/tasks", a.handleRequest(handler.CreateTask))
	a.Get("/projects/{projectID:[0-9]+}/tasks/{id:[0-9]+}", a.handleRequest(handler.GetTask))
	a.Put("/projects/{projectID:[0-9]+}/tasks/{id:[0-9]+}", a.handleRequest(handler.UpdateTask))
	a.Delete("/projects/{projectID:[0-9]+}/tasks/{id:[0-9]+}", a.handleRequest(handler.DeleteTask))
	a.Put("/projects/{projectID:[0-9]+}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.CompleteTask))
	a.Delete("/projects/{projectID:[0-9]+}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.UndoTask))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}