//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

package odds

import "fmt"
import "os"
import "testing"

//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

func TestMain(m *testing.M) {
  // call flag.Parse() here if TestMain uses flags
  fmt.Printf("hello\n");
  os.Exit(m.Run())
}

//////////////////////////////////////////////////////////////////////////////
//
//  perf
//
//////////////////////////////////////////////////////////////////////////////

func BenchmarkOddsHalfSlice(b *testing.B) {
  for i := 0; i < b.N; i++ {
    OddsHalfSlice(12, 12)
  }
}

func BenchmarkOddsHalfReversedArray(b *testing.B) {
  for i := 0; i < b.N; i++ {
    OddsHalfReversedArray(12, 12)
  }
}

func BenchmarkOddsHalfArray(b *testing.B) {
  for i := 0; i < b.N; i++ {
    OddsHalfArray(12, 12)
  }
}

func BenchmarkOddsHalf(b *testing.B) {
  for i := 0; i < b.N; i++ {
    OddsHalf(12, 12)
  }
}

func BenchmarkOddsHalfZip(b *testing.B) {
  for i := 0; i < b.N; i++ {
    OddsHalfZip(12, 12)
  }
}

func BenchmarkOdds(b *testing.B) {
  for i := 0; i < b.N; i++ {
    Odds(0.5, 12, 12)
  }
}
