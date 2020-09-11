package main

import (
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"

	"bytes"
	"log"
	"os"
	"time"
)

type ContactDetails struct {
	Email    string `json:"email"`
	Subject  string `json:"sublect"`
	Message  string `json:"message"`
	Username string `json:"message"`
}

type returningMessage struct {
	Result      int    `json:"result"`
	Token       string `json:"token"`
	Adress_list string `json:"adresses"`
	Message     string `json:"message"`
	Username    string `json:"username"`
}

func GetFileLogName() string {

	file_name := time.Now().Local().Format("2006_01_02")
	return "logs/info" + file_name + ".log"
}

var f *os.File

func Init() {
	// set location of log file

	fmt.Print("Start Init Log")

	logPath := GetFileLogName()

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// p.LogFile = f

	log.SetOutput(f)

}

func main() {
	Init()

	r := mux.NewRouter()
	r.HandleFunc("/testPost", testPost).Methods("POST")
	r.HandleFunc("/testPostError", postStatus0).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(NotFoundFuncHandler)
	r.Use(loggingMiddleware)

	s := &http.Server{
		Addr: "172.16.33.106:8080",
		//Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.Handler = r
	fmt.Printf("\nHTMLServer : Service started : Host=%v\n", s.Addr)
	log.Fatal(s.ListenAndServe())
	f.Close()
}

func testPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==testPost")
	//http.Error(w, err.Error("Not correct method"), http.StatusInternalServerError)

	details := returningMessage{
		Token:       r.FormValue("access-token"),
		Adress_list: r.FormValue("adress_list"),
		Username:    r.FormValue("username"),
		Message:     r.FormValue("message"),
		Result:      1,
	}

	// do something with details
	fmt.Println(details)

	js, err := json.Marshal(details)
	if err != nil {
		fmt.Println(js, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	//fmt.Println(string(js), err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	//fmt.Println(js)

}

func postStatus0(w http.ResponseWriter, r *http.Request) {
	fmt.Println("===postStatus0")
	http.Error(w, "Test Error", 0)
	return
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Mildware")
		// Do stuff here
		log.Println("RequestURI ", r.RequestURI)
		log.Println("Method ", r.Method)
		log.Println("Proto", r.Proto)
		log.Println("Header", r.Header)
		log.Println("Host", r.Host)
		log.Println("RemoteAddr ", r.RemoteAddr)
		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			log.Print("bodyErr ", bodyErr.Error())
			http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
			return
		}

		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		log.Printf("BODY: %q", rdr1)
		r.Body = rdr2

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func NotFoundFuncHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NotFoundFuncHandler")
	fmt.Println("Response", r)

	//  http.Error(w, "my own error message", http.StatusForbidden)

	// or using the default message error

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
