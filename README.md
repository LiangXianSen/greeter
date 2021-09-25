# Introduction

greeter is a demonstration of grpc service.


# Benchmark

```
goos: linux
goarch: amd64
pkg: github.com/LiangXianSen/greeter/hello
BenchmarkHelloGRPC-8    72829242                16.9 ns/op             0 B/op          0 allocs/op
BenchmarkHelloGRPC-8    72144979                16.2 ns/op             0 B/op          0 allocs/op
BenchmarkHelloGRPC-8    72489790                16.5 ns/op             0 B/op          0 allocs/op
BenchmarkHelloHTTP-8      298944              4125 ns/op            3492 B/op         51 allocs/op
BenchmarkHelloHTTP-8      312489              4248 ns/op            3490 B/op         51 allocs/op
BenchmarkHelloHTTP-8      290841              3866 ns/op            3493 B/op         51 allocs/op
PASS
ok      github.com/LiangXianSen/greeter/hello   8.309s
goos: linux
goarch: amd64
pkg: github.com/LiangXianSen/greeter/test/hello
BenchmarkHelloGRPC-8       15736             77408 ns/op           11526 B/op        194 allocs/op
BenchmarkHelloGRPC-8       15768             78378 ns/op           11526 B/op        194 allocs/op
BenchmarkHelloGRPC-8       14536             75797 ns/op           11527 B/op        194 allocs/op
BenchmarkHelloHTTP-8        3439            314545 ns/op           32162 B/op        203 allocs/op
BenchmarkHelloHTTP-8        8484            215427 ns/op           32056 B/op        203 allocs/op
BenchmarkHelloHTTP-8        7856           1718799 ns/op           32107 B/op        203 allocs/op
PASS
ok      github.com/LiangXianSen/greeter/test/hello      22.873s
```