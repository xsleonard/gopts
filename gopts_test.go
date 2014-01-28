package gopts

import (
    "fmt"
    "testing"
)

func ExampleOptions() {
    port := Option("Port", uint16(7777))
    e := &struct{ Port uint16 }{
        uint16(8888),
    }
    fmt.Printf("%d\n", e.Port)
    Set(e, port())
    fmt.Printf("%d\n", e.Port)
    Set(e, port(uint16(6666)))
    fmt.Printf("%d\n", e.Port)
    // Output:
    // 8888
    // 7777
    // 6666
}

func ExampleNestedOptions() {
    port := Option("Config.Port", uint16(7777))
    e := &struct{ Config struct{ Port uint16 } }{
        struct{ Port uint16 }{
            uint16(8888),
        },
    }
    fmt.Printf("%d\n", e.Config.Port)
    Set(e, port())
    fmt.Printf("%d\n", e.Config.Port)
    Set(e, port(uint16(6666)))
    fmt.Printf("%d\n", e.Config.Port)
    // Output:
    // 8888
    // 7777
    // 6666
}

func ExampleNestedOptionsWithPointer() {
    port := Option("Config.Port", uint16(7777))
    e := &struct{ Config *struct{ Port uint16 } }{
        &struct{ Port uint16 }{
            uint16(8888),
        },
    }
    fmt.Printf("%d\n", e.Config.Port)
    Set(e, port())
    fmt.Printf("%d\n", e.Config.Port)
    Set(e, port(uint16(6666)))
    fmt.Printf("%d\n", e.Config.Port)
    // Output:
    // 8888
    // 7777
    // 6666
}

var (
    Port    = Option("Config.Port", uint16(7777))
    Network = Option("Network", "tcp")
    // ...
)

type ServerConfig struct {
    Port uint16
    // ...
}

type Server struct {
    Network string
    Config  ServerConfig
    // ...
}

func NewServer(opts ...OptSetter) *Server {
    s := &Server{}
    // set defaults
    Set(s, Port(), Network())
    // override defaults with user's values
    Set(s, opts...)
    return s
}

func TestSet(t *testing.T) {
    s := NewServer(Network("udp"))
    // s.Config.Port == 7777
    // s.Network == "udp"
    if s.Config.Port != 7777 {
        t.Fail()
    }
    if s.Network != "udp" {
        t.Fail()
    }

    // Set multiple options
    prevOpts := Set(s, Port(uint16(6666)), Network("tcp"))
    // s.Config.Port == 6666
    // s.Network == "tcp"

    if s.Config.Port != 6666 {
        t.Fail()
    }
    if s.Network != "tcp" {
        t.Fail()
    }
    if len(prevOpts) != 2 {
        t.Fail()
    }

    // Reset to the previous value. These two statements are equivalent.
    prevPort := Set(s, prevOpts...)[0]
    // s.Config.Port == 7777
    if s.Config.Port != 7777 {
        t.Fail()
    }
    if s.Network != "udp" {
        t.Fail()
    }

    // Setting a single value
    prevPort(s)
    // s.Config.Port == 6666
    if s.Config.Port != 6666 {
        t.Fail()
    }

    // Another example for setting a single value
    Network("udp")(s)
    if s.Network != "udp" {
        t.Fail()
    }
}
