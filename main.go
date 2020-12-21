package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/spf13/viper"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()

	// .env
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading .env")
	}
	gomniSecurityKey := viper.GetString("SECURITY_KEY")
	gomniClientID := viper.GetString("CLIENT_ID")
	gomniSecret := viper.GetString("SECRET")

	// setup gomniauth
	gomniauth.SetSecurityKey(gomniSecurityKey)
	gomniauth.WithProviders(github.New(gomniClientID, gomniSecret,
		"http://localhost:8080/auth/callback/github"))

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)

	staticFiles := http.FileServer(http.Dir("./templates/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", staticFiles))

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting the webserver on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
