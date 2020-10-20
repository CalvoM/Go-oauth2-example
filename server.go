package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"github.com/gorilla/mux"
)

var manager = manage.NewDefaultManager()
var clientStore = store.NewClientStore()

//ClientCred data to be sent to browser
type ClientCred struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

var srvToken = server.NewServer(server.NewConfig(), manager)

func main() {

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
	})
	manager.MapClientStorage(clientStore)
	srvToken.SetPasswordAuthorizationHandler(PasswordHandler)
	srvToken.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srvToken.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./UI"))
	r.PathPrefix("/UI/").Handler(http.StripPrefix("/UI/", fs))
	r.PathPrefix("/api/auth/").HandlerFunc(authHandler).Methods("POST")
	r.PathPrefix("/api/").HandlerFunc(apiHandler)
	r.PathPrefix("/").HandlerFunc(serveTemplate)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:3001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

//Handle templates
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	basePath := filepath.Join("UI", "templates", "base.html")
	reqPath := filepath.Join("UI", "templates", filepath.Clean(r.URL.Path))
	info, err := os.Stat(reqPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles(basePath, reqPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiHandler")
	fmt.Println(r)
}
func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/auth/create_client/" {
		fmt.Println(r.Context().Deadline())
		_ = r.ParseForm()
		res := ClientCred{}
		formData := map[string]map[string][]string{"form": r.PostForm}
		res.ID, res.Secret = getCredentials(formData["form"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	} else if r.URL.Path == "/api/auth/generate_token/" {
		err := srvToken.HandleTokenRequest(w, r)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if r.URL.Path == "/api/auth/test" {
		token, err := srvToken.ValidationBearerToken(r)
		if err != nil {
			log.Fatal(err.Error())
		}
		data := map[string]interface{}{
			"client_id": token.GetClientID(),
			"user_id":   token.GetUserID(),
		}
		json.NewEncoder(w).Encode(data)
	}
}
func getCredentials(data map[string][]string) (string, string) {
	clientID := generateSecurityCredentials(18)
	clientSecret := generateSecurityCredentials(36)
	clientStore := store.NewClientStore()
	err := clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: "http://localhost:3001",
		UserID: "Test",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	manager.MapClientStorage(clientStore)
	return clientID, clientSecret
}
func generateSecurityCredentials(size int) string {
	key := make([]byte, size)
	_, _ = rand.Read(key[:])
	return base64.URLEncoding.EncodeToString(key[:])
}

// PasswordHandler handle password authorization login for user
func PasswordHandler(username, password string) (userID string, err error) {
	userID = "Test"
	err = nil
	return
}
