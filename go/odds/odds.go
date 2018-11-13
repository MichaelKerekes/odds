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
  if (xs == nil || ys == nil) {
    return nil
  } else {
    return cons(f(xs.head(), ys.head()), zipWith(f, xs.tail(), ys.tail()))
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Odds
//
//////////////////////////////////////////////////////////////////////////////

func oneList(w int) *List { return cons(1, replicate(w, 0)) }

func Odds(pw Double, w int, l int) *List {
  if (w == 0 || l == 0) {
    return oneList(w)
  } else {
    var ws = cons(0, mapList(func(x Double) Double { return x * pw       }, Odds(pw, w - 1, l    )))
    var ls =         mapList(func(x Double) Double { return x * (1 - pw) }, Odds(pw, w    , l - 1))
    return zipWith(func(x Double, y Double) Double { return x + y }, ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalf
//
//////////////////////////////////////////////////////////////////////////////

func averageList(xs *List, ys *List) *List {
  if (xs == nil || ys == nil) {
    return nil
  } else {
    return cons(0.5 * (xs.head() + ys.head()), averageList(xs.tail(), ys.tail()))
  }
}

func OddsHalf(w int, l int) *List {
  if (w == 0 || l == 0) {
    return oneList(w)
  } else {
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

func oneArray(w int) []Double {
  xs := make([]Double, w + 1)
  xs[0] = 1
  return xs
}

func averageArray(w int, ws []Double, ls []Double) []Double {
  averages := make([]Double, w + 1)
  averages[0] = 0.5 * ls[0]
  for i := 1; i <= w; i++ { averages[i] = 0.5 * (ws[i - 1] + ls[i]) }
  return averages
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfArray
//
//////////////////////////////////////////////////////////////////////////////

func OddsHalfArray(w int, l int) []Double {
  if w == 0 || l == 0 {
    return oneArray(w)
  } else {
      var ws = OddsHalfArray(w - 1, l    )
      var ls = OddsHalfArray(w    , l - 1)
      return averageArray(w, ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfSlice
//
//////////////////////////////////////////////////////////////////////////////

func oneSlice(ws []Double) {
  ws[0] = 1
  for i := 1; i < len(ws); i++ { ws[i] = 0; }
}

func averageSlice(averages []Double, xs []Double, ys []Double) {
  for i := 0; i < len(averages); i++ { averages[i] = 0.5 * (xs[i] + ys[i]) }
}

func oddsHalfSliceInternal(ws []Double, l int) {
  var w = len(ws)
  if w == 1 || l == 0 {
    oneSlice(ws)
  } else {
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
