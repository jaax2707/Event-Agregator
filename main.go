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
var array []int

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

func filter(s string) {
	diff_time := cas2 - cas1
	pocet_Intervalov := diff_time / interval / 60
	//t := time.Unix(1491400800, 0).Format(time.RFC822)
	//e := time.Unix(1491408000, 0).Format(time.RFC822)

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
		array = append(array, counter)
		fmt.Print(array)
		counter = 0
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func newTest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[len("/"):] != "favicon.ico" {
		switch r.Method {
		case "GET":
			//Counter(w,r)
			cas1, _ = strconv.Atoi(r.URL.Query().Get("cas1"))
			cas2, _ = strconv.Atoi(r.URL.Query().Get("cas2"))
			interval, _ = strconv.Atoi(r.URL.Query().Get("interval"))
			//fmt.Println(cas1)
			//fmt.Println(cas2)
			//fmt.Println(interval)
			title := r.URL.Path[len("/"):]

			if Exists(title) == true {

				if temp != title {
					//counter = 0
					f := title + ".txt"
					body, _ := ioutil.ReadFile(f)
					//buf := bytes.NewBuffer(body)
					//counter, _ = int(binary.ReadVarint(buf))
					t := time.Now().Format("02 Jan 06 15:04 UTC")
					telo = string(body) + t + ";"
					//filter(telo)
					p := &Page{Title: title, Body: []byte(telo)}
					p.save()

				}else{
					t := time.Now().Format("02 Jan 06 15:04 UTC")
					telo = telo + t + ";"
					//filter(telo)
					p := &Page{Title: title, Body: []byte(telo)}
					p.save()
				}
				view(w, r)
			}
		case "POST":
			save(w, r)
		}
	}
}
//func Counter(w http.ResponseWriter, req *http.Request) {
//	switch req.Method {
//	case "GET": // increment n
//		counter++
//	case "POST": // set n to posted value
//		buf := new(bytes.Buffer)
//		io.Copy(buf, req.Body)
//		body := buf.String()
//		if _, err := strconv.Atoi(body); err != nil {
//			fmt.Fprintf(w, "bad POST: %v\nbody: [%v]\n", err, body)
//		} else {
//			counter=0
//			fmt.Fprint(w, "counter reset\n")
//		}
//	}
//	//fmt.Fprintf(w, "counter = %d\n", counter)
//}

func main() {
	http.HandleFunc("/", newTest)
	http.ListenAndServe(":8000", nil)
	fmt.Println(cas1)
	fmt.Println(cas2)
	fmt.Println(interval)
}