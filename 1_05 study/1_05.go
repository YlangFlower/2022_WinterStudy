package main

import (
	"bufio"
	"fmt"
	"os"
)

import (
	"math"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscanln(reader, &n, &m)

	var board = make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscanf(reader, "%v\n", &board[i])
	}
	var num = n * m

	for i := 0; i < n-7; i++ {
		for j := 0; j < m-7; j++ {
			var cnt1 = float64(0) // B -> W
			var cnt2 = float64(0) // W -> B
			for k := i; k < i+8; k++ {
				for l := j; l < j+8; l++ {

					if (k+l)%2 == 0 { // 좌표의 합이 짝수인 곳
						if string(board[k][l]) == "B" {
							cnt1++
						} else {
							cnt2++
						}
					} else { // 좌표의 합이 홀수인 곳
						if string(board[k][l]) == "B" {
							cnt2++
						} else {
							cnt1++
						}
					}

				}
			}
			if num > int(math.Min(cnt1, cnt2)) {
				num = int(math.Min(cnt1, cnt2))
			}
		}
	}
	fmt.Println(num)
}
