# tobiichi
Some personally accumulated Golang development tools.

## parallel

Simple concurrency toolkit with Golang.

### Line

Three steps complete a piece of concurrent code for processing streaming data, and terminate when `panic`/`error`/`context canceled`.
``` Golang
// Step 1: Set concurrency, the size is 2.
line, supplier := NewLine(2), make(chan interface{})

// Step 2: Set the execution action. params := <- supplier
line.Run(context.Background(), supplier, func(params interface{}) error {
    fmt.Print(params)
    return nil
})

// Step 3: Fill in the operating parameters. The close operation is necessary.
for i := 0; i < 10; i++ {
    supplier <- i
}
close(supplier)
```

If you want to process the return value, the following format is recommended.
``` Golang
// Step 1: Set concurrency, the size is 2. It is recommended to use channels to pass the return value.
line, supplier, result := NewLine(2), make(chan interface{}), make(chan interface{})

// Step 2: Set the execution action. params := <- supplier
line.Run(context.Background(), supplier, func(params interface{}) error {
    result <- params
    return nil
})

// Step 3: Fill in the operating parameters and wait for the process to end. The close operation is necessary.
go func() {
    defer func() {
        // Omitted panic recieve
        close(result)
    }()
    for i := 0; i < 10; i++ {
        supplier <- i
    }
    close(supplier)
    line.Wait()
}()

for item := range result {
    fmt.Print(item)
}
```