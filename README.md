# hchannel

#### example

```go
package main

import (
	"fmt"

	"github.com/hqpko/hchannel"
)

func main() {
	c := hchannel.NewChannel(16, func(i interface{}) {
		// do something
	}).Run()
	defer c.Close()

	if ok := c.Input(1); ok {
		fmt.Println("input success")
	}
}

```

> 不使用 mustInput，会在 close 时造成 panic

> close 时基于 `range chan` 机制，会等待执行完所有已注入的消息
