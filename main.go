package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"fmt"
	"strconv"
	"os"
	"time"
	"strings"
)

var cas1 int
var cas2 int
var interval int
var counter int
var temp string
var telo string
var array []string

type Page struct {
	Title string
	Body []byte
}
type Page2 struct {
	Title2 string
	Body2 []string
}

func (p *Page) save () error {
	f := p.Title + ".txt"
	return ioutil.WriteFile(f, p.Body, 0600)
}

func load(title string) (*Page2, error) {
	return &Page2{Title2: title, Body2: array}, nil
}

func load2(title string) (*Page, error) {
	return &Page{Title: title, Body: []byte(strconv.Itoa(len(array)))}, nil
}

func view(w http.ResponseWriter, r *http.Request) {

	if cas1 != 0 || cas2 != 0 || interval != 0{
		title := r.URL.Path[len("/"):]
		p, _ := load(title)
		t, _ := template.ParseFiles("test.html")
		t.Execute(w, p)
		temp = title
	}else{
		title := r.URL.Path[len("/"):]
		p, _ := load2(title)
		t, _ := template.ParseFiles("test2.html")
		t.Execute(w, p)
		temp = title
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	body := strconv.Itoa(counter)
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w,r,"/"+title, http.StatusFound)
}

func Exists(name string) bool {
	if _, err := os.Stat(name+".txt"); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func filter(s string, title string) {
	if cas1 != 0 || cas2 != 0 || interval != 0{
	diff_time := cas2 - cas1
	pocet_Intervalov := diff_time / interval / 60

	for i := 0; i < pocet_Intervalov; i++ {

		init1 := time.Unix(1491400800+int64(i*10*60), 0).Format(time.RFC822)
		init2 := time.Unix(1491400800+int64(10*60*(i+1)), 0).Format(time.RFC822)

		str := strings.Split(s, ";")
		for j := 0; j < len(str); j++ {
			start, _ := time.Parse(time.RFC822, init1)
			end, _ := time.Parse(time.RFC822, init2)
			in, _ := time.Parse(time.RFC822, str[j])

			if inTimeSpan(start, end, in) {
				counter++
			}
		}
		array = append(array,"{   " + title + ":   " + strconv.Itoa(counter) + "   }")
		counter = 0
		}
	}else{
		array = strings.Split(s,";")
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func newTest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[len("/"):] != "favicon.ico" {
		switch r.Method {
		case "GET":
			cas1, _ = strconv.Atoi(r.URL.Query().Get("cas1"))
			cas2, _ = strconv.Atoi(r.URL.Query().Get("cas2"))
			interval, _ = strconv.Atoi(r.URL.Query().Get("interval"))

			title := r.URL.Path[len("/"):]

			if Exists(title) == true {

				if temp != title {
					f := title + ".txt"
					body, _ := ioutil.ReadFile(f)
					t := time.Now().Format("02 Jan 06 15:04 UTC")
					telo = string(body) + t + ";"
					p := &Page{Title: title, Body: []byte(telo)}
					p.save()

				}else{
					t := time.Now().Format("02 Jan 06 15:04 UTC")
					telo = telo + t + ";"
					p := &Page{Title: title, Body: []byte(telo)}
					p.save()
				}
				filter(telo, title)
				view(w, r)
			}
			array = nil
		case "POST":
			save(w, r)
		}
	}
}

func main() {
	http.HandleFunc("/", newTest)
	http.ListenAndServe(":8000", nil)
	fmt.Println(cas1)
	fmt.Println(cas2)
	fmt.Println(interval)
}