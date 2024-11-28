package main

import (
	"fmt"
	"time"
)

// O select aguarda até que um canal esteja pronto para leitura, bloquenado a execução da thread.
// Se o default for implementado, ele será executado quando nenhum canal estiver pronto para a leitura, evitando que a execução da thread seja bloqueada.
func main() {
	c := make(chan string)

	fmt.Println("running blocking select...")
	go sayHello(c)
	select {
	case msg := <-c:
		fmt.Println(msg)
	}

	fmt.Println("running non-blocking select...")
	go sayHello(c)
	select {
	case msg := <-c:
		fmt.Println(msg)
	default:
		fmt.Println("no message received")
	}
}

func sayHello(c chan<- string) {
	time.Sleep(2 * time.Second)
	c <- "hello! (after 2 seconds)"
}
