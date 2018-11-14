//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

#![allow(non_upper_case_globals)]
#![allow(non_snake_case)]
#![allow(dead_code)]

use std::fmt;
use std::sync::Arc;
use std::time::{Duration, Instant};

//////////////////////////////////////////////////////////////////////////////
//
//  Primitive Types
//
//////////////////////////////////////////////////////////////////////////////

type Int    = i64;
type Double = f64;
type Ref<A> = Arc<A>;

//////////////////////////////////////////////////////////////////////////////
//
//  Num
//
//////////////////////////////////////////////////////////////////////////////

trait Num : Copy {
  const zero : Self;
  const half : Self;
  const one  : Self;
  fn add(Self, Self) -> Self;
  fn sub(Self, Self) -> Self;
  fn mul(Self, Self) -> Self;
}

impl Num for Double {
  const zero : Double = 0.0;
  const half : Double = 0.5;
  const one  : Double = 1.0;
  fn add(x : Self, y : Self) -> Self { x + y }
  fn sub(x : Self, y : Self) -> Self { x - y }
  fn mul(x : Self, y : Self) -> Self { x * y }
}

//////////////////////////////////////////////////////////////////////////////
//
//  List
//
//////////////////////////////////////////////////////////////////////////////

enum List<A> {
  Nil,
  Cons(A, ListRef<A>)
}

use List::{Nil, Cons};

type ListRef<A> = Ref<List<A>>;

fn nil<A>() -> ListRef<A> { Ref::new(Nil) }
fn cons<A>(x : A, xs : &ListRef<A>) -> ListRef<A> { Ref::new(Cons(x, Ref::clone(xs))) }

impl<A> fmt::Display for List<A> where A : fmt::Display {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    write!(f, "[");
    writeElements(f, true, self);
    write!(f, "]")
  }
}

fn writeElements<A>(f: &mut fmt::Formatter, first : bool, xs : &List<A>) -> () where A : fmt::Display {
  match xs {
    Nil         => (),
    Cons(x, xs) => {
      if !first { write!(f, ", "); }
      write!(f, "{:8.6}", x);
      writeElements(f, false, xs)
    }
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  List operators
//
//////////////////////////////////////////////////////////////////////////////

fn replicate<A>(n : Int, x : A) -> ListRef<A> where A : Copy {
  if n == 0 { nil() } else { cons(x, &replicate(n - 1, x)) }
}

fn map<A, B, F>(f : F, xs : &List<A>) -> ListRef<B> where F : Fn(&A) -> B {
  match xs {
    Nil         => nil(),
    Cons(x, xs) => cons(f(x), &map(f, xs))
  }
}

fn zipWith<A, B, C, F>(f : F, xs : &List<A>, ys : &List<B>) -> ListRef<C> where F : Fn(&A, &B) -> C {
  match (xs, ys) {
    (Cons(x, xs), Cons(y, ys)) => cons(f(x, y), &zipWith(f, xs, ys)),
    _                          => nil()
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  odds
//
//////////////////////////////////////////////////////////////////////////////

fn oneList<A>(w : Int) -> ListRef<A> where A : Num { cons(A::one, &replicate(w, A::zero)) }

fn odds<A>(pw : A, w : Int, l : Int) -> ListRef<A> where A : Num {
  if w == 0 || l == 0 {
    oneList(w)
  } else {
    let ws = cons(A::zero, &map(|x| A::mul(*x, pw                ), &odds(pw, w - 1, l    )));
    let ls =                map(|x| A::mul(*x, A::sub(A::one, pw)), &odds(pw, w    , l - 1));
    zipWith(|x, y| A::add(*x, *y), &ws, &ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  oddsHalf
//
//////////////////////////////////////////////////////////////////////////////

fn average<A>(xs : &List<A>, ys : &List<A>) -> ListRef<A> where A : Num {
  match (xs, ys) {
    (Cons(x, xs), Cons(y, ys)) => cons(A::mul(A::half, A::add(*x, *y)), &average(xs, ys)),
    _                          => nil()
  }
}

fn oddsHalf<A>(w : Int, l : Int) -> ListRef<A> where A : Num {
  if w == 0 || l == 0 {
    oneList(w)
  } else {
    let ws = cons(A::zero, &oddsHalf(w - 1, l    ));
    let ls =                oddsHalf(w    , l - 1);
    average(&ws, &ls)
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Slice
//
//////////////////////////////////////////////////////////////////////////////

fn toList<A>(xs : &[A]) -> ListRef<A> where A : Copy {
  if xs.len() == 0 {
    nil()
  } else {
    cons(xs[0], &toList(&xs[1..]))
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Array
//
//////////////////////////////////////////////////////////////////////////////

type Array<A> = Box<[A]>;

fn array<A, F>(count : usize, f : F) -> Array<A> where F : Fn(usize) -> A {
  let mut xs = Vec::with_capacity(count);
  for i in 0..count { xs.push(f(i)) };
  xs.into_boxed_slice()
}

fn pushFront<A>(x : A, xs : &[A]) -> Array<A> where A : Copy {
  array(xs.len() + 1, |i| if i == 0 { x } else { xs[i - 1] } )
}

//////////////////////////////////////////////////////////////////////////////
//
//  oddsHalfArray
//
//////////////////////////////////////////////////////////////////////////////

fn oneArray<A>(w : Int) -> Array<A> where A : Num {
  let count = (w + 1) as usize;
  let mut xs = Vec::with_capacity(count);
  xs.push(A::one);
  for _ in 1..count { xs.push(A::zero) };
  xs.into_boxed_slice()
}

fn averageArray<A>(w : Int, ws : &[A], ls : &[A]) -> Array<A> where A : Num {
  let count = (w + 1) as usize;
  let mut averages = Vec::with_capacity(count);
  averages.push(A::mul(A::half, ls[0]));
  for i in 1..count { averages.push(A::mul(A::half, A::add(ws[i - 1], ls[i]))) };
  averages.into_boxed_slice()
}

fn oddsHalfArrayInternal<A>(w : Int, l : Int) -> Array<A> where A : Num {
  if w == 0 || l == 0 {
    oneArray(w)
  } else {
    let ws = oddsHalfArrayInternal(w - 1, l    );
    let ls = oddsHalfArrayInternal(w    , l - 1);
    averageArray(w, &ws, &ls)
  }
}

fn oddsHalfArray<A>(w : Int, l : Int) -> ListRef<A> where A : Num {
  toList(&oddsHalfArrayInternal(w, l))
}

//////////////////////////////////////////////////////////////////////////////
//
//  oddsHalfSlice
//
//////////////////////////////////////////////////////////////////////////////

fn oneSlice<A>(xs : &mut [A]) where A : Num {
  xs[0] = A::one;
  for i in 1..xs.len() { xs[i] = A::zero }
}

fn averageSlice<A>(xs : &mut [A], ys : &[A]) where A : Num {
  for i in 0..xs.len() { xs[i] = A::mul(A::half, A::add(xs[i], ys[i])) }
}

fn oddsHalfSliceInternal<A>(ws : &mut [A], l : Int) where A : Num {
  let w1 = ws.len();
  if w1 == 1 || l == 0 {
    oneSlice(ws)
  } else {
    let mut ls = vec![A::zero; w1];
    ws[0] = A::zero;
    oddsHalfSliceInternal(&mut ws[1..], l    );
    oddsHalfSliceInternal(&mut ls     , l - 1);
    averageSlice(ws, &ls)
  }
}

fn oddsHalfSlice<A>(w : Int, l : Int) -> ListRef<A> where A : Num {
  let mut ws = vec![A::zero; (w + 1) as usize];
  oddsHalfSliceInternal(&mut ws, l);
  toList(&ws)
}

//////////////////////////////////////////////////////////////////////////////
//
//  time
//
//////////////////////////////////////////////////////////////////////////////

fn toSeconds(d : Duration) -> f64 { d.as_secs() as f64 + 1e-9 * (d.subsec_nanos() as f64) }

fn time<A, F>(string : &str, f : F) where F : Fn() -> A, A : fmt::Display {
  let start    = Instant::now();
  let x        = f();
  let duration = toSeconds(start.elapsed());
  println!("{}: {:7.2}ms {}", string, 1000.0 * duration, x);
}

//////////////////////////////////////////////////////////////////////////////
//
//  main
//
//////////////////////////////////////////////////////////////////////////////

fn main() {
  let w = 12;
  let l = 12;

  time("odds          ", || odds                    (Double::half, w, l));
  time("oddsHalf      ", || oddsHalf      ::<Double>(              w, l));
  time("oddsHalfArray ", || oddsHalfArray ::<Double>(              w, l));
  time("oddsHalfSlice ", || oddsHalfSlice ::<Double>(              w, l));
}
