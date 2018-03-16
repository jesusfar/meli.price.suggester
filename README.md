#  Price suggester

Price suggester is a tool for suggest prices given a category ID of Mercado Libre Items.

## Installation

1.- Download

```
$ go get github.com/jesusfar/price.suggester

``` 
2.- Install dependencies

This project uses go dep, if you don't have it, you have to install go dep first.

```
$ cd $GOPATH/src/github.com/jesusfar/price.suggester
```

Install dependencies with:

```
$ dep ensure
```
3.- Running tests

Run unit test
```
$ go test -v ./...
```

Run benchmark tests
```
$ go test -v ./... -bench .
```
## How to use

Before to use the suggester, you need to fetch and train the data set of items.

### Fetch the data set

In order to fetch the items, we are using a Systematic Random Sampling method.

For example for each category we need to know the total amount of items and calc the size of the sampling data, 
then we get a random offset (K) to start the fetching items based on proportion (P).  

Fetching items for categories

```
$ go run main.go fetch

```
Fetching items for specific category.

```
$ go run main.go fetch MLA1743

```
### Train the data set

In order to suggest the prices, we need to train the data set of sampling data items.

```
$ go run main.go train

```
### Suggesting prices

Finally, we can suggest prices given a category ID. 

```
$ go run main.go suggest MLA1743

```
### Serve API

```
$ go run main.go serve

```
Test endpoint
```
$ curl -v http://localhost:8080/categories/MLA100028/prices
```
### Demo 
```
$ curl -v http://ec2-18-216-251-218.us-east-2.compute.amazonaws.com:8080/categories/MLA100028/prices
```
### License

MIT license
