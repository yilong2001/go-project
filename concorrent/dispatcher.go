package main


import(
"fmt"
//"strings"
)

type Dispatcher struct {
    WorkerPool chan chan Task
    maxWorkers int
    wokers []*Worker
    quit chan bool
}

func NewDispatcher(maxWorkers int) *Dispatcher {
    pool := make(chan chan Task, maxWorkers)
    return &Dispatcher{WorkerPool: pool,
        maxWorkers:maxWorkers, 
        wokers:[]*Worker{},
        quit:make(chan bool)}
}

func (this *Dispatcher) Run() {
    for i:=0; i<this.maxWorkers; i++ {
        worker := NewWorker(this.WorkerPool)
        worker.Run()
        this.wokers = append(this.wokers, worker)
    }

    go this.dispatch()
}

func (this *Dispatcher) Stop() {
    for i:=0; i<this.maxWorkers; i++ {
        this.wokers[i].Stop()
    }

    go func() {
        this.quit <- true
    }()
}

func (this *Dispatcher) dispatch() {
    for {
        select {
        case task := <- TaskQueue:
            go func(task Task) {
                //等待空闲woker(即任务池)
                taskPool := <- this.WorkerPool
                taskPool <- task
            }(task)
        case <- this.quit:
            fmt.Println(" dispatch quit. ")
            return
        }
    }
}

