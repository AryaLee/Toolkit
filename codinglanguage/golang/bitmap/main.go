package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/kelindar/bitmap"
)

func main() {
	t := bitmap.Bitmap{}
	v, z := t.Max()
	fmt.Println("max: ", v, z, t.Count())
	a, x, b := net.ParseCIDR("100.64.0.10/24")
	fmt.Println(a, x, b)
	c, d := x.Mask.Size()
	fmt.Println(x.IP, x.Mask)
	fmt.Println(c, d)

	t.Set(0)
	t.Set(1<<(d-c) - 1)

	fmt.Println(t.ToBytes())
	fmt.Println(json.Marshal(t))
	v, z = t.Max()
	fmt.Println("max: ", v, z, t.Count())
	os.Exit(0)
	t.Range(func(x uint32) {
		fmt.Println("range x", x)
		t.Remove(x)
	})

	t.Range(func(x uint32) {
		fmt.Println("range2 x", x)
	})
	fmt.Println(t.ToBytes())

	v, z = t.Max()
	fmt.Println("max: ", v, z, t.Count())
	p := bitmap.Bitmap{}
	t.Clone(&p)
	p.Ones()
	p.Xor(t)
	fmt.Println(cmp.Equal(p, t))
	v, z = p.Max()
	fmt.Println("max: ", v, z, p.Count())
	v, z = t.Min()
	fmt.Println("t min: ", v, z, p.Count())
	t.Remove(0)
	v, z = t.Min()
	fmt.Println("t min: ", v, z, p.Count())
}
