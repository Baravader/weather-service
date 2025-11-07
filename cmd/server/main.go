package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const listeningPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			log.Println(err)
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(s)
	fmt.Println(jobs[0].ID())
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("start server on port", listeningPort)
		err := http.ListenAndServe(listeningPort, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("starts schedule,job:", jobs[0].ID())
		s.Start()
	}()
	wg.Wait()
}

func initJobs(scheduler gocron.Scheduler) ([]gocron.Job, error) {

	j, err := scheduler.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("works!")
			},
		),
	)

	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}

//func runCron() {
//	if err != nil {
//		// handle error
//	}
//	// each job has a unique id
//	fmt.Println(j.ID())
//
//	// start the scheduler
//	s.Start()
//}
