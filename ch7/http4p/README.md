# TEST

List all items (Read):  
http://localhost:8000/list  
```sh
curl http://localhost:8000/list
```

Get the price of a single item (Read):  
http://localhost:8000/price?item=shoes  
```sh
curl http://localhost:8000/price?item=socks
curl http://localhost:8000/price?item=nonexistent  # Returns 404
```

Create a new item (Create):  
http://localhost:8000/create?item=hat&price=12.50  
```sh
curl "http://localhost:8000/create?item=hat&price=12.50"
# Creating 'hat' again returns 409 Conflict
```

Update item price (Update):  
http://localhost:8000/update?item=shoes&price=55.00  
```sh
curl "http://localhost:8000/update?item=shoes&price=55.00"
# Updating nonexistent item returns 404
curl "http://localhost:8000/update?item=shoes&price=invalid"  # Returns 400 Bad Request
```

Delete an item (Delete):  
http://localhost:8000/delete?item=socks  
```sh
curl http://localhost:8000/delete?item=socks
# Deleting nonexistent item returns 404
```
