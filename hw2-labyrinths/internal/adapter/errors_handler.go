package adapter

import (
	"fmt"
	"os"
)

func ShowError(err error) {
	fmt.Fprintf(os.Stderr, "Error:%v\n", err.Error())
}
