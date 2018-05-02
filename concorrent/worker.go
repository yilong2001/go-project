package main


import(
//"strings"
//"fmt"
"log"
)

type Task struct {
	Data UserData
}

var TaskQueue chan Task;

func CreateTaskQueue(maxtasks int) {
    TaskQueue = make(chan Task, maxtasks)
}


type Worker struct {
    WorkerPool chan chan Task
    TaskPool chan Task
    quit chan bool
}

func NewWorker(pool chan chan Task) *Worker {
    return &Worker{
        WorkerPool: pool,
        TaskPool: make(chan Task),
        quit:     make(chan bool),
    }
}


func (this *Worker) Run() {
    go func() {
        for {
            this.WorkerPool <- this.TaskPool
            select {
                //等待任务,woker首先出队,然后才有任务
                case task:= <- this.TaskPool:
                    if err := task.Data.Send(); err != nil {
                        log.Println("send wrong : ", err)
                    }
                case <- this.quit:
                    log.Println(" worker quit. ")
                    return
            }
        }
    }()
}

func (this *Worker) Stop() {
    go func() {
        this.quit <- true
    }()
}
