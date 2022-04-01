# cmd

Executing commands with verbose execution error output

### func RunOut(cmdCtx context.Context, param []string) ([]byte, error)
Returns + prints to standard input the result of executing the command passed as the 'param' argument.

### func Run(cmdCtx context.Context, param []string) ([]byte, error)
Returns the result of executing the command passed as the 'param' argument.

If the command fails, the Run() and RunOut() functions return in the error parameter a copy of the output result from the standard error stream.
See example below. 

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ds248a/cmd"
)

func main() {
	ctx, close := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	defer close()

	_, err := cmd.Run(ctx, []string{"ls", "-!", "."})
	fmt.Printf("err: %v\n", err)
/*
  remote out -> redirect to Err Message
------------------------

  cmd err
------------------------
err: ls: invalid option -- '!'
Try 'ls --help' for more information.
*/

	_, err := cmd.RunOut(ctx, []string{"ls", "-!", "."})
	fmt.Printf("err: %v\n", err)
/*
  remote out
------------------------
ls: invalid option -- '!'
Try 'ls --help' for more information.

  cmd err
------------------------
err: ls: invalid option -- '!'
Try 'ls --help' for more information.
*/
}
```
