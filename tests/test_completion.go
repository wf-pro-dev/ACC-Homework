package main

import (
	"fmt"

	"github.com/williamfotso/acc/cmd/completion"
)

func main() {

	deadline_completion, _ := completion.DeadlineCompletion("1")
	fmt.Println(deadline_completion)
}
