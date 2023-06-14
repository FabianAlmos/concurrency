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

func Timer() {
	for {
		<-time.After(10 * time.Minute)
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
	for i := 1; i <= len(suppliers.Suppliers); i++ {
		job := Job{i, suppliers.Suppliers[i]}
		jobs <- job
	}
	close(jobs)
}

func updateDB(supplier_id int, menuItem model.Menu) {
	db := db.GetDB()
	stmt, err := db.Prepare("INSERT INTO menu_items (image, ingredients, name, price, type) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(menuItem.Image, menuItem.Ingredients, menuItem.Name, menuItem.Price, menuItem.Type)
	if err != nil {
		panic(err)
	}
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
		updateDB(id, menu.Menu[i])
	}
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		processMenu(job.id)
	}
	wg.Done()
}

func createWorkerPool() {
	var wg *sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(wg)
	}
	wg.Wait()
}
