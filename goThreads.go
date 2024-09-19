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

func Resender(pro chan Pac){
	for{
		var tmp = <-pro

		pro <- tmp
	}
}

func Client(pro chan Pac) {
	var syn = rand.IntN(100)
	var ack Pac = Pac{-2, -2, ""}

	for ack == (Pac{-2, -2, ""}) {
		fmt.Println("Pac send")
		pro <- Pac{-1, syn, "ABC"}
		time.Sleep(1 * time.Second)
		go ClientCheck(&ack, syn, pro)
		fmt.Println("sub routine activated")
		time.Sleep(5 * time.Millisecond)
		if ack.ack != syn+1 {
			fmt.Println("message not approved")
			ack = Pac{-2, -2, ""}
		}
	}

	time.Sleep(500*time.Millisecond)

	fmt.Println("client thinks connection established")

	pro <- Pac{ack.sec + 1, ack.ack, "ABC"}

	var serverDown bool = false

	var messageRecieved bool = false

	for {
		if messageRecieved == false {

			var tmp0 = <-pro

			var tmp1 = <-pro

			var tmp2 = <-pro

			if tmp0.sec > tmp1.sec {
				var tmp = tmp0
				tmp0 = tmp1
				tmp1 = tmp
			}
			if tmp1.sec > tmp2.sec {
				var tmp = tmp1
				tmp1 = tmp2
				tmp2 = tmp
			}

			pro <- Pac{tmp0.sec + 1, -1, ""}

			pro <- Pac{tmp1.sec + 1, -1, ""}

			pro <- Pac{tmp2.sec + 1, -1, ""}

			fmt.Printf("The message was: %s \n", tmp0.data+tmp1.data+tmp2.data)
			messageRecieved = true
		}

		if serverDown == false {
			var end0 = <-pro
			pro <- Pac{end0.sec + 1, -1, ""}
			serverDown = true
		}

		time.Sleep(1 * time.Second)

		//2000 is shutdown for client

		pro <- Pac{-1, 2000, ""}

		var end1 = <-pro

		if end1.ack == 2001 {
			break
		}
	}

	fmt.Println("client has shutdown")
}

func ClientCheck(ackPointer *Pac, syn int, pro chan Pac) {
	temp := <-pro //Kill package: Pac{-1,-1,""}

	if temp.ack != -1 && temp.sec != -1 || syn != temp.sec { //Kill package actiaved
		*ackPointer = temp
	}

}

func Server(pro chan Pac) {

	temp := <-pro
	var syn = rand.IntN(100) + 100

	connectionEstablished := false

	for connectionEstablished == false {
		fmt.Println("Server, sends pac")
		pro <- Pac{temp.sec + 1, syn, temp.data}
		time.Sleep(200 * time.Millisecond)
		temp2 := <-pro
		fmt.Println("Server, recieve pac")
		if temp2.ack == syn+1 && temp2.sec == temp.sec+1 {
			fmt.Println("Server, sub activates")
			connectionEstablished = true
		}
	}

	fmt.Println("Connection fully established")

	var message [3]string

	message[0] = "I "
	message[1] = "am "
	message[2] = "God"

	var seq [3]int

	for i := 0; i < 3; i++ {
		seq[i] = rand.IntN(100) + 100*i
	}

	var serverReadyToShutDown bool = false

	for {
		if seq[0] != -1 {
			pro <- Pac{-1, seq[0], message[0]}
		}

		if seq[1] != -1 {
			pro <- Pac{-1, seq[1], message[1]}
		}

		if seq[2] != -1 {
			pro <- Pac{-1, seq[2], message[2]}
		}

		time.Sleep(500 * time.Millisecond)

		if seq[0] != -1 && seq[1] != -1 && seq[2] != -1 { //Checking to see if we recieve all acknoladgements.
			fmt.Printf("seq[0] = %d. seq[0] = %d. seq[0] = %d.\n",seq[0],seq[1],seq[2])
			var tmp0 = <-pro

			if tmp0.ack == seq[0]+1 {
				seq[0] = -1
			} else if tmp0.ack == seq[1]+1 {
				seq[1] = -1
			} else if tmp0.ack == seq[2]+1 {
				seq[2] = -1
			}
		}

		if seq[1] != -1 {
			var tmp1 = <-pro

			if tmp1.ack == seq[0]+1 {
				seq[0] = -1
			} else if tmp1.ack == seq[1]+1 {
				seq[1] = -1
			} else if tmp1.ack == seq[2]+1 {
				seq[2] = -1
			}
		}

		if seq[2] != -1 {
			var tmp2 = <-pro

			if tmp2.ack == seq[0]+1 {
				seq[0] = -1
			} else if tmp2.ack == seq[1]+1 {
				seq[1] = -1
			} else if tmp2.ack == seq[2]+1 {
				seq[2] = -1
			}
		}

		if seq[0] == -1 && seq[1] == -1 && seq[2] == -1 {
			fmt.Println("Server, ack for message recieved")

			if serverReadyToShutDown == false {
				pro <- Pac{-1, 1000, ""}

				var end0 = <-pro

				if end0.ack == 1001 {
					var end1 = <-pro
					pro <- Pac{end1.sec + 1, -1, ""}
					break
				} else {
					pro <- Pac{-1, 1000, ""}
				}
			}
		}
	}
	fmt.Println("server has shutdown")

}

type Pac struct {
	ack  int
	sec  int
	data string
}
