package main

import "fmt"

func main() {
	scores := [8]float64{32,45,43,56,76,88,65,45}
    avg := 0.0

	for i:=0; i<8; i++ {
		avg+=scores[i]
	}

	avg = avg/8.0

    fmt.Printf("The average is %.2f\n",avg)
}
