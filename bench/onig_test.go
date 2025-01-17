package bench

import (
	"github.com/tmikus/onig-go/v2"
	"strings"
	"testing"
)

func BenchmarkOnigLiteral(b *testing.B) {
	x := strings.Repeat("x", 50) + "y"
	b.StopTimer()
	re := onig.MustNewRegex("y")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if re.MustFindMatch(x) == nil {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkOnigNotLiteral(b *testing.B) {
	x := strings.Repeat("x", 50) + "y"
	b.StopTimer()
	re := onig.MustNewRegex(".y")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if re.MustFindMatch(x) == nil {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkOnigMatchClass(b *testing.B) {
	b.StopTimer()
	x := strings.Repeat("xxxx", 20) + "w"
	re := onig.MustNewRegex("[abcdw]")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if re.MustFindMatch(x) == nil {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkOnigMatchClass_InRange(b *testing.B) {
	b.StopTimer()
	// 'b' is between 'a' and 'c', so the charclass
	// range checking is no help here.
	x := strings.Repeat("bbbb", 20) + "c"
	re := onig.MustNewRegex("[ac]")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if re.MustFindMatch(x) == nil {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkOnigAnchoredLiteralShortNonMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := onig.MustNewRegex("^zbc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigAnchoredLiteralLongNonMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 15; i++ {
		x = append(x, x...)
	}
	re := onig.MustNewRegex("^zbc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigAnchoredShortMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := onig.MustNewRegex("^.bc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigAnchoredLongMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 15; i++ {
		x = append(x, x...)
	}
	re := onig.MustNewRegex("^.bc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigOnePassShortA(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := onig.MustNewRegex("^.bc(d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigNotOnePassShortA(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := onig.MustNewRegex(".bc(d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigOnePassShortB(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := onig.MustNewRegex("^.bc(?:d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigNotOnePassShortB(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := onig.MustNewRegex(".bc(?:d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigOnePassLongPrefix(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := onig.MustNewRegex("^abcdefghijklmnopqrstuvwxyz.*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigOnePassLongNotPrefix(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := onig.MustNewRegex("^.bcdefghijklmnopqrstuvwxyz.*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.MustFindMatch(string(x))
	}
}

func BenchmarkOnigMatchParallelShared(b *testing.B) {
	x := []byte("this is a long line that contains foo bar baz")
	re := onig.MustNewRegex("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.MustFindMatch(string(x))
		}
	})
}

func benchmarkOnig(b *testing.B, re string, n int) {
	r := onig.MustNewRegex(re)
	t := makeText(n)
	b.ResetTimer()
	b.SetBytes(int64(n))
	for i := 0; i < b.N; i++ {
		if r.MustFindMatch(string(t)) != nil {
			b.Fatal("match!")
		}
	}
}

func BenchmarkOnigMatchEasy0_32(b *testing.B)   { benchmarkOnig(b, easy0, 32<<0) }
func BenchmarkOnigMatchEasy0_1K(b *testing.B)   { benchmarkOnig(b, easy0, 1<<10) }
func BenchmarkOnigMatchEasy0_32K(b *testing.B)  { benchmarkOnig(b, easy0, 32<<10) }
func BenchmarkOnigMatchEasy0_1M(b *testing.B)   { benchmarkOnig(b, easy0, 1<<20) }
func BenchmarkOnigMatchEasy0_32M(b *testing.B)  { benchmarkOnig(b, easy0, 32<<20) }
func BenchmarkOnigMatchEasy0i_32(b *testing.B)  { benchmarkOnig(b, easy0i, 32<<0) }
func BenchmarkOnigMatchEasy0i_1K(b *testing.B)  { benchmarkOnig(b, easy0i, 1<<10) }
func BenchmarkOnigMatchEasy0i_32K(b *testing.B) { benchmarkOnig(b, easy0i, 32<<10) }
func BenchmarkOnigMatchEasy0i_1M(b *testing.B)  { benchmarkOnig(b, easy0i, 1<<20) }
func BenchmarkOnigMatchEasy0i_32M(b *testing.B) { benchmarkOnig(b, easy0i, 32<<20) }
func BenchmarkOnigMatchEasy1_32(b *testing.B)   { benchmarkOnig(b, easy1, 32<<0) }
func BenchmarkOnigMatchEasy1_1K(b *testing.B)   { benchmarkOnig(b, easy1, 1<<10) }
func BenchmarkOnigMatchEasy1_32K(b *testing.B)  { benchmarkOnig(b, easy1, 32<<10) }
func BenchmarkOnigMatchEasy1_1M(b *testing.B)   { benchmarkOnig(b, easy1, 1<<20) }
func BenchmarkOnigMatchEasy1_32M(b *testing.B)  { benchmarkOnig(b, easy1, 32<<20) }
func BenchmarkOnigMatchMedium_32(b *testing.B)  { benchmarkOnig(b, medium, 32<<0) }
func BenchmarkOnigMatchMedium_1K(b *testing.B)  { benchmarkOnig(b, medium, 1<<10) }
func BenchmarkOnigMatchMedium_32K(b *testing.B) { benchmarkOnig(b, medium, 32<<10) }
func BenchmarkOnigMatchMedium_1M(b *testing.B)  { benchmarkOnig(b, medium, 1<<20) }
func BenchmarkOnigMatchMedium_32M(b *testing.B) { benchmarkOnig(b, medium, 32<<20) }
func BenchmarkOnigMatchHard_32(b *testing.B)    { benchmarkOnig(b, hard, 32<<0) }
func BenchmarkOnigMatchHard_1K(b *testing.B)    { benchmarkOnig(b, hard, 1<<10) }
func BenchmarkOnigMatchHard_32K(b *testing.B)   { benchmarkOnig(b, hard, 32<<10) }
func BenchmarkOnigMatchHard_1M(b *testing.B)    { benchmarkOnig(b, hard, 1<<20) }
func BenchmarkOnigMatchHard_32M(b *testing.B)   { benchmarkOnig(b, hard, 32<<20) }
func BenchmarkOnigMatchHard1_32(b *testing.B)   { benchmarkOnig(b, hard1, 32<<0) }
func BenchmarkOnigMatchHard1_1K(b *testing.B)   { benchmarkOnig(b, hard1, 1<<10) }
func BenchmarkOnigMatchHard1_32K(b *testing.B)  { benchmarkOnig(b, hard1, 32<<10) }
func BenchmarkOnigMatchHard1_1M(b *testing.B)   { benchmarkOnig(b, hard1, 1<<20) }
func BenchmarkOnigMatchHard1_32M(b *testing.B)  { benchmarkOnig(b, hard1, 32<<20) }
