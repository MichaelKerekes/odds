//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

package main

import "fmt"
import "../odds"

//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

func main() {
  var pw = odds.Double(0.5)
  var w  = 12
  var l  = 12

  fmt.Printf("Odds                  %v\n", odds.Odds                 (pw, w, l))
  fmt.Printf("OddsHalfZip           %v\n", odds.OddsHalfZip          (    w, l))
  fmt.Printf("OddsHalf              %v\n", odds.OddsHalf             (    w, l))
  fmt.Printf("OddsHalfArray         %v\n", odds.OddsHalfArray        (    w, l))
  fmt.Printf("OddsHalfReversedArray %v\n", odds.OddsHalfReversedArray(    w, l))
  fmt.Printf("OddsHalfSlice         %v\n", odds.OddsHalfSlice        (    w, l))
}
