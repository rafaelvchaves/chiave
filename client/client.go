package client


func main() {
	chiaveProxy := NewProxy(nil, 5, 3)
	chiaveProxy.Increment("key")
}