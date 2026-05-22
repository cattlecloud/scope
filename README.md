# scope

[![Go Reference](https://pkg.go.dev/badge/cattlecloud.net/go/scope.svg)](https://pkg.go.dev/cattlecloud.net/go/scope)
[![License](https://img.shields.io/github/license/cattlecloud/scope?color=7C00D8&style=flat-square&label=License)](https://github.com/cattlecloud/scope/blob/main/LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/cattlecloud/scope/ci.yaml?style=flat-square&color=0FAA07&label=Tests)](https://github.com/cattlecloud/scope/actions/workflows/ci.yaml)

`scope` is a substitute for the Go context package.

It provides a more convenient context API than the standard library, while
maintaining 100% compatibility.

### Requirements

The minimum Go version is `go1.26`.

### Install

The `scope` package can be added to a project with `go get`.

```shell
go get -u cattlecloud.net/go/scope@latest
```

### Examples

##### New

```go
ctx := scope.New()
```

##### TTL

```go
ctx, cancel := scope.TTL(5 * time.Second)
// ctx is canceled after 5 seconds
defer cancel()
```

##### Deadline

```go
ctx, cancel := scope.Deadline(time.Now().Add(10 * time.Second))
// ctx is canceled at the specified time
defer cancel()
```

##### Cancelable

```go
ctx, cancel := scope.Cancelable()
// ctx can be canceled manually
defer cancel()
```

##### WithCancel

```go
ctx, cancel := scope.WithCancel(parentCtx)
defer cancel()
```

##### WithTTL

```go
ctx, cancel := scope.WithTTL(parentCtx, 3 * time.Second)
// parentCtx with a 3 second timeout
defer cancel()
```

##### WithValue

```go
ctx := scope.WithValue(parentCtx, "userID", 123)
```

##### Value

```go
userID := scope.Value[int](ctx, "userID")
```

##### Join

```go
ctx1, cancel1 := scope.WithCancel(scope.New())
ctx2, cancel2 := scope.TTL(5 * time.Second)

joined, cancel := scope.Join(ctx1, ctx2)
// joined is canceled when ctx1 or ctx2 is canceled
defer cancel()
defer cancel1()
defer cancel2()
```

###### Deadline

```go
ctx1, _ := scope.Deadline(time.Now().Add(10 * time.Second))
ctx2, _ := scope.Deadline(time.Now().Add(20 * time.Second))

joined, _ := scope.Join(ctx1, ctx2)
deadline, ok := joined.Deadline() // deadline is 10 seconds, ok is true
```

###### Done

```go
joined, cancel := scope.Join(ctx1, ctx2)
<-joined.Done() // blocks until either ctx1 or ctx2 is done
```

###### Err

```go
joined, cancel := scope.Join(ctx1, ctx2)
<-joined.Done()
err := joined.Err() // returns the error from the first canceled context
```

###### Value

```go
ctx1 := scope.WithValue(scope.New(), "key", "value1")
ctx2 := scope.WithValue(scope.New(), "key", "value2")

joined, _ := scope.Join(ctx1, ctx2)
val := joined.Value("key") // returns value1 (ctx1's value is checked first)
```

### License

The `cattlecloud.net/go/scope` module is open source under the [BSD](LICENSE) license.
