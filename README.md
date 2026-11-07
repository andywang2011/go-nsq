# go-nsq

Thin NSQ producer wrapper (package `nsq`) around [go-nsq](https://github.com/nsqio/go-nsq).

## Install

```bash
go get github.com/andywang2011/go-nsq@v0.1.0
```

## Quick start

```go
import (
	"context"
	"fmt"

	"github.com/andywang2011/go-nsq"
)

cfg := nsq.ConfigDefault()
cfg.Address = "127.0.0.1:65530"

p, err := nsq.NewProducer(context.Background(), cfg)
if err != nil {
	fmt.Println("unavailable")
	return
}
defer p.Stop()
fmt.Println("ok")
```

## Examples

```bash
export NSQD_ADDR=127.0.0.1:4150
cd go-nsq/examples/producer && go run .
cd go-nsq/examples/consumer && go run .
```

## License

MIT — see [LICENSE](LICENSE).
