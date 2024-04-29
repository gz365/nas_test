package main

import (
	"fmt"
	"smb_gbox/client"
	"smb_gbox/conf"
)

func main() {

	c, err := client.NewClient(conf.GetConf())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// c.Upload()
	c.ReadDir()

	defer c.Disconnect()
}
