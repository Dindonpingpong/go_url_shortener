package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)


type Config struct {
	Address string
}

type ShortedURL string

type URL string

type DB struct {
	Storage map[ShortedURL]URL
}

func New() *DB {
	DB := &DB{Storage: make(map[ShortedURL]URL)}

	return DB
}

func (DB *DB) SaveShortedURL(sURL ShortedURL, URL URL) {
	DB.Storage[sURL] = URL
}

func (DB *DB) GetURL(sURL ShortedURL) (URL, error) {
	URL, ok := DB.Storage[sURL]

	if !ok {
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

type App struct {
	Config *Config
	DB *DB
}


func (a *App) Start() {
	http.HandleFunc("/", a.root)

	http.ListenAndServe(a.Config.Address, nil)
}

func (a *App) root(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
		case http.MethodGet:

			sURL := ShortedURL(path.Base(r.URL.Path))

			URL, error := a.DB.GetURL(sURL)

			if error != nil {
				http.Error(w, error.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Location", string(URL))

			w.WriteHeader(http.StatusTemporaryRedirect)
		case http.MethodPost:
			b, _ := ioutil.ReadAll(r.Body)
			URL := URL(string(b))
			sURL := generateShortURL(URL)

			a.DB.SaveShortedURL(sURL, URL)

			w.WriteHeader(http.StatusCreated)
			
			fmt.Fprintf(w, "http://%s/%s", r.Host, sURL)
		default:
			http.Error(w, "Only GET and POST requests are allowed!", http.StatusMethodNotAllowed)
	}
}

func generateShortURL(URL URL) ShortedURL {
	h := md5.New()

	h.Write([]byte(URL))

	sum := hex.EncodeToString(h.Sum(nil))	

	sURL := ShortedURL(sum)

	return sURL
}