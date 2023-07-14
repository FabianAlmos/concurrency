package worker

import (
	"concuLec/db"
	model "concuLec/models"
	response "concuLec/responses"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	id       int
	supplier response.SupplierResponse
}

type Result struct {
	id int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)
var suppliers = &response.SupplierCollectionResponse{}

func Start() {
	for {
		jobs = make(chan Job, 10)
		<-time.After(5 * time.Second)
		getSuppliers()
		allocate()
		createWorkerPool()
	}
}

func getSuppliers() {
	resp, err := http.Get("https://foodapi.golang.nixdev.co/suppliers")
	if err != nil {
		fmt.Println(err)
	}

	err = json.NewDecoder(resp.Body).Decode(suppliers)
	if err != nil {
		fmt.Println(err)
	}
}

func allocate() {
	for i := 0; i < len(suppliers.Suppliers); i++ {
		job := Job{i, suppliers.Suppliers[i]}
		jobs <- job
	}
	close(jobs)
}

func updateDB(menu_id int, menuItem model.Menu) int {
	db := db.GetDB()
	stmt, err := db.Prepare("UPDATE menu_items SET image = $1, ingredients = $2, name = $3, price = $4, type = $5 WHERE menu_id = $6")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(menuItem.Image, menuItem.Ingredients, menuItem.Name, menuItem.Price, menuItem.Type, menu_id)
	if err != nil {
		panic(err)
	}
	return menu_id
}

func processMenu(id int) {
	resp, err := http.Get(fmt.Sprintf("https://foodapi.golang.nixdev.co/suppliers/%d/menu", id))
	if err != nil {
		fmt.Println(err)
	}

	menu := &response.MenuResponse{}
	err = json.NewDecoder(resp.Body).Decode(menu)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(menu.Menu); i++ {
		result := Result{id: updateDB(id, menu.Menu[i])}
		results <- result
	}
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		processMenu(job.id)
	}
	wg.Done()
}

func createWorkerPool() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(wg)
	}
	wg.Wait()
}
