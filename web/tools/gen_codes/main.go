package main

import (
	"fmt"
	"test.com/project-common/encrypts"
)

func main() {
	// member ids, org ids, project ids, stage ids, task ids
	ids := []int64{1013, 1014, 29, 30, 13046, 13047, 13048, 88, 89, 90, 91, 98, 99, 100, 101, 102, 120}
	for _, id := range ids {
		code := encrypts.EncryptNoErr(id)
		fmt.Printf("id=%d => code=%s\n", id, code)
	}
}
