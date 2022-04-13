package main

import (
	proxy "ReverseProxy"
)

func main() {
	proxy.Start("gommehd.net", 25565)
}
