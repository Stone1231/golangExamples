package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	// printGo()
	// sumGo()
	// sumNoGo()
	// exitGoroutine()
	// goSched()
	// multiGo()
	// chan1()
	// chan2()
	// showLenCap()
	// onlySendReceive()
	// selectChannel()
	// callNewConsumer()
	// semaphore()
	// quitByClosedChannel()
	// timeoutBySelect()
	// multiProcess()
	// waitGoFinish()
	// emptyStruct()
	// chanClose()
	// twoGosync()
	syncOnce()
}

func printGo() {
	go func() {
		for {
			fmt.Println("A")
		}
	}()

	for {
		fmt.Println("B")
	}
}

///WaitGroup
func sum(id int) {
	var x int64
	for i := 0; i < math.MaxUint32; i++ {
		x += int64(i)
	}
	println(id, x)
}

func sumGo() {
	wg := new(sync.WaitGroup)
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done() // = Add(-1)
			sum(id)
		}(i)
	}
	wg.Wait()
}

func sumNoGo() {
	for i := 0; i < 3; i++ {
		sum(i)
	}
}

///Goexit
func exitGoroutine() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer println("A.defer")
		func() {
			defer println("B.defer")
			runtime.Goexit() // 终止当前 goroutine
			println("B")     // 不会执行
		}()
		println("A") // 不会执行
	}()
	wg.Wait()
}

///runtime.Gosched 當前goroutine暫停,放回佇列等待下次被呼叫執行
func goSched() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	var A = func() {
		defer wg.Done()
		for i := 0; i < 6; i++ {
			println(i)
			time.Sleep(1 * time.Second)
			if i == 3 {
				runtime.Gosched()
			}
		}
	}
	var B = func() {
		defer wg.Done()
		println("Hello, World!")
	}
	//先A後B
	go B()
	go A()

	wg.Wait()
}

func multiGo() {
	wg := new(sync.WaitGroup)
	wg.Add(3)
	var A = func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			println("A", i)
			//time.Sleep(1 * time.Second)
		}
	}
	var B = func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			println("B", i)
			//time.Sleep(1 * time.Second)
		}
	}

	var C = func() {
		defer wg.Done()
		go B()
		for i := 0; i < 100; i++ {
			println("C", i)
			//time.Sleep(1 * time.Second)
		}
	}

	println("gogogo")
	//先A後B

	go C()
	go B()
	go A()

	wg.Wait()
}

///
func chan1() {
	data := make(chan int)  //數據交換
	exit := make(chan bool) //退出通知

	go func() {
		for d := range data { //佇列接收通知直到close
			time.Sleep(time.Second * 2)
			fmt.Println(d)

		}
		fmt.Println("recv over.")
		exit <- true //發出退出通知
	}()

	data <- 1 //發出數據
	data <- 2
	data <- 3

	close(data) //關閉佇列

	fmt.Println("send over.")

	<-exit //等待退出通知
	fmt.Println("exit")
}

func chan2() {
	data := make(chan int, 3) //緩衝區可儲存3個元素
	exit := make(chan bool)
	data <- 1 //緩衝區未滿前不會阻塞
	data <- 2
	data <- 3
	//data <- 4 // all goroutines are asleep - deadlock!//如果緩衝區已滿,阻塞
	go func() {
		for d := range data { // 緩衝區為空前不會阻塞
			fmt.Println(d)
		}
		exit <- true
	}()
	data <- 4
	data <- 5
	//close(data) //會error
	for index := 0; index < 3000; index++ {
		data <- index
	}
	fmt.Println("before close")
	close(data) //如果沒close會阻塞在上面的for …… range
	fmt.Println("after close")
	<-exit
	fmt.Println("exit")
}

//
func showLenCap() {
	//ok-idiom 模式判斷channel是否關閉
	// for {
	// 	if d, ok := <-data; ok {
	// 		fmt.Println(d)
	// 	} else {
	// 		break
	// 	}
	// }

	d1 := make(chan int)
	d2 := make(chan int, 3)
	// d1 <- 2 //無緩衝又沒接收,會error
	d2 <- 1
	fmt.Println(len(d1), cap(d1)) // 0  0
	fmt.Println(len(d2), cap(d2)) // 1  3
}

//
func onlySendReceive() {
	c := make(chan int, 3)
	var send chan<- int = c // send-only
	var recv <-chan int = c // receive-only

	send <- 1
	// <-send // Error: receive from send-only type chan<- int

	<-recv
	// recv <- 2 // Error: send to receive-only type <-chan int
}

//select channel
func selectChannel() {
	a, b := make(chan int, 3), make(chan int)
	go func() {
		v, ok, s := 0, false, ""
		for {
			select {
			case v, ok = <-a: //加上ok布林值表示是否成功接收
				s = "a"
			case v, ok = <-b:
				s = "b"
			}
			if ok {
				fmt.Println(s, v)
			} else {
				os.Exit(0)
			}
		}
	}()
	for i := 0; i < 5; i++ {
		select {
		case a <- i:
		case b <- i:
		}
	}
	close(a)
	select {} //沒有可用channel, 阻塞main goroutine
}

//簡單工廠模式
func NewConsumer() chan int {
	data := make(chan int, 3)
	go func() {
		for d := range data {
			fmt.Println(d)
		}
		os.Exit(0)
	}()
	return data
}
func callNewConsumer() {
	data := NewConsumer()
	data <- 1
	data <- 2
	close(data)
	select {}
}

///channel 實現號誌(Semaphore)
//Semaphore是一件可以容納N人的房間，如果人不滿就可以進去，
//如果人滿了，就要等待有人出來。
//對於N=1的情況，稱為binary semaphore。
//一般的用法是，用於限制對於某一資源的同時訪問。
func semaphore() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	sem := make(chan int, 1)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done()
			sem <- 1 //發送給sem, 阻塞或成功
			for x := 0; x < 3; x++ {
				fmt.Println(id, x)
			}
			<-sem //接收數據後,使其他阻塞可以發送數據
		}(i)
	}
	wg.Wait()
}

//使用closed channel發出退出通知
func quitByClosedChannel() {
	var wg sync.WaitGroup
	quit := make(chan bool)
	for i := 0; i < 2; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			task := func() {
				println(id, time.Now().Nanosecond())
				time.Sleep(time.Second)
			}
			for {
				select {
				case <-quit: //closed channel不會阻塞,可用作退出通知
					return
				default: //執行正常任務
					task()
				}
			}
		}(i)
	}

	time.Sleep(time.Second * 5) //讓goroutine執行一段時間
	close(quit)                 //發出退出通知
	wg.Wait()
	println("the end")
}

//select 實現超時
func timeoutBySelect() {
	w := make(chan bool)
	c := make(chan int, 2)
	go func() {
		select {
		case v := <-c:
			fmt.Println(v)
		case <-time.After(time.Second * 2):
			fmt.Println("timeout 2.")
		}
		w <- true
	}()
	//c <- 1 // 註解掉引發timeout
	<-w
}

///
type Request struct {
	data []int
	ret  chan int
}

func NewRequest(data ...int) *Request {
	return &Request{data, make(chan int, 1)}
}
func Process(req *Request) {
	x := 0
	for _, i := range req.data {
		x += i
	}
	time.Sleep(2 * time.Second)
	req.ret <- x
}
func multiProcess() {
	req := NewRequest(10, 20, 30)
	go Process(req)
	fmt.Println("go others...")
	fmt.Println(<-req.ret)
}

//
func waitGoFinish() {
	c := make(chan int) // Allocate a channel.
	// Start the sort in a goroutine; when it completes, signal on the channel.
	go func() {
		//list.Sort()
		time.Sleep(2 * time.Second)
		fmt.Println("B")
		c <- 1 // Send a signal; value does not matter.
	}()
	fmt.Println("A")
	//doSomethingForAWhile()
	<-c // Wait for sort to finish; discard sent value.
	fmt.Println("C")
}

func emptyStruct() {
	a := struct{}{}
	b := struct{}{}

	log.Println(a == b)            // true
	log.Printf("%p, %p\n", &a, &b) // 0x7bb7f8, 0x7bb7f8
}

func chanClose() {
	done := make(chan struct{})

	fmt.Printf("start")

	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()

	go func() {
		<-done
		fmt.Printf("done-1\n")
	}()

	go func() {
		<-done
		fmt.Printf("done-2\n")
	}()

loop:
	for {
		select {
		case <-done:
			time.Sleep(time.Second * 1)
			break loop
		}
	}

	fmt.Printf("end\n")
}

func twoGosync() {
	var x, y int
	go func() {
		x = 1                   //A1
		fmt.Print("y:", y, " ") //A2
	}()

	go func() {
		y = 1                   //B1
		fmt.Print("x:", x, " ") //B2
	}()

	//A1,B1,A2,B2 OR B1,A1,A2,B2..各種順序都有可能

	for {
	}
}

func syncOnce() {

	var do = func(o *sync.Once) {

		fmt.Println("Start do")

		o.Do(func() {

			fmt.Println("Doing something...") //只會執行一次

		})

		fmt.Println("Do end")
	}

	o := &sync.Once{}

	go do(o)

	go do(o)

	time.Sleep(time.Second * 2)
}