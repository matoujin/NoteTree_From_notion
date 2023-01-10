package main

import (
	"fmt"
	"sync"
	"time"
)

var JobChannel = make(chan Task)

type Task struct {
	arg string
}

type Pool struct {
	//entrypoint
	taskchannel chan Task
	//Max worker
	workerNum int
}

func NewTask(argu string) {
	fmt.Println("Find newTask, please waiting... ")
	t := Task{
		arg: argu,
	}
	JobChannel <- t
}

func NewPool(workNum int) *Pool {

	p := Pool{
		taskchannel: make(chan Task),
		workerNum:   workNum,
	}
	return &p
}

func (p *Pool) worker(wg *sync.WaitGroup) {

	defer wg.Done()
	for task := range p.taskchannel {
		SBlockChildren(task.arg)
	}
	//wg.Done()
}

//协程池工作
func (p *Pool) poolStartWork(waittime int) {

	var wg sync.WaitGroup

	for i := 0; i < p.workerNum; i++ {
		wg.Add(1)
		go p.worker(&wg)
	}

	go func() {
		for task := range JobChannel {
			p.taskchannel <- task
		}
	}()

	//超时机制
	var done = make(chan int)
	go func() {
		wg.Wait()
		done <- 1

	}()

loop:
	select {
	case <-done:
		goto loop
	case <-time.After(time.Duration(waittime) * time.Second):
		fmt.Printf("*Timed out waiting for wait group.*\nThe work may have done.\n\n")
		close(p.taskchannel)
		fmt.Println("Channel 1 has closed\n")
		close(JobChannel)
		fmt.Println("Channel 2 has closed\n")

	}

}

func InitPool(waittime int) {

	p := NewPool(8)
	p.poolStartWork(waittime)

}
