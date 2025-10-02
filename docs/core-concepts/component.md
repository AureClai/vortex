## Components structs

With an exemple Component name Example.

### The main struct : Example

The main struct is the struct that defines the Go object for the Example component. In good Vortex practice, these should not be used to store states (as it does not retrigger render) and to pass info from parent (use props instead)

```Go
type Example struct {
    BaseComponant // Type of baseComponent Base, Stateful, etc
    state ExampleState // if stateful
    props ExampleProps
    arg1 interface{}
    arg2 interface{}
}
```

The animation objects should be in that struct.

### The props struct : ExampleProps

```Go
//go:vortex-props
type ExampleProps struct {
    prop1 prop1Type
}
```

Vortex developpment tool `vortex dev` create a builder for the ExampleComponent in real time when it detects a Props struct. The declaration `go:vortex-props` aims to make the Go AST detect those props and build the functions :

```
// file : example.vtb.go
// Automatically generated builder

type ExempleBuilder struct {
    props ExampleProps
}

func Build_Example() *ExempleBuilder {
    return &ExempleBuilder{
        props : make(ExampleProps)
    }
}

func (e ExempleBuilder) Prop1(prop1 prop1Type) *ExempleBuilder {
    e.props.prop1 = prop1
}
```

### The state struct

```

The state struct represents the inner data structure of the component than trigger a reRender on update.

```
