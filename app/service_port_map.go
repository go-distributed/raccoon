package app

var ServicePortMap map[string]string

func init() {
	ServicePortMap = make(map[string]string)
	ServicePortMap["test service"] = ":8080"
}
