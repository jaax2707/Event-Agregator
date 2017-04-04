package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"fmt"
	"strconv"
	"bytes"
	"io"
)

var cas1 int
var cas2 int
var interval int
var counter int
var temp string

type Page struct {
	Title string
	Body []byte
}



func (p *Page) save () error {
	f := p.Title + ".txt"
	return ioutil.WriteFile(f, p.Body, 0600)
}

func load(title string) (*Page, error) {
	f := title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func view(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	p, _ := load(title)


	t, _ := template.ParseFiles("test.html")
	t.Execute(w, p)
	temp = title
	Counter(w,r)
}

func save(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	body := strconv.Itoa(counter)
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w,r,"/"+title, http.StatusFound)
}

func newTest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		cas1, _ = strconv.Atoi(r.URL.Query().Get("cas1"))
		cas2, _ = strconv.Atoi(r.URL.Query().Get("cas2"))
		interval, _ = strconv.Atoi(r.URL.Query().Get("interval"))
		//fmt.Println(cas1)
		//fmt.Println(cas2)
		//fmt.Println(interval)
		title := r.URL.Path[len("/"):]
		body := strconv.Itoa(counter)
		if temp != title{
			//counter = 0
			f := title + ".txt"
			body, _ := ioutil.ReadFile(f)
			counter, _ = strconv.Atoi(string(body))
			//buf := bytes.NewBuffer(body)
			//counter, _ = int(binary.ReadVarint(buf))
		}
		p := &Page{Title: title, Body: []byte(body)}
		p.save()

		view(w, r)

	case "POST":
		save(w, r)

	}
}
func Counter(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": // increment n
		counter++
	case "POST": // set n to posted value
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		body := buf.String()
		if _, err := strconv.Atoi(body); err != nil {
			fmt.Fprintf(w, "bad POST: %v\nbody: [%v]\n", err, body)
		} else {
			counter=0
			fmt.Fprint(w, "counter reset\n")
		}
	}
	fmt.Fprintf(w, "counter = %d\n", counter)
}

func main() {
	http.HandleFunc("/", newTest)
	http.ListenAndServe(":8000", nil)
	fmt.Println(cas1)
	fmt.Println(cas2)
	fmt.Println(interval)
}