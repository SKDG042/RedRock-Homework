package main

import "Redrock/message-board/api"

func main() {
	h := api.InitRouter()
	h.Spin()
}
