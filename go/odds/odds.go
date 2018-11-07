//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

package odds

import "fmt"

//////////////////////////////////////////////////////////////////////////////
//
//  Double
//
//////////////////////////////////////////////////////////////////////////////

type Double float64

func (x Double) String() string { return fmt.Sprintf("%.6f", float64(x)) }

//////////////////////////////////////////////////////////////////////////////
//
//  List
//
//////////////////////////////////////////////////////////////////////////////

type List struct {
  x  Double
  xs *List
}

func cons(x Double, xs *List) *List { return &List{x, xs} }

func (xs *List) head  () Double { return xs.x  }
func (xs *List) tail  () *List  { return xs.xs }
func (xs *List) String() string { return "[" + xs.elementsString(true) + "]" }

func (xs *List) elementsString(first bool) string {
  if xs == nil {
    return ""
  } else if first {
    return xs.head().String() + xs.tail().elementsString(false)
  } else {
    return " " + xs.head().String() + xs.tail().elementsString(false)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  List operators
//
//////////////////////////////////////////////////////////////////////////////

func replicate(n int, x Double) *List {
  switch {
    case n == 0: return nil
    default    : return cons(x, replicate(n - 1, x))
  }
}

func mapList(f func(Double) Double, xs *List) *List {
  switch {
    case xs == nil: return nil
    default       : return cons(f(xs.head()), mapList(f, xs.tail()))
  }
}

func zipWith(f func(Double, Double) Double, xs *List, ys *List) *List {
  switch {
    case xs == nil: return ys
    case ys == nil: return xs
    default: return cons(f(xs.head(), ys.head()), zipWith(f, xs.tail(), ys.tail()))
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Odds
//
//////////////////////////////////////////////////////////////////////////////

func oneList(w int) *List { return cons(1, replicate(w, 0)) }

func Odds(pw Double, w int, l int) *List {
  switch {
    case w == 0: return oneList(w)
    case l == 0: return oneList(w)
    default:
      var ws = cons(0, mapList(func(x Double) Double { return x * pw       }, Odds(pw, w - 1, l    )))
      var ls =         mapList(func(x Double) Double { return x * (1 - pw) }, Odds(pw, w    , l - 1))
      return zipWith(func(x Double, y Double) Double { return x + y }, ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfZip
//
//////////////////////////////////////////////////////////////////////////////

func average(x Double, y Double) Double { return 0.5 * (x + y) }

func OddsHalfZip(w int, l int) *List {
  switch {
    case w == 0: return oneList(w)
    case l == 0: return oneList(w)
    default:
      var ws = cons(0, OddsHalfZip(w - 1, l    ))
      var ls =         OddsHalfZip(w    , l - 1)
      return zipWith(average, ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalf
//
//////////////////////////////////////////////////////////////////////////////

func averageList(xs *List, ys *List) *List {
  switch {
    case xs == nil: return nil
    case ys == nil: return nil
    default: return cons(0.5 * (xs.head() + ys.head()), averageList(xs.tail(), ys.tail()))
  }
}

func OddsHalf(w int, l int) *List {
  switch {
    case w == 0: return oneList(w)
    case l == 0: return oneList(w)
    default:
      var ws = cons(0, OddsHalf(w - 1, l    ))
      var ls =         OddsHalf(w    , l - 1)
      return averageList(ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Array
//
//////////////////////////////////////////////////////////////////////////////

func oneArray(n int) []Double {
  xs := make([]Double, n + 1)
  xs[0] = 1
  return xs
}

func averageArray(xs []Double, ys []Double) []Double {
  averages := make([]Double, len(xs))
  for i := 0; i < len(averages); i++ { averages[i] = 0.5 * (xs[i] + ys[i]) }
  return averages
}

func reverseArray(xs []Double) []Double {
  n  := len(xs)
  ys := make([]Double, n)
  for i := 0; i < n; i++ { ys[i] = xs[n - 1 - i] }
  return ys
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfArray
//
//////////////////////////////////////////////////////////////////////////////

func OddsHalfArray(w int, l int) []Double {
  switch {
    case w == 0: return oneArray(w)
    case l == 0: return oneArray(w)
    default:
      var ws = append([]Double{0}, OddsHalfArray(w - 1, l    )...)
      var ls =                     OddsHalfArray(w    , l - 1)
      return averageArray(ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfReversedArray
//
//////////////////////////////////////////////////////////////////////////////

func oneReversedArray(n int) []Double {
  xs := make([]Double, n + 1)
  xs[n] = 1
  return xs
}

func oddsHalfReversedArrayInternal(w int, l int) []Double {
  switch {
    case w == 0: return oneReversedArray(w)
    case l == 0: return oneReversedArray(w)
    default:
      var ws = append(oddsHalfReversedArrayInternal(w - 1, l    ), 0)
      var ls =        oddsHalfReversedArrayInternal(w    , l - 1)
      return averageArray(ws, ls)
  }
}

func OddsHalfReversedArray(w int, l int) []Double {
  return reverseArray(oddsHalfReversedArrayInternal(w, l))
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfSlice
//
//////////////////////////////////////////////////////////////////////////////

func oneSlice(xs []Double) {
  xs[0] = 1
  for i := 1; i < len(xs); i++ { xs[i] = 0; }
}

func averageSlice(averages []Double, xs []Double, ys []Double) {
  for i := 0; i < len(averages); i++ { averages[i] = 0.5 * (xs[i] + ys[i]) }
}

func oddsHalfSliceInternal(ws []Double, l int) {
  var w = len(ws)
  switch {
    case w == 1: oneSlice(ws)
    case l == 0: oneSlice(ws)
    default:
      ls := make([]Double, w)
      ws[0] = 0
      oddsHalfSliceInternal(ws[1:], l    )
      oddsHalfSliceInternal(ls    , l - 1)
      averageSlice(ws, ws, ls)
  }
}

func OddsHalfSlice(w int, l int) []Double {
  ws := make([]Double, w + 1)
  oddsHalfSliceInternal(ws, l)
  return ws
}
