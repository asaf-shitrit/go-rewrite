[![Project Tests](https://github.com/asaf-shitrit/html-overwrite/actions/workflows/go.yml/badge.svg)](https://github.com/asaf-shitrit/html-overwrite/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/asaf-shitrit/html-overwrite/branch/main/graph/badge.svg?token=4BAB8KMGCJ)](https://codecov.io/gh/asaf-shitrit/html-overwrite)
<h2 align="center">Simple HTML Modification in Go</h2>
Do you grin at the sight of html.Node ? Me too.

Modifying HTML in Go should be simple.


**🧘🏻 Human friendly**: query language for humans that takes a minute to grasp.

**🎓 Simplicity**: no robust solution here just one single goal

**🏎 Fast**: has support for zero allocation Set/Append

## Getting Started

```go
package main
import "github.com/asaf-shitrit/go-rewrite"

func main(){

    example := `
       <html>
            <body>
                <h1 id="content">Hello !</h1>
            </body>
       </html>
    `

    // load
    doc := html_overwrite.Load(strings.NewReader(example))

    // mutate
    doc.Set("id=content", "Bye Bye")

    // get back
    newStringDoc := doc.String()
}
```

## Fast Set/Append

```go
import "github.com/asaf-shitrit/go-rewrite"

func main(){
    
    // example stream
    res, err := http.DefaultClient.Get("somesite.com")
    output := &bytes.Buffer{}

    // append
    html_overwrite.Append(res.Body, output, "tag=head", injectedValue)
}
```
## Query Language
```

The two html identifiers supported right now are:
- id
- class name


```
// All matchers are split up by a comma (,)

// Pattern:
[matcher]=[value],[matcher]=[value]

// Example (By ID):
id=content

// Example (By Class):
class=great-name

// Example (Multiple Matchers)
id=content,class=great-name

## Stream Set & Append Benchmarks

Useful for stream cases where a single 
performant mutation is required.

std go lib:
```
BenchmarkAppend
BenchmarkAppend/std_lib
BenchmarkAppend/std_lib         	  153255	      7917 ns/op	   12232 B/op	      69 allocs/op
BenchmarkSet
BenchmarkSet/std_lib
BenchmarkSet/std_lib            	  153068	      7913 ns/op	   12512 B/op	      82 allocs/op
```

stream implementation (zero allocations):
 ```
BenchmarkAppend
BenchmarkAppend/stream
BenchmarkAppend/stream          	  398082	      3024 ns/op	       0 B/op	       0 allocs/op
BenchmarkSet
BenchmarkSet/stream
BenchmarkSet/stream             	  409041	      2886 ns/op	       0 B/op	       0 allocs/op
PASS
```


## State
- Will be actively developed by me based on features I require in personal projects
- Tests will be added to each new feature release 
- Feel free to create new issues if you find any 😃
