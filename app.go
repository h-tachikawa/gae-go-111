package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.HandleFunc("/", indexHandler)

	r.Route("/articles", func(r chi.Router) {
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(ArticleCtx)
			r.Get("/", getArticle)
		})
	})

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		ctx := context.WithValue(r.Context(), "article", articleID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	article := ctx.Value("article")
	fmt.Fprint(w, "result: ", article)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World")
}