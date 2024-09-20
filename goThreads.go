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
		if rand.IntN(100) >= 95 {
			pro <- tmp
		}
	}
}

func Client(pro chan Pac) {
	var seq = rand.IntN(100)

	var synAck Pac = Pac{-1, -1}

	for synAck.ack == -1 { //Send syn and get synack
		sendSyn(pro, seq)

		go getSynAck(pro, &synAck, seq)

		time.Sleep(300*time.Millisecond)
	} //Now we have recieved a acknoladgement from the server.

	Ack := Pac{synAck.seq+1,synAck.ack}

	for {
		sendAck(pro,Ack)
		fmt.Println("Client believes a connection has been established")
		getSynAck(pro,&synAck,seq)
		time.Sleep(300*time.Millisecond)
	}

}

func sendAck(pro chan Pac, Ack Pac){
	fmt.Println("Client, sendAck")
	pro <- Ack //Sending ack

	time.Sleep(500 * time.Millisecond)
}

func getSynAck(pro chan Pac, synAck *Pac, seq int) {
	fmt.Println("Client, getSynAck")
	var recieve Pac = Pac{-1, -1}

	recieve = <-pro //Getting syn-ack

	if recieve.ack == seq+1 {
		*synAck = recieve
	}
}

func sendSyn(pro chan Pac, seq int) {
	fmt.Println("Client, sendSyn")
	pro <- Pac{-1, seq} //Sending syn

	time.Sleep(500 * time.Millisecond)
}

func Server(pro chan Pac) {
	var seq = rand.IntN(100) + 100

	firstMsg := <-pro

	Ack := Pac{-1,-1}

	for Ack.ack == -1{
		sendSynAck(pro, seq,firstMsg.seq+1)

		go recieveAck(pro,seq,&Ack)
		time.Sleep(300*time.Millisecond)
	}

	fmt.Println("Server believes connection has been established")
}

func recieveAck(pro chan Pac, seq int, Ack *Pac){
	fmt.Println("Server, recieveAck")
	recieve := <- pro

	if  recieve.ack == seq+1{
		*Ack = recieve
	}
}

func sendSynAck(pro chan Pac, seq int, ack int){
	fmt.Println("Server, sendSynAck")
	pro <- Pac{ack,seq}
	time.Sleep(500 * time.Millisecond)
}

type Pac struct {
	ack int
	seq int
}
