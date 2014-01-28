/* Helpers for struct configuration.

Use it like so:

    import "gopts"

    var (
        Port = gopts.Option("Config.Port", uint16(7777))
        Network = gopts.Option("Network", "tcp")
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
*/
package gopts

import (
    "log"
    "reflect"
    "strings"
)

// Accepts 0 or 1 arguments.  If no arguments, the default value registered
// with Option is used.
type Opt func(...interface{}) OptSetter

// Sets a field's value on a wrapped struct.  interface{} must be an
// addressable type, i.e. a pointer to a struct.
type OptSetter func(interface{}) OptSetter

// Creates an Opt for field name and default value defVal.  Name must be
// capitalized since only public fields can be set with reflection.
// Nested field names are supported with "A.B.C" syntax.
func Option(name string, defVal interface{}) Opt {
    // Verify that all names will be exportable
    if len(name) == 0 {
        log.Panic("Option name cannot be empty")
    }
    fields := strings.Split(name, ".")
    for _, f := range fields {
        c := byte(f[0])
        if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
            log.Panicf("%s is not an exported field name", name)
        }
    }
    return newOpt(name, defVal)
}

func newOpt(name string, defVal interface{}) Opt {
    return func(_val ...interface{}) OptSetter {
        // Unpack the arguments or use the default value
        var val reflect.Value
        if len(_val) == 0 {
            val = reflect.ValueOf(defVal)
        } else if len(_val) != 1 {
            log.Panic("Multivariable options not supported")
        } else {
            val = reflect.ValueOf(_val[0])
        }
        return newOptSetter(name, val)
    }
}

func newOptSetter(name string, val reflect.Value) OptSetter {
    fields := strings.Split(name, ".")
    return func(self interface{}) OptSetter {
        obj := reflect.ValueOf(self)
        for _, f := range fields[:len(fields)-1] {
            obj = reflect.Indirect(obj).FieldByName(f)
            if !obj.IsValid() {
                log.Panicf("Unknow field %s", f)
            }
        }
        // Use reflection to set the name on the self struct
        // self must be a pointer to a struct
        fieldName := fields[len(fields)-1]
        field := reflect.Indirect(obj).FieldByName(fieldName)
        if !field.IsValid() {
            log.Panicf("Unknown field %s", name)
        }
        if !field.CanSet() {
            log.Panicf("Can't set field on Option target, use a pointer")
        }
        prev := reflect.ValueOf(field.Interface())
        field.Set(val)
        return newOptSetter(name, prev)
    }
}

// Sets multiple Options on a struct pointer.  Returns the previous Options'
// values in order they were passed to this function.
func Set(obj interface{}, opts ...OptSetter) []OptSetter {
    prevs := make([]OptSetter, 0, len(opts))
    for _, o := range opts {
        prevs = append(prevs, o(obj))
    }
    return prevs
}
