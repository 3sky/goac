package main

func main() {
	a := App{}
	
    a.Initialize("./SimpleDB.db")

	a.Run(":5000")
	defer a.DB.Close()
}