package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

func main() {
	person := &Person{
		Name:  "张三",
		Id:    20,
		Email: "222@qq.com",
		Phones: []*PhoneNumber{
			&PhoneNumber{
				Number: "111111",
				Type:   PhoneType_MOBILE,
			},
			&PhoneNumber{
				Number: "222222",
				Type:   PhoneType_HOME,
			},
		},
	}

	// 编码
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	// 解码
	msg := &Person{}
	err = proto.Unmarshal(data, msg)
	if err != nil {
		fmt.Println("unmarshal err:", err)
	}

	fmt.Println(person)
	fmt.Println(len(person.String()))
	fmt.Println(len(data))
	fmt.Println(msg)

}
