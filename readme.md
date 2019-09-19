# Generator

[![GoDoc](https://godoc.org/github.com/bmdelacruz/generator?status.svg)](https://godoc.org/github.com/bmdelacruz/generator)

A simple Go package which provides utility for creating generic generators like the one from Javascript. It uses channels to send and receive data to and from the generator function.

Please take note that its behaviour is not a complete duplicate of Javascript. It's implemented as a language feature in Javascript (They have a dedicated keyword, `yield`, for yielding values and syntax, `function*`, for marking functions as generators) and error handling in Go completely different from Javascript's.


## Sample usage

```go
g := generator.New(
  func(gc *generator.Controller) (interface{}, error) {
    for i := 0; i < 5; i++ {
      gc.Yield(i)
    }
    return nil, nil
  },
)
for value, isDone, _ := g.Next(nil); !isDone; value, isDone, _ = g.Next(nil) {
  fmt.Println(value)
}

// Output:
// 0
// 1
// 2
// 3
// 4
```

## Author
Created by Bryan Dela Cruz &lt;bryanmdlx@gmail.com&gt;

## Questions

Do you have questions regarding this package? Please file an issue with a `question` label.

## Contributors

None yet!

### How to contribute

Discovered a bug in this package? Want to change generator's behaviour? Have you thought of a more efficient way to implement the generator? File an issue with appropriate label/s and let's talk about it!

Already forked this repository to add the stuff you want? Pull requests are very much appreciated!