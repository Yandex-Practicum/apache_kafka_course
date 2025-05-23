package main

import (
	"fmt"
	"time"
)

func main() {
	sum := func(done <-chan struct{}, in <-chan int, n int) (<-chan int, <-chan struct{}) {
		out := make(chan int)
		heartbeat := make(chan struct{})

		go func() {
			defer close(out)
			defer close(heartbeat)
			pulse := time.Tick(time.Second * 1)
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case i := <-in:
					select {
					case <-done:
					case out <- i:
						time.Sleep(time.Second * 2)
					}
				}
			}
		}()

		return out, heartbeat
	}

	gen := func(done <-chan struct{}, ints ...int) <-chan int {
		out := make(chan int)

		go func() {
			defer close(out)
			for i := range ints {
				select {
				case <-done:
					return
				case out <- i:
				}
			}
		}()

		return out
	}

	done := make(chan struct{})
	defer close(done)

	in := gen(done, 2, 3, 4, 5, 6)
	ch, heartbeat := sum(done, in, 8)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-heartbeat:
				fmt.Println("heartbeat")
			}
		}
	}()

	for i := range ch {
		fmt.Println(i)
	}
}
