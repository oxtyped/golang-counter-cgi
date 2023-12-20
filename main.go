package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"strconv"
)

func doPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form
	for k := range form {
		fmt.Fprintln(w, "Form", k+":", form.Get(k))
	}
	post := r.PostForm
	for k := range post {
		fmt.Fprintln(w, "PostForm", k+":", post.Get(k))
	}

}

// doGet implements a handler to return an incremented value of a counter
func doGet(w http.ResponseWriter, r *http.Request) {

	// open a file locally that contains the integer value and increment
	// accordingly
	f, err := os.OpenFile("counter.dat", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Seek(0, 0)

	var counter int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		counterValue := scanner.Text()
		counter, err = strconv.Atoi(counterValue)
		if err != nil {
			log.Fatal(err)
		}
		counter += 1
		_, err = fmt.Fprintln(w, counter)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
	err = f.Truncate(0)
	f.Seek(0, 0)
	_, err = fmt.Fprintf(f, strconv.Itoa(counter))
	if err != nil {
		log.Fatalf("error writing: %#v", err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")

	doGet(w, r)
}

func main() {
	err := cgi.Serve(http.HandlerFunc(handler))
	if err != nil {
		fmt.Println(err)
	}
}
