package gopts

import "fmt"

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
