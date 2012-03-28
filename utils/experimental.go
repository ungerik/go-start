package utils

func IntRange(start int, end int) <-chan int {
	return IntRangeStep(start, end, 1)
}

func IntRangeStep(start int, end int, step int) <-chan int {
	result := make(chan int)
	go func() {
		for i := start; i < end; i += step {
			result <- i
		}
		close(result)
	}()
	return result
}
