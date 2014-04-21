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
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-distributed/raccoon/app"
	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/instance"
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

	da := app.NewDemoApp(c)
	c.AddListener(controller.AddRouterEventType, da.AddRouterListener)
	c.AddListener(controller.AddInstanceEventType, da.AddInstanceListener)

	return nil
}

func plotRouter() error {
	if len(os.Args) < 5 {
		return fmt.Errorf("Usage: demo r <cAddr> <port> <id>")
	}

	addr, err := getInterfaceAddr()
	if err != nil {
		return err
	}

	cAddr := os.Args[2]
	port := os.Args[3]
	rAddr := addr + port
	id := os.Args[4]

	if id != "rand" && id != "rr" {
		return fmt.Errorf("for the sake of this demo, only ids of 'rand' and 'rr' are supported: %v", id)
	}

	// start router
	r, err := router.New(id, rAddr, cAddr)
	if err != nil {
		return err
	}

	err = r.Start()
	if err != nil {
		return err
	}

	err = r.RegisterOnCtler()

	return err
}

// Service: "test service"
func plotInstance() error {
	if len(os.Args) < 4 {
		return fmt.Errorf("Usage: demo i <cAddr> <id>")
	}

	cAddr := os.Args[2]
	id := os.Args[3]
	service := "test service"

	// start http server
	addr, err := getInterfaceAddr()
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, world %v\n", id)
	})

	l, err := net.Listen("tcp", addr+":0")
	if err != nil {
		return err
	}
	go http.Serve(l, nil)

	iAddr := l.Addr().String()
	//fmt.Println("http address:", iAddr)

	// create instance
	instance, err := instance.NewInstance(id, service, iAddr)
	if err != nil {
		return err
	}

	// register instance to controller

	err = instance.RegisterOnCtler(cAddr)
	if err != nil {
		return err
	}

	resp, err := http.Get("http://" + iAddr + "/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(reply))

	return err
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
