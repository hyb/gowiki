package main

import (
	"http"
	"io/ioutil"
	"os"
	"url"
)

const (
	viewRoot     = "views/"
	lenPath      = len("/view/")
)

// Page object for storing the title and body of a wiki page
type Page struct {
    Title string
    Body  []byte
}

func (p *Page) save() os.Error {
	filename := viewRoot + url.QueryEscape(p.Title) + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, os.Error) {
	filename := viewRoot + url.QueryEscape(title) + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, os.Error) {
	title := r.URL.Path[lenPath:]
	return title, nil
}


var viewCache = make(map[string] []byte)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, _ := getTitle(w, r)
	cached := viewCache[title]
	var respBytes []byte = nil


	// Load the cached page if it was already looked up
	if cached != nil {
		respBytes = cached
	} else {
		common := &CommonData{}
		p, err := loadPage(title)
		if err != nil {
			p = &Page{Title: title}
			respBytes, err = RenderPage("wiki_notfound", common, p)
		} else {
			respBytes, err = RenderPage("wiki_view", common, p)
		}

		// Cache the page for the next request
		viewCache[title] = respBytes
	}

	w.Write(respBytes)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, _ := getTitle(w, r)
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	respBytes, err := RenderPage("wiki_edit", &CommonData{}, p)
	w.Write(respBytes)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	viewCache[title] = nil
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
