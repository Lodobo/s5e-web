package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type PageData struct {
	Title string
}

func replaceText(i string) string {
	switch i {
	case "0":
		return ":)"
	case "1":
		return ":("
	case "2":
		return ":0"
	case "3":
		return ":|"
	case "4":
		return ":/"
	case "5":
		return ":D"
	case "6":
		return "xD"
	case "7":
		return "0_0"
	case "8":
		return "D:"
	case "9":
		return "-_-"
	case "10":
		return ":-)"
	case "11":
		return ">:)"
	case "12":
		return ">:("
	case "13":
		return "^-^"
	case "14":
		return ":^)"
	case "15":
		return "x_x"
	}
	return ":)"
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/components/header.html",
		"web/templates/components/footer.html",
		"web/templates/pages/home.html",
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("parsing templates: %v", err)
		return
	}

	data := PageData{
		Title: "Home",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getCount(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")
	fmt.Fprint(w, replaceText(count))
}

func getAbout(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/components/header.html",
		"web/templates/components/footer.html",
		"web/templates/pages/about.html",
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("parsing templates: %v", err)
		return
	}

	data := PageData{
		Title: "About",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/counter", getCount)
	mux.HandleFunc("/about", getAbout)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:         ":3333",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Println("Server starting on :3333")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server stopped")
}
