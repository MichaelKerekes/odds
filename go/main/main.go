//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

package main

import "fmt"
import "time"
import "../odds"

//////////////////////////////////////////////////////////////////////////////
//
//  timer
//
//////////////////////////////////////////////////////////////////////////////

func timer(name string, f func() fmt.Stringer) {
  start := time.Now()
  xs    := f()
  delta := time.Now().Sub(start).Seconds()
  fmt.Printf("%s: %7.2f ms %v\n", name, 1000.0 * delta, xs)
}

//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

func main() {
  var pw = odds.Double(0.5)
  var w  = 12
  var l  = 12

  timer("Odds         ", func() fmt.Stringer { return odds.Odds         (pw, w, l) } )
  timer("OddsHalf     ", func() fmt.Stringer { return odds.OddsHalf     (    w, l) } )
  timer("OddsHalfArray", func() fmt.Stringer { return odds.OddsHalfArray(    w, l) } )
  timer("OddsHalfSlice", func() fmt.Stringer { return odds.OddsHalfSlice(    w, l) } )
}
