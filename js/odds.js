//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

'use strict'

//////////////////////////////////////////////////////////////////////////////
//
//  now
//
//////////////////////////////////////////////////////////////////////////////

const isNodeJS = function () { try { return Object.prototype.toString.call(global.process) === '[object process]' } catch(e) { return false } } ()

function nowNode() {
  const time = process.hrtime()
  return time[0] + (time[1] / 1e9)
}

function nowChrome() { return performance.now() / 1000 }

const now = isNodeJS ? nowNode : nowChrome

//////////////////////////////////////////////////////////////////////////////
//
//  time
//
//////////////////////////////////////////////////////////////////////////////

function time(string, f) {
  const start = now()
  let xs = f()
  const end   = now()
  const delta = ((end - start) * 1000 * 100 | 0) / 100

  console.log(`${string}: ${delta} ms ${listToString(xs)}`)
}

//////////////////////////////////////////////////////////////////////////////
//
//  List
//
//////////////////////////////////////////////////////////////////////////////

class Cons {
  constructor(x, xs) {
    this.x  = x
    this.xs = xs
  }
}

const nil = null
function cons(x, xs) { return new Cons(x, xs) }

function listToString(xs) {
  if (xs === nil) {
    return "[]"
  } else {
    return "[ " + elements(xs) + " ]"
  }

  function elements(xs) {
    if (xs.xs === nil) {
      return xs.x.toFixed(6)
    } else {
      return xs.x.toFixed(6) + ", " + elements(xs.xs)
    }
  }
}

function map(f, xs) {
  if (xs === nil) {
    return nil
  } else {
    return cons(f(xs.x), map(f, xs.xs))
  }
}

function zipWith(f, xs, ys) {
  if (xs === nil || ys == nil) {
    return nil
  } else {
    return cons(f(xs.x, ys.x), zipWith(f, xs.xs, ys.xs))
  }
}

function replicate(n, x) { 
  if (n === 0) { 
    return nil 
  } else { 
    return cons(x, replicate(n - 1, x)); 
  } 
}

//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

function oneList(w) { return cons(1, replicate(w, 0)) }

function odds(pw, w, l) {
  if (w === 0 || l === 0) {
    return oneList(w)
  } else {
    var ws = cons(0, map((x) => x * pw      , odds(pw, w - 1, l    )))
    var ls =         map((x) => x * (1 - pw), odds(pw, w    , l - 1))
    return zipWith((x, y) => x + y, ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  oddsHalf
//
//////////////////////////////////////////////////////////////////////////////

function averageList(xs, ys) {
  if (xs === nil || ys == nil) {
    return nil
  } else {
    return cons(0.5 * (xs.x + ys.x), averageList(xs.xs, ys.xs))
  }
}

function oddsHalf(w, l) {
  if (w === 0 || l === 0) {
    return oneList(w)
  } else {
    var ws = cons(0, oddsHalf(w - 1, l    ))
    var ls =         oddsHalf(w    , l - 1)
    return averageList(ws, ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Array
//
//////////////////////////////////////////////////////////////////////////////

function array(count, f) {
  const xs = new Float64Array(count)
  for (let i = 0; i < count; i++) { xs[i] = f(i) }
  return xs
}

function arrayToList(xs) {
  let ys = nil
  for (let i = xs.length - 1; i >= 0; i--) ys = cons(xs[i], ys)
  return ys
}

//////////////////////////////////////////////////////////////////////////////
//
//  oddsArray
//
//////////////////////////////////////////////////////////////////////////////

function oneArray(w) {
  const xs = new Float64Array(w + 1)
  xs[0] = 1
  return xs
}

function averageArray(w, ws, ls) {
  const averages = new Float64Array(w + 1)
  averages[0] = 0.5 * ls[0]
  for (let i = 1; i <= w; i++) { averages[i] = 0.5 * (ws[i - 1] + ls[i]) }
  return averages
}

function oddsHalfArrayInternal(w, l) {
  if (w === 0 || l === 0) {
    return oneArray(w)
  } else {
    var ws = oddsHalfArrayInternal(w - 1, l    )
    var ls = oddsHalfArrayInternal(w    , l - 1)
    return averageArray(w, ws, ls)
  }
}

function oddsHalfArray(w, l) { 
  return arrayToList(oddsHalfArrayInternal(w, l)) 
}

//////////////////////////////////////////////////////////////////////////////
//
//
//
//////////////////////////////////////////////////////////////////////////////

function oddsTest() {
  const w = 12
  const l = 12
  time("odds          ", () => odds         (0.5, w, l))
  time("oddsHalf      ", () => oddsHalf     (     w, l))
  time("oddsHalfArray ", () => oddsHalfArray(     w, l))
}

//
// Tests
//

oddsTest()
