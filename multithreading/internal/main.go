package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	threadMain()
}

// ------------------------------------------------------
// # 1. THREAD MAIN
// ------------------------------------------------------
//
// A go routine está sendo criada através da thread main, mas nenhuma operação da go routine
// é realiza. O programa encerra imediatamente porque a thread main finalizou sua execução.

func threadMain() {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// ------------------------------------------------------
// # 2. SLEEP
// ------------------------------------------------------
//
// Um time.sleep mantém a thread main ativa por um curto período de tempo,
// interrompendo a execução da go routine na metade.

func sleep() {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	time.Sleep(500 * time.Millisecond)
}

// ------------------------------------------------------
// # 3. WAIT GROUP
// ------------------------------------------------------
//
// Agora o WaitGroup está segurando a thread main, aguardando o sinal que será enviado
// pela go routine. Dessa forma, é garantido que o programa só irá encerrar após a
// execução completa da go routine.

func waitgroup() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
		wg.Done()
	}()
	wg.Wait()
}

// ------------------------------------------------------
// # 4. FOREVER
// ------------------------------------------------------
//
// O comportamento do wait group anterior pode ser reproduzido utilizando um channel.
// A thread main está de pé porque está aguardando que algo seja produzido no
// channel "forever" para consumir em seguida.

func forever() {
	forever := make(chan string)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
		forever <- "finalizado"
	}()
	fmt.Println(<-forever)
}

// ------------------------------------------------------
// # X. CONTEXT
// ------------------------------------------------------
//
// Mesmo comportamento utilizando o context. A thread main ficará de pé enquanto
// o context não for cancelado.

func contextExample() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
		cancel()
	}()
	<-ctx.Done()
}

// ------------------------------------------------------
// # 5. CHANNELS
// ------------------------------------------------------
//
// A go routine envia dados para o canal enquanto que a thread main as consome.
// O close(ch) evita que ocorra um deadlock, pois a thread main tentaria consumir
// o channel eternamente.

func pubsub() {
	ch := make(chan int)
	go func(ch chan<- int) {
		for i := 0; i < 10; i++ {
			ch <- int(i)
		}
		close(ch)
	}(ch)
	for i := range ch {
		fmt.Println(i)
	}
}

// ------------------------------------------------------
// # 6. CHANNELS + WAIT GROUPS
// ------------------------------------------------------
//
// Implementação com uma go routine produtora, uma consumidora e um wait group
// segurando o encerramento da thread main.

func pubsubwg() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ch chan<- int, wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}(ch, &wg)
	go func(ch <-chan int) {
		for i := range ch {
			fmt.Println(i)
		}
		wg.Done()
	}(ch)
	wg.Wait()
}

// ------------------------------------------------------
// # 7. LOAD BALANCER
// ------------------------------------------------------
//
// Implementação de 100.000 threads para atender 1.000.000 de requests.
// O wait group assegura que nenhum request será perdido, pois impede que a thread
// main encerre antes que os workers terminem de executar.

func loadbalancer() {
	REQUESTS := 1000001
	WORKERS := 100000
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(REQUESTS)
	for i := 0; i < WORKERS; i++ {
		go worker(i, ch, &wg)
	}
	for i := 0; i < REQUESTS; i++ {
		ch <- i
	}
	wg.Wait()
}

func worker(id int, ch chan int, wg *sync.WaitGroup) {
	for x := range ch {
		fmt.Printf("worker #%d received %d\n", id, x)
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}
}

// ------------------------------------------------------
// # 7. LOAD BALANCER + ATOMIC
// ------------------------------------------------------
//
// A implementação realiza uma contagem de requisições. A concorrência entre as threads
// pode causar um erro na contagem. O mutex faz um lock na variável de soma, impedindo que
// outras thread alterem o valor antes que a thread atual faça um unlock. O pacote "atomic"
// simplifica o uso de Mutex.

func atomicCount() {
	count := uint64(0)
	REQUESTS := 1000001
	WORKERS := 100000
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(REQUESTS)
	for i := 0; i < WORKERS; i++ {
		go workerAtomic(i, &count, ch, &wg)
	}
	for i := 0; i < REQUESTS; i++ {
		ch <- i
	}
	wg.Wait()
	fmt.Printf("Total requests: %d\n", count)
}

func workerAtomic(id int, count *uint64, ch chan int, wg *sync.WaitGroup) {
	for x := range ch {
		atomic.AddUint64(count, 1)
		// *count += uint64(1)
		fmt.Printf("worker #%d received %d\n", id, x)
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}
}

// ------------------------------------------------------
// # 8. PROMISE RACE
// ------------------------------------------------------
//
// Comportamento do Promise.race em go utilizando select. O select irá executar
// quando o primeiro canal enviar e o resultado irá para result.

func promiseRace() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func(ch chan string) {
		time.Sleep(100 * time.Millisecond)
		ch <- "channel 1"
	}(ch1)

	go func(ch chan string) {
		time.Sleep(200 * time.Millisecond)
		ch <- "channel 2"
	}(ch2)

	var result string
	select {
	case m1 := <-ch1:
		result = m1
	case m2 := <-ch2:
		result = m2
	case <-time.After(300 * time.Millisecond):
		panic("timeout")
	}

	fmt.Println(result)
}

// ------------------------------------------------------
// # 9. PROMISE ALL - FOR + SELECT
// ------------------------------------------------------
//
// Comportamento do Promise.all em go utilizando select. O select é executado duas
// vezes, uma para cada go routine.

func promiseAllForSelect() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func(ch chan string) {
		time.Sleep(200 * time.Millisecond)
		ch <- "channel 1"
	}(ch1)

	go func(ch chan string) {
		time.Sleep(100 * time.Millisecond)
		ch <- "channel 2"
	}(ch2)

	result := [2]string{}
	for i := 0; i < 2; i++ {
		select {
		case m1 := <-ch1:
			result[0] = m1
		case m2 := <-ch2:
			result[1] = m2
		case <-time.After(1000 * time.Millisecond):
			panic("timeout")
		}
	}

	fmt.Println(runtime.NumGoroutine())
	fmt.Println(result)
}
