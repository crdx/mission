package task

import "fmt"

type Tasks []Task

func (self Tasks) Print(long bool) {
	for _, task := range self {
		if long {
			fmt.Println(task.GetLongString())
		} else {
			fmt.Println(task.GetShortString())
		}
	}
}
