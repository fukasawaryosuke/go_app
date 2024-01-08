package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"net/http"
	"html/template"
)

type Page struct {
	Title  string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	//読み書き用の権限でファイルを作成します。既に存在している場合は、元のファイルを切り捨てます（トランケート）。
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// 与えられた文字列をファイルに書き込みます。戻り値は、書き込んだバイト数と書き込み中に起きたエラーです。
	_, err = io.WriteString(file,string(p.Body))
	return err
}

func loadPage (title string)(*Page,error){
	fmt.Println("loadPage")

	filename := title + ".txt"
	fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("err")
		return nil, err
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]

	p, err := loadPage(title)
	if err != nil {
		fmt.Println("view handler err")
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		fmt.Println("edit handler err")
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")

	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		fmt.Println("save handler err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
