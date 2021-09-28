# Introduction

greeter is a demonstration of grpc service.


# Benchmark

```
goos: linux
goarch: amd64
pkg: github.com/LiangXianSen/greeter/hello
BenchmarkHelloGRPC-8    73838512                16.1 ns/op             0 B/op          0 allocs/op
BenchmarkHelloGRPC-8    74298560                15.9 ns/op             0 B/op          0 allocs/op
BenchmarkHelloGRPC-8    74460302                16.1 ns/op             0 B/op          0 allocs/op
BenchmarkHelloHTTP-8       12663             95766 ns/op           16943 B/op        252 allocs/op
BenchmarkHelloHTTP-8       12612             95216 ns/op           16942 B/op        252 allocs/op
BenchmarkHelloHTTP-8       12627             94939 ns/op           16942 B/op        252 allocs/op
PASS
ok      github.com/LiangXianSen/greeter/hello   11.630s
?       github.com/LiangXianSen/greeter/prometheus      [no test files]
2021/09/28 17:52:02 HTTP service listen on :8000
2021/09/28 17:52:02 GRPC service listen on :8080
goos: linux
goarch: amd64
pkg: github.com/LiangXianSen/greeter/test/hello
BenchmarkHelloGRPC-8       14277             77401 ns/op           11788 B/op        198 allocs/op
BenchmarkHelloGRPC-8       15445             77374 ns/op           11709 B/op        198 allocs/op
BenchmarkHelloGRPC-8       15496             77018 ns/op           11708 B/op        198 allocs/op
BenchmarkHelloHTTP-8        2965            524253 ns/op           45908 B/op        408 allocs/op
BenchmarkHelloHTTP-8        5061            263300 ns/op           45549 B/op        407 allocs/op
BenchmarkHelloHTTP-8        4798            246794 ns/op           45407 B/op        407 allocs/op
PASS
ok      github.com/LiangXianSen/greeter/test/hello      11.778s
```