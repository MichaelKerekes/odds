//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

import std.stdio;
import std.format;
import core.time;

private:

//////////////////////////////////////////////////////////////////////////////
//
//  List
//
//////////////////////////////////////////////////////////////////////////////

class List(T) {
  T    x;
  List xs;

  this(T _x, List _xs) {
    x  = _x;
    xs = _xs;
  }

  override string toString() const {
    return "[" ~ elementsString(true, this) ~ "]";
  }
}

string elementsString(T)(bool first, List!(T) xs) {
  if (xs is null) {
    return "";
  } else if (first) {
    return format!"%6f"(xs.x) ~ elementsString(false, xs.xs);
  } else {
    return " " ~ format!"%6f"(xs.x) ~ elementsString(false, xs.xs);
  }
}

List!(T) nil(T)() { return null; };
List!(T) cons(T)(T x, List!(T) xs) { return new List!(T)(x, xs); }

T        head(T)(List!(T) xs) { return xs.x;  }
List!(T) tail(T)(List!(T) xs) { return xs.xs; }
 
//////////////////////////////////////////////////////////////////////////////
//
//  List operators
//
//////////////////////////////////////////////////////////////////////////////

List!(T) replicate(T)(int n, T x) {
  if (n == 0) {
    return nil!(T)();
  } else {
    return cons!(T)(x, replicate!(T)(n - 1, x));
  }
}

List!(T) mapList(T)(T delegate(T) f, List!(T) xs) {
  if (xs is null) {
    return nil!(T)();
  } else {
    return cons(f(xs.head()), mapList(f, xs.tail()));
  }
}

List!(T) zipWith(T)(T delegate(T, T) f, List!(T) xs, List!(T) ys) {
  if (xs is null || ys is null) {
    return nil!(T)();
  } else {
    return cons(f(xs.head(), ys.head()), zipWith(f, xs.tail(), ys.tail()));
  }
}

List!(T) toList(T)(T[] xs) {
  auto ys = nil!(T)();
  for (int i = xs.length - 1; i >= 0 ; i--) { ys = cons(xs[i], ys); }
  return ys;
}

//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

List!(T) oneList(T)(int w) { return cons(1, replicate!(T)(w, 0.0)); }

List!(T) odds(T)(T pw, int w, int l) {
  if (w == 0 || l == 0) {
    return oneList!(T)(w);
  } else {
    auto ws = cons(0, mapList((T x) => x * pw      , odds(pw, w - 1, l    )));
    auto ls =         mapList((T x) => x * (1 - pw), odds(pw, w    , l - 1));
    return zipWith((T x, T y) => x + y, ws, ls);
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalf
//
//////////////////////////////////////////////////////////////////////////////

List!(T) averageList(T)(List!(T) xs, List!(T) ys) {
  if (xs is null || ys is null) {
    return nil!(T)();
  } else {
    return cons(0.5 * (xs.head() + ys.head()), averageList(xs.tail(), ys.tail()));
  }
}

List!(T) oddsHalf(T)(int w, int l) {
  if (w == 0 || l == 0) {
    return oneList!(T)(w);
  } else {
    auto ws = cons(0, oddsHalf!(T)(w - 1, l    ));
    auto ls =         oddsHalf!(T)(w    , l - 1);
    return averageList(ws, ls);
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Array
//
//////////////////////////////////////////////////////////////////////////////

T[] oneArray(T)(int w) {
  auto xs = new T[w + 1];
  xs[0] = 1;
  for (int i = 1; i <= w; i++) { xs[i] = 0; }
  return xs;
}

T[] averageArray(T)(int w, T[] ws, T[] ls) {
  auto averages = new T[w + 1];
  averages[0] = 0.5 * ls[0];
  for (int i = 1; i <= w; i++) { averages[i] = 0.5 * (ws[i - 1] + ls[i]); }
  return averages;
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfArray
//
//////////////////////////////////////////////////////////////////////////////

T[] oddsHalfArrayInternal(T)(int w, int l) {
  if (w == 0 || l == 0) {
    return oneArray!(T)(w);
  } else {
      auto ws = oddsHalfArrayInternal!(T)(w - 1, l    );
      auto ls = oddsHalfArrayInternal!(T)(w    , l - 1);
      return averageArray(w, ws, ls);
  }
}

List!(T) oddsHalfArray(T)(int w, int l) { 
  return toList(oddsHalfArrayInternal!(T)(w, l));
}

//////////////////////////////////////////////////////////////////////////////
//
//  OddsHalfSlice
//
//////////////////////////////////////////////////////////////////////////////

void oneSlice(T)(T[] ws) {
  ws[0] = 1;
  for (int i = 1; i < ws.length; i++) { ws[i] = 0; }
}

void averageSlice(T)(T[] averages, T[] xs, T[] ys) {
  for (int i; i < averages.length; i++) { averages[i] = 0.5 * (xs[i] + ys[i]); }
}

void oddsHalfSliceInternal(T)(T[] ws, int l) {
  auto w = ws.length;
  if (w == 1 || l == 0) {
    oneSlice(ws);
  } else {
      auto ls = new T[w];
      ws[0] = 0;
      oddsHalfSliceInternal!(T)(ws[1..$], l    );
      oddsHalfSliceInternal!(T)(ls    , l - 1);
      averageSlice(ws, ws, ls);
  }
}

List!(T) oddsHalfSlice(T)(int w, int l) {
  auto ws = new T[w + 1];
  oddsHalfSliceInternal!(T)(ws, l);
  return toList(ws);
}

//////////////////////////////////////////////////////////////////////////////
//
//  timer
//
//////////////////////////////////////////////////////////////////////////////

void timer(string name, List!(double) delegate() f) {
  const start = MonoTime.currTime;
  const xs    = f();
  const delta = 1.0e-9 * (MonoTime.currTime - start).total!"nsecs";
  //const delta    = to!("seconds", T, Duration)(duration);
  // !!!! toString of generic List doesn't work
  writefln("%s: %7.2f ms %s", name, 1000.0 * delta, xs.toString());
}

//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

void main() {
  auto pw = 0.5;
  auto w  = 12;
  auto l  = 12;

  timer("odds         ", () => odds                  (pw, w, l) );
  timer("oddsHalf     ", () => oddsHalf     !(double)(    w, l) );
  timer("oddsHalfArray", () => oddsHalfArray!(double)(    w, l) );
  timer("oddsHalfSlice", () => oddsHalfSlice!(double)(    w, l) );
}
