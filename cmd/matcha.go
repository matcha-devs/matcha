package main

type Matcha struct {
	db DatabaseInterface
}

func NewApp(db DatabaseInterface) *Matcha {
	return &Matcha{db}
}
