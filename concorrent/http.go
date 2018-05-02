package main

import(
"fmt"
"flag"
"strings"
"log"
"net/http"
)

var (
	maxworkers   = flag.Int("workers", 50, "max workers.")
	maxTasks = flag.Int("tasks", 1000, "max tasks.")
	listenAddr        = flag.String("l", ":3000", "Listen address.")
)

type UserData struct {
	Ip string
}

func (this *UserData) Send() error {
	//fmt.Println(this.Ip)
	statis.Success++
	return nil
}

var dispatcher *Dispatcher

type Statis struct {
	Total int
	Success int 
	Fail int 
}

func (this *Statis) ToString() string {
	return fmt.Sprintf("total:%d, success:%d, fail:%d", 
		this.Total, this.Success, this.Fail)
}

var statis *Statis = &Statis{Total:0, Success:0, Fail:0}

func init() {
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/stop", stopHandler)
	http.HandleFunc("/statis", statisHandler)

	dispatcher = NewDispatcher(*maxworkers)
	dispatcher.Run()

	CreateTaskQueue(*maxTasks)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    //if r.Method != "POST" {
    //    w.WriteHeader(http.StatusMethodNotAllowed)
    //    return
    //}

    addr := strings.Split(r.RemoteAddr, ":") 
    ip := addr[0]

    statis.Total++

    task := Task{Data:UserData{Ip:ip}}
    select {
	case TaskQueue <- task:
	default:
		statis.Fail++
    	fmt.Println("TaskQueue is full !")
    	w.WriteHeader(http.StatusBadRequest)
    	return
	}
    
    w.WriteHeader(http.StatusOK)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    //if r.Method != "POST" {
    //    w.WriteHeader(http.StatusMethodNotAllowed)
    //    return
    //}

    dispatcher.Stop()

    w.WriteHeader(http.StatusOK)
}

func statisHandler(w http.ResponseWriter, r *http.Request) {
    //if r.Method != "POST" {
    //    w.WriteHeader(http.StatusMethodNotAllowed)
    //    return
    //}

	w.Write([]byte(statis.ToString()))

	statis.Total = 0
	statis.Success = 0
	statis.Fail = 0

    //dispatcher.Stop()

    w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}