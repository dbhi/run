package lib

import "fmt"

func Exec(args []string) {

}

func exec(ts []Task) error {
	for _, t := range ts {
		fmt.Println(t)
	}
	return nil
}
