// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 195.

// Http4 is an e-commerce server that registers the /list and /price
// endpoint by calling http.HandleFunc.
// Practice 7.11: Implement create, update, and delete operations
// Practice 7.12: Render the item list as an HTML table
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// ----------------------------------------------------
// HTML Template Definition
// ----------------------------------------------------

// Define the HTML template for the item list.
// We parse it once globally for efficiency.
var itemListTemplate = template.Must(template.New("itemlist").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Store Inventory</title>
    <style>
        body { font-family: sans-serif; margin: 20px; }
        table {
            width: 50%; /* Adjust as needed */
            border-collapse: collapse;
            margin-top: 20px;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <h1>Store Inventory</h1>
    <table>
        <thead>
            <tr>
                <th>Item</th>
                <th>Price</th>
            </tr>
        </thead>
        <tbody>
            {{range $item, $price := .}}
            <tr>
                <td>{{$item}}</td>
                <td>{{$price}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>
`))

//!+main

func main() {
	db := database{items: map[string]dollars{"shoes": 50, "socks": 5}, mu: &sync.Mutex{}}
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//!-main

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	mu    *sync.Mutex
	items map[string]dollars
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	db.mu.Lock()
	itemsToRender := make(map[string]dollars)
	for k, v := range db.items {
		itemsToRender[k] = v
	}
	db.mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := itemListTemplate.Execute(w, itemsToRender); err != nil {
		log.Printf("Template execution error for /list: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.mu.Lock()
	price, ok := db.items[item]
	db.mu.Unlock()

	if ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

// create handles creating a new item (Create operation)
func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")

	if item == "" {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "item name is required\n")
		return
	}
	if priceStr == "" {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "price is required to create a new item\n")
		return
	}

	newPrice, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "invalid price format: %q\n", priceStr)
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.items[item]; ok {
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		fmt.Fprintf(w, "item %q already exists; use /update to modify\n", item)
		return
	}

	db.items[item] = dollars(newPrice)
	fmt.Fprintf(w, "Item %q created with price %s\n", item, dollars(newPrice))
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")

	if item == "" || priceStr == "" {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "item and price are required\n")
		return
	}

	newPrice, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "invalid price format: %q\n", priceStr)
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.items[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		fmt.Fprintf(w, "no such item: %q; use /create to add new items\n", item)
		return
	}

	db.items[item] = dollars(newPrice)
	fmt.Fprintf(w, "Item %q updated to price %s\n", item, dollars(newPrice))

}

// delete handles deleting an item (Delete operation)
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if item == "" {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "item name is required to delete\n")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.items[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	delete(db.items, item) // Go's built-in delete function for maps
	fmt.Fprintf(w, "Item %q deleted\n", item)
}
