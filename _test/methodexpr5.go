package main

import "fmt"

// Deep promotion chain: 70 single-field embeds exercise the
// walker at depths a fixed visit cap would truncate (native
// R3-CAP64). gc binds method expressions at any depth; the
// interpreter must terminate and agree.
type leaf struct{}

func (leaf) get() int { return 7 }

type L1 struct{ leaf }
type L2 struct{ L1 }
type L3 struct{ L2 }
type L4 struct{ L3 }
type L5 struct{ L4 }
type L6 struct{ L5 }
type L7 struct{ L6 }
type L8 struct{ L7 }
type L9 struct{ L8 }
type L10 struct{ L9 }
type L11 struct{ L10 }
type L12 struct{ L11 }
type L13 struct{ L12 }
type L14 struct{ L13 }
type L15 struct{ L14 }
type L16 struct{ L15 }
type L17 struct{ L16 }
type L18 struct{ L17 }
type L19 struct{ L18 }
type L20 struct{ L19 }
type L21 struct{ L20 }
type L22 struct{ L21 }
type L23 struct{ L22 }
type L24 struct{ L23 }
type L25 struct{ L24 }
type L26 struct{ L25 }
type L27 struct{ L26 }
type L28 struct{ L27 }
type L29 struct{ L28 }
type L30 struct{ L29 }
type L31 struct{ L30 }
type L32 struct{ L31 }
type L33 struct{ L32 }
type L34 struct{ L33 }
type L35 struct{ L34 }
type L36 struct{ L35 }
type L37 struct{ L36 }
type L38 struct{ L37 }
type L39 struct{ L38 }
type L40 struct{ L39 }
type L41 struct{ L40 }
type L42 struct{ L41 }
type L43 struct{ L42 }
type L44 struct{ L43 }
type L45 struct{ L44 }
type L46 struct{ L45 }
type L47 struct{ L46 }
type L48 struct{ L47 }
type L49 struct{ L48 }
type L50 struct{ L49 }
type L51 struct{ L50 }
type L52 struct{ L51 }
type L53 struct{ L52 }
type L54 struct{ L53 }
type L55 struct{ L54 }
type L56 struct{ L55 }
type L57 struct{ L56 }
type L58 struct{ L57 }
type L59 struct{ L58 }
type L60 struct{ L59 }
type L61 struct{ L60 }
type L62 struct{ L61 }
type L63 struct{ L62 }
type L64 struct{ L63 }
type L65 struct{ L64 }
type L66 struct{ L65 }
type L67 struct{ L66 }
type L68 struct{ L67 }
type L69 struct{ L68 }
type L70 struct{ L69 }

func main() {
	o := L70{}
	f := (*L70).get
	fmt.Println(f(&o))
}

// Output:
// 7
