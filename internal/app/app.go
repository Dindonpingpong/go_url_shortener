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

type ShortedUrl string

type Url string

type Db struct {
	Storage map[ShortedUrl]Url
}

func New() *Db {
	db := &Db{Storage: make(map[ShortedUrl]Url)}

	return db
}

func (db *Db) SaveShortedUrl(sUrl ShortedUrl, url Url) {
	db.Storage[sUrl] = url
}

func (db *Db) GetUrl(sUrl ShortedUrl) (Url, error) {
	url, ok := db.Storage[sUrl]

	if !ok {
		return "", fmt.Errorf("Url not found")
	}

	return url, nil
}

type App struct {
	Config *Config
	Db *Db
}


func (a *App) Start() {
	http.HandleFunc("/", a.root)

	http.ListenAndServe(a.Config.Address, nil)
}

func (a *App) root(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
		case http.MethodGet:

			sUrl := ShortedUrl(path.Base(r.URL.Path))

			url, error := a.Db.GetUrl(sUrl)

			if error != nil {
				http.Error(w, error.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Location", string(url))

			w.WriteHeader(http.StatusTemporaryRedirect)
		case http.MethodPost:
			b, _ := ioutil.ReadAll(r.Body)
			url := Url(string(b))
			sUrl := generateShortUrl(url)

			a.Db.SaveShortedUrl(sUrl, url)

			w.WriteHeader(http.StatusCreated)
			
			fmt.Fprintf(w, "http://%s/%s", r.Host, sUrl)
		default:
			http.Error(w, "Only GET and POST requests are allowed!", http.StatusMethodNotAllowed)
	}
}

func generateShortUrl(url Url) ShortedUrl {
	h := md5.New()

	h.Write([]byte(url))

	sum := hex.EncodeToString(h.Sum(nil))	

	sUrl := ShortedUrl(sum)

	return sUrl
}