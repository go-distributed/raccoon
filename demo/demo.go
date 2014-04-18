// Identities:
// - One controller
// - Two routers
// - Two service instances
//
// Demo scenarios:
//
// 1. Add one router for round robin.
// 2. Add another router for random select. Compare it with the first one.
// 3. Instance failure. Both routers must report the failure to controller
//    and then it will remove that instance from serving.
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/router"
)

func plotController() error {
	cAddr := os.Args[2]
	c, err := controller.New(cAddr)
	if err != nil {
		return err
	}

	err = c.Start()
	if err != nil {
		return err
	}

	return nil
}

func plotRouter() error {
	if len(os.Args) < 5 {
		return fmt.Errorf("Usage: demo r <cAddr> <rAddr> id")
	}

	cAddr := os.Args[2]
	rAddr := os.Args[3]
	id := os.Args[4]

	// start router
	r, err := router.New(rAddr)
	if err != nil {
		return err
	}

	err = r.Start()
	if err != nil {
		return err
	}

	// register router in controller
	regRouterArgs := &controller.RegRouterArgs{
		Id: id,
	}

	addr, err := getInterfaceAddr()
	if err != nil {
		return err
	}
	regRouterArgs.Addr = addr + rAddr
	//regRouterArgs.Addr = "127.0.0.1" + rAddr

	client, err := rpc.Dial("tcp", cAddr)
	if err != nil {
		return err
	}

	//fmt.Println("debug:", regRouterArgs.Addr, cAddr)

	err = client.Call("ControllerRPC.RegisterRouter", regRouterArgs, nil)
	if err != nil {
		return err
	}

	return nil
}

func plotInstance() error {
	panic("")
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: demo [c|r|i] <cAddr>")
	}

	switch os.Args[1] {
	case "c":
		err := plotController()
		if err != nil {
			log.Fatal("plotController:", err)
		}
	case "r":
		err := plotRouter()
		if err != nil {
			log.Fatal("plotRouter:", err)
		}
	case "i":
		err := plotInstance()
		if err != nil {
			log.Fatal("plotInstance:", err)
		}
	default:
		log.Fatal("Usage: demo [c|r|i] <cAddr>")
	}

	log.Println(os.Args[1], "successfully running")
	select {}
}

func getInterfaceAddr() (string, error) {

	intAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	var addr string
	for _, iAddr := range intAddrs {
		if !strings.HasPrefix(iAddr.String(), "127.") &&
			!strings.HasPrefix(iAddr.String(), "172.") {
			addr = iAddr.String()
			break
		}
	}

	if addr == "" {
		return "", fmt.Errorf("cannot found any addr: %v", intAddrs)
	}

	index := strings.Index(addr, "/")
	if index != -1 {
		return addr[:index], nil
	} else {
		return addr, nil
	}
}
