# Logging at high frequency

The Gameboy clock runs at 4.194304 MHz, which means that the CPU executes 4,194,304 instructions per second which correspongs to 1 instruction every 238ns

The fmt.Print takes around 1.5us to execute, which is 6 times longer than the time it takes to execute an instruction. It can be measured by running the following code:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    start := time.Now()
    fmt.Println("Hello, World!")
    elapsed := time.Since(start)
    fmt.Printf("fmt.Println took %s\n", elapsed)
}

=== RUN   TestTick
Hello, World!
fmt.Println took 1.125µs
--- PASS: TestTick (0.11s)
PASS
ok      github.com/codefrite/gameboy-go/gameboy (cached)
```

This pretty much means that we should not log anything to the console during the execution of the CPU and probably the execution of all the components that operate at the same frequency as the CPU.

The theoritical limit would be 10^9 ns / 1.125 x 10^3 ns = 888,888 logs per second if the CPU does nothing else.
