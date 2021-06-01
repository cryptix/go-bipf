package bipf_test

import (
	"encoding/hex"
	"math"
	"os"

	"github.com/ssb-ngi-pointer/go-bipf"
)

func ExampleMapOf() {
	theMap := bipf.MapOf(map[string]bipf.Valuer{
		"i1": bipf.Int32(1337),
		"s1": bipf.String("acab"),
		"d1": bipf.Double(23.42),
		"b1": bipf.Bool(true),
		"b2": bipf.Bool(false),
	}, "i1", "s1", "d1", "b1", "b2")

	hexd := hex.Dumper(os.Stdout)
	if err := theMap(hexd); err != nil {
		panic(err)
	}

	// Output:
	// 00000000  35 10 69 31 22 39 05 00  00 10 73 31 20 61 63 61  |5.i1"9....s1 aca|
	// 00000010  62 10 64 31 43 ec 51 b8  1e 85 6b 37 40 10 62 31  |b.d1C.Q...k7@.b1|
	// 00000020  0e 01 10 62 32 0e 00
}

func ExampleListOf() {
	theList := bipf.ListOf(
		bipf.Int32(1337),
		bipf.String("acab"),
		bipf.Double(math.NaN()),
		bipf.Double(23.42),
		bipf.Bool(false),
		bipf.Double(-0.001),
		bipf.Bool(true),
	)

	hexd := hex.Dumper(os.Stdout)
	if err := theList(hexd); err != nil {
		panic(err)
	}

	// Output:
	// 00000000  4c 22 39 05 00 00 20 61  63 61 62 43 01 00 00 00  |L"9... acabC....|
	// 00000010  00 00 f8 7f 43 ec 51 b8  1e 85 6b 37 40 0e 00 43  |....C.Q...k7@..C|
	// 00000020  fc a9 f1 d2 4d 62 50 bf  0e 01
}
