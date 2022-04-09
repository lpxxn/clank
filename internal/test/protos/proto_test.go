package protos_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lpxxn/clank/internal/test/protos/model"
)

func TestStudent1(t *testing.T) {
	s1 := &model.Student{
		Id:   1234567890,
		Name: "五六七",
		Age:  18,
	}
	sBody, err := proto.Marshal(s1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(sBody))
	revS1 := &model.Student{}
	if err := proto.Unmarshal(sBody, revS1); err != nil {
		t.Fatal(err)
	}
	t.Log(revS1)
}

func TestStudent2(t *testing.T) {
	s1 := &model.Student{
		Id:   1,
		Name: "孙悟空",
		Age:  300,
	}
	sBody, err := proto.Marshal(s1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%b", sBody)
	t.Logf("%d", sBody)
	t.Log(string(sBody))
	jBody, _ := json.Marshal(s1)
	t.Log(string(jBody))
	revS1 := &model.Student{}
	if err := proto.Unmarshal(sBody, revS1); err != nil {
		t.Fatal(err)
	}
	t.Log(revS1)
}

func TestStudentList1(t *testing.T) {
	s1 := &model.StudentList{
		Class: "三年级二班",
		Students: []*model.Student{
			&model.Student{Id: 123465, Name: "路飞", Age: 19},
			&model.Student{Id: 321, Name: "索龙", Age: 20},
			&model.Student{Id: 789, Name: "乔巴", Age: 6},
		},
		Teacher: "雷利",
		Score:   []int64{1, 2, 3, 4, 5, 6},
	}
	sBody, err := proto.Marshal(s1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(sBody))
	jBody, _ := json.Marshal(s1)
	t.Log(string(jBody))
	revS1 := &model.StudentList{}
	if err := proto.Unmarshal(sBody, revS1); err != nil {
		t.Fatal(err)
	}
	t.Log(revS1)
}

func TestSync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	go func() {
		// this go routine wait for changes to the sharedRsc
		fmt.Println("1 ini")
		c.L.Lock()
		fmt.Println("1 lock")

		for len(sharedRsc) == 0 {
			fmt.Println("1 wait")
			c.Wait()
		}
		fmt.Println(sharedRsc["rsc1"])
		c.L.Unlock()
		wg.Done()
	}()

	go func() {
		// this go routine wait for changes to the sharedRsc
		fmt.Println("2 ini")
		c.L.Lock()
		fmt.Println("2 lock")
		for len(sharedRsc) == 0 {
			fmt.Println("2 wait")
			c.Wait()
		}
		fmt.Println(sharedRsc["rsc2"])
		c.L.Unlock()
		wg.Done()
	}()

	// this one writes changes to sharedRsc
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("writes")
		c.L.Lock()
		sharedRsc["rsc1"] = "foo"
		sharedRsc["rsc2"] = "bar"
		c.Broadcast()
		c.L.Unlock()
	}()

	wg.Wait()
}

var sharedRsc = make(map[string]interface{})
