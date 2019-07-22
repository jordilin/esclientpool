# esclientpool

esclientpool maintains a pool of Elasticsearch clients that can be
used in scenarios where you need to use concurrency such as having an
HTTP endpoint that forwards requests to Elasticsearch.

## Building and testing

```bash
go build github.com/jordilin/esclientpool
```

```bash
go test github.com/jordilin/esclientpool
```

## License

This project is licensed under

* BSD-3-Clause ([LICENSE](LICENSE) or [https://opensource.org/licenses/BSD-3-Clause](https://opensource.org/licenses/BSD-3-Clause))
