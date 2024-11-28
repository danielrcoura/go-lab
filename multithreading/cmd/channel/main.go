package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// será que ao encerrar um canal, todas as go routines recebem o sinal de encerramento?

func main() {
	c := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// go workerRange(1, c, wg)
	// go workerRange(2, c, wg)
	go workerSelect(1, c, wg)
	go workerSelect(2, c, wg)

	c <- "message1"
	time.Sleep(1 * time.Second)
	c <- "message2"
	close(c)

	wg.Wait()
	fmt.Println("goroutines:", runtime.NumGoroutine())
}

// o case será executado se o canal estiver pronto para leitura, caso contrário o default é executado
func workerSelect(id int, c chan string, wg *sync.WaitGroup) {
	// o select não fica ouvindo os canais, ele faz uma única verificação e executa o case que estiver pronto
	// o for garante que a verificação aconteça continuamente
	for {
		select {
		case x, ok := <-c: // executado quando o canal está pronto para leitura
			if !ok {
				fmt.Printf("worker %d: channel closed\n", id)
				wg.Done()
				return
			}
			fmt.Printf("worker %d: received %s\n", id, x)
		// sem esse default o select aguarda até que um case esteja pronto, bloqueando a execução da thread
		default: // executado quando o canal não está pronto para leitura
			fmt.Printf("worker %d: channel is not ready for receiving\n", id)
			// sem esse sleep o loop fica muito rápido, executando o print acima várias vezes
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func workerRange(id int, c chan string, wg *sync.WaitGroup) {
	// recebe valores de um canal até que ele seja fechado
	for x := range c {
		fmt.Printf("worker %d: received %s\n", id, x)
	}
	fmt.Printf("worker %d: channel closed\n", id)
	wg.Done()
}
