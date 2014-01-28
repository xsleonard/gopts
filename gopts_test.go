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
    SetOptions(e, port())
    fmt.Printf("%d\n", e.Port)
    SetOptions(e, port(uint16(6666)))
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
    SetOptions(e, port())
    fmt.Printf("%d\n", e.Config.Port)
    SetOptions(e, port(uint16(6666)))
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
    SetOptions(e, port())
    fmt.Printf("%d\n", e.Config.Port)
    SetOptions(e, port(uint16(6666)))
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
    SetOptions(s, Port(), Network())
    // override defaults with user's values
    SetOptions(s, opts...)
    return s
}

func TestSetOptions(t *testing.T) {
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
    prevOpts := SetOptions(s, Port(uint16(6666)), Network("tcp"))
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

    // Sets a single option
    prevPort := SetOption(s, prevOpts[0])
    // s.Config.Port == 7777, back to default
    if s.Config.Port != 7777 {
        t.Fail()
    }

    // Reset to the previous value. These two statements are equivalent.
    SetOption(s, prevPort)
    // s.Config.Port == 6666
    if s.Config.Port != 6666 {
        t.Fail()
    }
    prevPort(s)
    // s.Config.Port == 6666
    if s.Config.Port != 6666 {
        t.Fail()
    }

    // Another example of avoid SetOption
    Network("udp")(s)
    if s.Network != "udp" {
        t.Fail()
    }
}
