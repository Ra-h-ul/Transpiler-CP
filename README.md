//input 1 

package main

import (
    "fmt"
)

func addsub(x int) (a, b int) {
    return x + 3, x - 2
}

func main() {
    y, z := addsub(4)
    fmt.Println("y =", y)
    fmt.Println("z =", z)
}






// input 2

package main

import (
	"fmt"
)

func main() {
	m := map[string]string{"first": "hi", "second": "you", "third": "there"}
	first := true
	for k, v := range m {
		if first {
			first = false
		} else {
			fmt.Print(" ")
		}
		fmt.Print(k + v)
	}
	fmt.Println()
}

}




## Standard library

- [x] `fmt.Println`
- [x] `fmt.Print`
- [ ] `fmt.Printf` (partially)
- [ ] `fmt.Sprintf`
- [x] `strings.Contains`
- [x] `strings.HasPrefix`
- [ ] `strings.HasSuffix`
- [ ] `strings.Index`
- [ ] `strings.Join`
- [ ] `strings.NewReader`
- [ ] `strings.Replace`
- [ ] `strings.Split`
- [ ] `strings.SplitN`
- [x] `strings.TrimSpace`
