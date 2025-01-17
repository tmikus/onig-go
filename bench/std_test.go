package bench

import (
	"regexp"
	"strings"
	"testing"
)

func BenchmarkStdLiteral(b *testing.B) {
	x := strings.Repeat("x", 50) + "y"
	b.StopTimer()
	re := regexp.MustCompile("y")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.MatchString(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkStdNotLiteral(b *testing.B) {
	x := strings.Repeat("x", 50) + "y"
	b.StopTimer()
	re := regexp.MustCompile(".y")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.MatchString(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkStdMatchClass(b *testing.B) {
	b.StopTimer()
	x := strings.Repeat("xxxx", 20) + "w"
	re := regexp.MustCompile("[abcdw]")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.MatchString(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkStdMatchClass_InRange(b *testing.B) {
	b.StopTimer()
	// 'b' is between 'a' and 'c', so the charclass
	// range checking is no help here.
	x := strings.Repeat("bbbb", 20) + "c"
	re := regexp.MustCompile("[ac]")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.MatchString(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkStdAnchoredLiteralShortNonMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := regexp.MustCompile("^zbc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdAnchoredLiteralLongNonMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 15; i++ {
		x = append(x, x...)
	}
	re := regexp.MustCompile("^zbc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdAnchoredShortMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := regexp.MustCompile("^.bc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdAnchoredLongMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 15; i++ {
		x = append(x, x...)
	}
	re := regexp.MustCompile("^.bc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdOnePassShortA(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := regexp.MustCompile("^.bc(d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdNotOnePassShortA(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := regexp.MustCompile(".bc(d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdOnePassShortB(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := regexp.MustCompile("^.bc(?:d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdNotOnePassShortB(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := regexp.MustCompile(".bc(?:d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdOnePassLongPrefix(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := regexp.MustCompile("^abcdefghijklmnopqrstuvwxyz.*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdOnePassLongNotPrefix(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := regexp.MustCompile("^.bcdefghijklmnopqrstuvwxyz.*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(x)
	}
}

func BenchmarkStdMatchParallelShared(b *testing.B) {
	x := []byte("this is a long line that contains foo bar baz")
	re := regexp.MustCompile("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.Match(x)
		}
	})
}

func benchmarkStd(b *testing.B, re string, n int) {
	r := regexp.MustCompile(re)
	t := makeText(n)
	b.ResetTimer()
	b.SetBytes(int64(n))
	for i := 0; i < b.N; i++ {
		if r.Match(t) {
			b.Fatal("match!")
		}
	}
}

func BenchmarkStdMatchEasy0_32(b *testing.B)   { benchmarkStd(b, easy0, 32<<0) }
func BenchmarkStdMatchEasy0_1K(b *testing.B)   { benchmarkStd(b, easy0, 1<<10) }
func BenchmarkStdMatchEasy0_32K(b *testing.B)  { benchmarkStd(b, easy0, 32<<10) }
func BenchmarkStdMatchEasy0_1M(b *testing.B)   { benchmarkStd(b, easy0, 1<<20) }
func BenchmarkStdMatchEasy0_32M(b *testing.B)  { benchmarkStd(b, easy0, 32<<20) }
func BenchmarkStdMatchEasy0i_32(b *testing.B)  { benchmarkStd(b, easy0i, 32<<0) }
func BenchmarkStdMatchEasy0i_1K(b *testing.B)  { benchmarkStd(b, easy0i, 1<<10) }
func BenchmarkStdMatchEasy0i_32K(b *testing.B) { benchmarkStd(b, easy0i, 32<<10) }
func BenchmarkStdMatchEasy0i_1M(b *testing.B)  { benchmarkStd(b, easy0i, 1<<20) }
func BenchmarkStdMatchEasy0i_32M(b *testing.B) { benchmarkStd(b, easy0i, 32<<20) }
func BenchmarkStdMatchEasy1_32(b *testing.B)   { benchmarkStd(b, easy1, 32<<0) }
func BenchmarkStdMatchEasy1_1K(b *testing.B)   { benchmarkStd(b, easy1, 1<<10) }
func BenchmarkStdMatchEasy1_32K(b *testing.B)  { benchmarkStd(b, easy1, 32<<10) }
func BenchmarkStdMatchEasy1_1M(b *testing.B)   { benchmarkStd(b, easy1, 1<<20) }
func BenchmarkStdMatchEasy1_32M(b *testing.B)  { benchmarkStd(b, easy1, 32<<20) }
func BenchmarkStdMatchMedium_32(b *testing.B)  { benchmarkStd(b, medium, 32<<0) }
func BenchmarkStdMatchMedium_1K(b *testing.B)  { benchmarkStd(b, medium, 1<<10) }
func BenchmarkStdMatchMedium_32K(b *testing.B) { benchmarkStd(b, medium, 32<<10) }
func BenchmarkStdMatchMedium_1M(b *testing.B)  { benchmarkStd(b, medium, 1<<20) }
func BenchmarkStdMatchMedium_32M(b *testing.B) { benchmarkStd(b, medium, 32<<20) }
func BenchmarkStdMatchHard_32(b *testing.B)    { benchmarkStd(b, hard, 32<<0) }
func BenchmarkStdMatchHard_1K(b *testing.B)    { benchmarkStd(b, hard, 1<<10) }
func BenchmarkStdMatchHard_32K(b *testing.B)   { benchmarkStd(b, hard, 32<<10) }
func BenchmarkStdMatchHard_1M(b *testing.B)    { benchmarkStd(b, hard, 1<<20) }
func BenchmarkStdMatchHard_32M(b *testing.B)   { benchmarkStd(b, hard, 32<<20) }
func BenchmarkStdMatchHard1_32(b *testing.B)   { benchmarkStd(b, hard1, 32<<0) }
func BenchmarkStdMatchHard1_1K(b *testing.B)   { benchmarkStd(b, hard1, 1<<10) }
func BenchmarkStdMatchHard1_32K(b *testing.B)  { benchmarkStd(b, hard1, 32<<10) }
func BenchmarkStdMatchHard1_1M(b *testing.B)   { benchmarkStd(b, hard1, 1<<20) }
func BenchmarkStdMatchHard1_32M(b *testing.B)  { benchmarkStd(b, hard1, 32<<20) }
