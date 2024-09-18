package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)




func main(){

	var pro chan Pac
	pro = make(chan Pac)


	go Server(pro)

	go Client(pro)


	for {

	}
}

func Client(pro chan Pac){
	var syn = rand.IntN(100)
	var ack Pac = Pac{-2,-2,""}



	for ack == (Pac{-2,-2,""}){
		fmt.Println("starts")
		pro <- Pac{-1, syn, "ABC"}
		fmt.Println("b")
		time.Sleep(1 * time.Second)
		go ClientCheck(&ack,syn, pro)
		time.Sleep(5 * time.Millisecond)
		if ack.ack != syn+1{
			ack = Pac{-2,-2,""}
		}
	}
	fmt.Println("Made it out")
	fmt.Printf("Ack: ack %d, sec %d, data %s \n", ack.ack, ack.sec, ack.data)

	pro <- Pac{ack.sec+1, ack.ack, "ABC"}

	for {

	}
}

func ClientCheck(ackPointer *Pac, syn int, pro chan Pac){
	temp := <- pro //Kill package: Pac{-1,-1,""}

	if temp.ack != -1 && temp.sec != -1 || syn != temp.sec{ //Kill package actiaved
		fmt.Println("Might run")
		*ackPointer = temp
	}

	
}

func Server(pro chan Pac){

	temp := <- pro
	fmt.Println("server got message")
	var syn = rand.IntN(100)+100



	pro <- Pac{temp.sec+1, syn, temp.data}

	temp2 := <- pro

	if temp2.ack == syn+1 && temp2.sec == temp.ack+1{

	}


	for {

	}
}

type Pac struct {
	ack int
	sec int
	data string
}