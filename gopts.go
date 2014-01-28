// Helpers for struct configuration
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
type OptSetter func(interface{})

// Creates an Opt for field name and default value defVal.  Name must be
// capitalized since only public fields can be set with reflection.
// Nested field names are supported.  Nested fields
func Option(name string, defVal interface{}) Opt {
    // Pre-verify that all names will be exportable
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

        return func(self interface{}) {
            obj := reflect.ValueOf(self)
            for _, f := range fields[:len(fields)-1] {
                obj = reflect.Indirect(obj).FieldByName(f)
                if !obj.IsValid() {
                    log.Panicf("Unknow field %s", f)
                }
            }
            // Use reflection to set the name on the self struct
            // self must be a pointer to a struct
            name := fields[len(fields)-1]
            field := reflect.Indirect(obj).FieldByName(name)
            if !field.IsValid() {
                log.Panicf("Unknown field %s", name)
            }
            if !field.CanSet() {
                log.Panicf("Can't set field on Option target, use a pointer")
            }
            field.Set(val)
        }
    }
}

// Sets multiple Options on a struct pointer
func SetOptions(obj interface{}, opts ...OptSetter) {
    for _, o := range opts {
        o(obj)
    }
}
