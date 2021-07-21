package api

import (
	"fmt"
	"testing"
)

func Test_ListServices(t *testing.T) {
	services, err := ListServices("127.0.0.1:2379", []string{"micro.financial", "micro.timeline"})
	fmt.Println(services, err)
}
