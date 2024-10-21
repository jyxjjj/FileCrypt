package main

import (
    "fmt"
    "github.com/ccding/go-stun/stun"
)

func main() {

    client := stun.NewClient()
    client.SetServerHost("stun.miwifi.com", 3478)
    nat, host, err := client.Discover()
    if err != nil {
        fmt.Println("BehaviorTest failed:", err)
        return
    }
    fmt.Println("NAT Type:", nat.String())
    fmt.Println("Host:", host.String())
    behavior, err2 := client.BehaviorTest()
    if err2 != nil {
        fmt.Println("BehaviorTest failed:", err)
        return
    }
    fmt.Println("MappingType:", behavior.MappingType.String())
    fmt.Println("FilteringType:", behavior.FilteringType.String())
}
