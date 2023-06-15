package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	pb "github.com/zhashkevych/trinity/internal/models" // replace with your actual path
)

func main() {
	// Connect to a NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to the 'calculated-prices' subject
	_, err = nc.Subscribe("calculated-prices", func(m *nats.Msg) {
		// Unmarshal the message data into a PoolPairList
		poolPairList := &pb.PoolPairList{}
		err = proto.Unmarshal(m.Data, poolPairList)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		// Print the PoolPairList
		for i, poolPair := range poolPairList.PoolPairs {
			fmt.Println(" ---------- ")
			fmt.Printf("PoolPair %d:\n", i+1)
			fmt.Printf("DexID: %v\n", poolPair.DexId)
			fmt.Printf("ID: %v\n", poolPair.Id)
			fmt.Printf("EffectivePrice0: %v\n", poolPair.EffectivePrice0)
			fmt.Printf("EffectivePrice1: %v\n", poolPair.EffectivePrice1)
			fmt.Printf("Reserve0: %v\n", poolPair.Reserve0)
			fmt.Printf("Reserve1: %v\n", poolPair.Reserve1)
			// Add more print statements for the other fields
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	// Keep the connection alive until Ctrl+C is pressed
	select {}
}
