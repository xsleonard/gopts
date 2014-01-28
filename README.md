gopts
=====

Helper for defining configurable options for your structs.

This package was inspired by http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html

Documentation: [![GoDoc](https://godoc.org/github.com/xsleonard/gopts?status.png)](http://godoc.org/github.com/xsleonard/gopts)

## Example

```go
import "gopts"

var (
    Port = gopts.Option("Config.Port", uint16(7777))
    Network = gopts.Option("Network", "tcp")
)

type ServerConfig struct {
    Port uint16
}

type Server struct {
    Network string
    Config ServerConfig
}

func NewServer(opts ...gopts.OptSetter) *Server {
    s := &Server{}
    // set defaults
    gopts.Set(s, Port(), Network())
    // override defaults with user's values
    gopts.Set(s, opts...)
    return s
}

func updateServer() {
    s := NewServer(Network("udp"))
    // s.Config.Port == 7777
    // s.Network == "udp"

    // Set multiple options
    prevOpts := gopts.Set(s, Port(uint16(6666)), Network("tcp"))
    // s.Config.Port == 6666
    // s.Network == "tcp"

    // Reset to the previous values. These two statements are equivalent.
    prevOpts = gopts.Set(s, prevOpts...)
    // s.Config.Port == 7777
    // s.Config.Network == "udp"

    // Shorthand for setting a single option:
    prevOpts[0](s)
    // s.Config.Port == 6666
}
```
