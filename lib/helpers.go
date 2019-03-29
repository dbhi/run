package lib

import (
	"fmt"
	"os"

	au "github.com/logrusorgru/aurora"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(au.Red(err))
		os.Exit(1)
	}
}
