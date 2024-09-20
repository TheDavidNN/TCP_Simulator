package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {

	var pro chan Pac
	pro = make(chan Pac)

	go Server(pro)

	go Client(pro)

	go Resender(pro)

	for {

	}
}

func Resender(pro chan Pac) {
	for {
		var tmp = <-pro
		if rand.IntN(100) >= 50 {
			pro <- tmp
		}
	}
}

func Client(pro chan Pac) {
	var syn = rand.IntN(100)
	var ack Pac = Pac{-2, -2, ""}

	var tmp0 Pac

	/*
		fmt.Println("Pac send")
			pro <- Pac{-1, syn, "ABC"}
			time.Sleep(500 * time.Millisecond)
			go ClientCheck(&ack, syn, pro)
			fmt.Println("sub routine activated")
			time.Sleep(500 * time.Millisecond)
			if ack.ack != syn+1 {
				fmt.Println("message not approved")
				ack = Pac{-2, -2, ""}
	*/

	var sendFirstMsg bool = false

	for ack == (Pac{-2, -2, ""}) {

		if !sendFirstMsg {
			fmt.Println("client, send message")
			ack = SendMessageRecieveAnswer(pro, Pac{-1, syn, "ABC"}, syn+1, 800, &sendFirstMsg)
			fmt.Println("client, recieve message")
		}

		pro <- Pac{ack.sec + 1, ack.ack, "ABC"}
		time.Sleep(1 * time.Second)
		tmp0 = <-pro
		if tmp0.ack == syn+1 {
			ack = Pac{-2, -2, ""}
		}
	}

	time.Sleep(1 * time.Second)

	fmt.Println("client thinks connection established")

}

func CheckIfMessage(pro chan Pac, check *bool) {
	tmp := <-pro

	pro <- tmp

	*check = true
}

func SendMessage(pro chan Pac, packag Pac) {
	pro <- packag
}

func CheckMessage(pro chan Pac, packag Pac, exspectedAck int) bool {
	if packag.ack == exspectedAck {
		return true
	}
	return false
}

func SendMessageRecieveAnswer(pro chan Pac, packag Pac, exspectedAck int, tim int, noLongerDoThis *bool) Pac {

	for {
		SendMessage(pro, packag)
		fmt.Printf("Sending message packag: sec: %d\n", packag.sec)

		time.Sleep(time.Duration(tim) * time.Millisecond)
		var boo bool = false
		go CheckIfMessage(pro, &boo)
		time.Sleep(100 * time.Millisecond)
		var tmp Pac = Pac{-1, -1, ""}
		if boo {
			tmp = <-pro
			fmt.Printf("recieve message %d\n", tmp.ack)
		}
		if CheckMessage(pro, tmp, exspectedAck) {
			*noLongerDoThis = true
			return tmp
		}
	}
	return Pac{-1, -1, ""}

}

func Server(pro chan Pac) {

	temp := <-pro
	var syn = rand.IntN(100) + 100

	connectionEstablished := false

	var t bool = false

	for connectionEstablished == false {
		fmt.Println("Server, sends pac")
		var temp2 = SendMessageRecieveAnswer(pro, Pac{temp.sec + 1, syn, temp.data}, syn+1, 400, &t)
		fmt.Println("Server, recieve pac")
		if temp2.ack == syn+1 && temp2.sec == temp.sec+1 {
			fmt.Println("Server, sub activates")
			connectionEstablished = true
		}
	}

	fmt.Println("Connection fully established")
}

type Pac struct {
	ack  int
	sec  int
	data string
}
