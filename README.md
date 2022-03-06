# closer

### Основной функционал

## func RunOut(cmdCtx context.Context, param []string) ([]byte, error)
Возвращает + выводит в стандартный поток ввода результат выполнения комманды,
переданной в качестве аргумента param.

## func Run(cmdCtx context.Context, param []string) ([]byte, error)
Возвращает результат выполнения комманды, переданной в качестве аргумента param.

В случае не удачного выполнения команды, финкции Run и RunOut возвращают в параметре error копию результата
вывода из стандартного потока ошибок. Смотрите пример ниже.

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

	out, err := cmd.RunOut(ctx, []string{"ls", "-!", "."})

	fmt.Printf("err: %v\n", err)
	fmt.Printf("out: %s\n", string(out))

	// Результат выполнения
	// err: ls: invalid option -- '!'
	// Try 'ls --help' for more information.

}

```

### Объектная и кроспакетная обработка

Пример выше - образец кроспакетной обработки. Следующий вариант - объектный.

```go
import (
  "github.com/ds248a/closer"
)

func main() {
  c := closer.NewCloser()
  c.Add(cc.Clear)
  c.Add(pg.Close)
  c.ListenSignal(syscall.SIGTERM)
}
```
