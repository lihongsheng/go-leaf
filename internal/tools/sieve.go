package tools

func sieveGenerate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func sieveFilter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- i
		}
	}
}

func Sieve(max int64) {
	ch := make(chan int)
	go sieveGenerate(ch)
	for i := int64(0); i < max; i++ {
		prime := <-ch
		print(prime, "\n")
		ch1 := make(chan int)
		go sieveFilter(ch, ch1, prime)
		ch = ch1
	}
}
