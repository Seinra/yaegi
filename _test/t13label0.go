package main

import "fmt"

// Test for labels defined in switch case bodies (caseBody): `next:`
// must resolve from the nested switch. gc prints 0-true-"i".
func scan(s string) string {
	format := s
	nvalues := 0
	ok := false
out:
	for {
		var c byte
		if len(format) == 0 {
			c = 0
		} else {
			c = format[0]
			format = format[1:]
		}
		switch c {
		case '%':
			nvalues++
			ok = true
		case 0:
			break out
		case ' ', '\t':
			format = format
			ok = true
		next:
			for {
				if len(format) == 0 {
					break out
				}
				c2 := format[0]
				switch c2 {
				case ' ', '\t':
					// nop
				default:
					break next
				}
				format = format[1:]
			}
		default:
			break out
		}
	}
	return fmt.Sprintf("%d-%v-%q", nvalues, ok, format)
}

func main() { fmt.Println(scan("   hi")) }

// Output:
// 0-true-"i"
