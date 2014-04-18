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

func main() {

}
