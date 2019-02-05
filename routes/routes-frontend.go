package routes

import (
  "fmt"
  "os"
  "log"
  "net/http"
  "html/template"
  "path/filepath"

  "github.com/gorilla/mux"
)

var frontend_routes map[string]string

func FrontEndRoutes(r *mux.Router) {
  frontend_routes = map[string]string{
    "/"         : "index.html",
    "/index"    : "index.html",
    "/login"    : "auth/login.html",
    "/signup"   : "auth/signup.html",
    "/404"      : "error/404.html",
    "/500"      : "error/500.html",
  }
  r.PathPrefix("/").HandlerFunc(serveTemplate)
}

func serveTemplate(res http.ResponseWriter, req *http.Request) {
  fmt.Println("[", req.Method, "] frontend url", req.URL.Path)
  lp := filepath.Join("templates", "layout.html")
  fp := filepath.Join(
    "templates",
    filepath.Clean(
      frontend_routes[req.URL.Path],
    ),
  )

  // Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
  if err != nil {
    fmt.Println(err.Error())
    if os.IsNotExist(err) {
      fp = filepath.Join(
        "templates",
        filepath.Clean(
          frontend_routes["/404"],
        ),
      )
    }
  }

  if info.IsDir() {
    fmt.Println("[", req.Method, "] frontend url", req.URL.Path, "is not found")
    fp = filepath.Join(
      "templates",
      filepath.Clean(
        frontend_routes["/404"],
      ),
    )
  }

  tmpl, err := template.ParseFiles(lp, fp)
  if err != nil {
    log.Println(err.Error())
    fp = filepath.Join(
      "templates",
      filepath.Clean(
        frontend_routes["/500"],
      ),
    )
  }

  if err := tmpl.ExecuteTemplate(res, "layout", nil); err != nil {
    log.Println(err.Error())
    http.Error(res, http.StatusText(500), 500)
  }
}
