package web

import (
	"analiser/pkg/lib"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"
)

// Routing (using gorilla/mux)
// простейший пример https://gowebexamples.com/routes-using-gorilla-mux/
// документация https://github.com/gorilla/mux?tab=readme-ov-file

type application struct {
}

func Process() {

	app := &application{}

	router := mux.NewRouter()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/table", app.table).Methods("GET")
	router.HandleFunc("/weeks", app.weeks)
	router.HandleFunc("/days", app.days)
	router.HandleFunc("/total", app.total)
	router.HandleFunc("/category/{name}", app.category)
	router.HandleFunc("/manage/", app.manage)
	router.HandleFunc("/day/{date}", app.day)
	router.Use(loggingMiddleware)

	handler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	router.PathPrefix("/static/").Handler(handler)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	shutdown(router)
	//err := http.ListenAndServe(":4000", router)
	//log.Fatal(err)
}

var startTime = time.Now()

// loggingMiddleware https://github.com/gorilla/mux?tab=readme-ov-file#middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime = time.Now()
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			filename := "data/log.txt"
			content := time.Now().Format("2006-01-02 15:04:05") + " " + r.RequestURI + "\n"
			lib.AppendToFile(filename, content)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func (app *application) display(tpl string, w http.ResponseWriter, data map[string]any) {

	var funcMap = template.FuncMap{
		"lower":  strings.ToLower,
		"repeat": func(s string) string { return strings.Repeat(s, 2) },
	}

	files := []string{
		"./templates/" + tpl,
		"./templates/base.layout.html",
	}
	ts, err := template.New(tpl).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if data == nil {
		data = map[string]any{}
	}
	data["categories"] = lib.Categories()
	data["days"] = lib.Last7days()

	err = ts.Execute(w, data)

	fmt.Printf("tpl %v, time: %v\n", tpl, time.Since(startTime))

	if err != nil {
		app.serverError(w, err)
	}
}

// shutdown Graceful
// Go 1.8 introduced the ability to gracefully shutdown a *http.Server. Here's how to do that alongside mux:
func shutdown(router *mux.Router) {
	//var wait time.Duration
	//flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	//flag.Parse()
	wait := time.Second * 15

	srv := &http.Server{
		Addr: ":4000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
