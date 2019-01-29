package main

// import "os"

func main() {
	a := App{}
	
    a.Initialize("./SimpleDB.db")

    a.Run(":5000")
}