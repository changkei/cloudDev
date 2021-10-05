package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{}
	db.data = map[string]dollars{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	data     map[string]dollars
	dataLock sync.RWMutex
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	db.dataLock.RLock()
	fmt.Fprintf(w, "Available Items:\n")
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	db.dataLock.RUnlock()
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.dataLock.RLock()
	if price, ok := db.data[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	db.dataLock.RUnlock()
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	if priceDol, err := strconv.ParseFloat(price, 2); err == nil {
		db.dataLock.Lock()
		db.data[item] = dollars(priceDol)
		db.dataLock.Unlock()
		fmt.Fprintf(w, "%s: %s\n", item, price)
	} else {
		fmt.Fprintf(w, "Please enter a valed number: %q\n", price)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	iprice := req.URL.Query().Get("price")
	if price, ok := db.data[item]; ok {
		dol, _ := strconv.ParseFloat(iprice, 32)
		dollar := dollars(dol)
		db.dataLock.Lock()
		db.data[item] = dollar
		db.dataLock.Unlock()
		fmt.Fprintf(w, "Old price %s\n", price)
		fmt.Fprintf(w, "New price %s\n", db.data[item])
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db.data[item]; ok {
		db.dataLock.Lock()
		delete(db.data, item)
		db.dataLock.Unlock()
		fmt.Fprintf(w, "Deleted %s: %s\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
