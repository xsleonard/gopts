gopts
=====

Helper for defining configurable options for your structs

Documentation: [![GoDoc](https://godoc.org/github.com/xsleonard/gopts?status.png)](http://godoc.org/github.com/xsleonard/gopts)

## Example

```go
import "gopts"

var (
    Port = gopts.Opt("Config.Port", uint16(7777))
    Network = gopts.Opt("Network", "tcp")
    // ...
)

type ServerConfig struct {
    Port uint16
    // ...
}

type Server struct {
    Network string
    Config ServerConfig
    // ...
}

func NewServer(opts ...gopts.OptSetter) *Server {
    s := &Server{}
    // set defaults
    gopts.SetOptions(Port(), Network())
    // override defaults with user's values
    gopts.SetOptions(opts...)
    return s
}

func updateServer() {
    s := NewServer(Network("udp"))
    // s.Config.Port == 7777
    // s.Network == "udp"

    // Set multiple options
    prevOpts := gopts.SetOptions(s, Port(uint16(6666)), Network("tcp"))
    // s.Config.Port == 6666
    // s.Network == "tcp"

    // Sets a single option
    prevPort := gopts.SetOption(s, prevOpts[0])
    // s.Config.Port == 7777, back to default

    // Reset to the previous value. These two statements are equivalent.
    gopts.SetOption(s, prevPort)
    prevPort(s)
    // s.Config.Port == 6666
}
```
