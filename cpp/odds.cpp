//----------------------------------------------------------------------------
//  Odds - !!!! move
//----------------------------------------------------------------------------

namespace Odds {
  //----------------------------------------------------------------------------
  //  Array test
  //----------------------------------------------------------------------------

  template<I32 n> void oddsA   (Array<F64, n + 1> xs) { xs[0] = 0; }
  template<>      void oddsA<0>(Array<F64, 0 + 1> xs) { xs[0] = 0; }

  Array<F64, 1> foo() {
    Array<F64, 1> xs;
    xs[1] = 0;
    return xs;
  }

  //----------------------------------------------------------------------------
  //  oddsT
  //----------------------------------------------------------------------------

  template<I32 n> void zeroWinsT(F64 (&xs)[n + 1]) {
    for (auto& x : xs) { x = 0; }
    xs[0] = 1;
  }

  template<I32 n> void averageT(F64 * xs, F64* ys) {
    for (I32 index = 0; index < n; index++) {
      xs[index] = 0.5 * (xs[index] + ys[index]);
    }
  }

  template<I32 w> void oddsT(I32 l, F64 (&xs)[w + 1]) {
    if (l == 0) {
      zeroWinsT<w>(xs);
    } else {
      // wins

      F64 ys[w]; oddsT<w - 1>(l, ys);

      // losses

      oddsT<w>(l - 1, xs);

      // sum

      xs[0] = 0.5 * xs[0];

      averageT<w>(xs + 1, ys);
    }
  }

  template<> void oddsT<0>(I32, F64 (&xs)[1]) {
    xs[0] = 1;
  }

  //----------------------------------------------------------------------------
  //  odds
  //----------------------------------------------------------------------------

  void zeroWins(I32 n, F64* xs) {
    xs[0] = 1;

    for (I32 index = 1; index <= n; index++) {
      xs[index] = 0;
    }
  }

  void average(I32 n, F64* xs, F64* ys) {
    for (I32 index = 0; index < n; index++) {
      xs[index] = 0.5 * (xs[index] + ys[index]);
    }
  }

  void odds(I32 w, I32 l, F64* xs) {
    if (w == 0 || l == 0) {
      zeroWins(w, xs);
    } else {
      // losses

      odds(w, l - 1, xs);

      // wins

      F64* ys = xs + w + 1; odds(w - 1, l, ys);

      // sum

      xs[0] = 0.5 * xs[0];

      average(w, xs + 1, ys);
    }
  }

  //----------------------------------------------------------------------------
  //  run
  //----------------------------------------------------------------------------

  constexpr I32 w = 12;
  constexpr I32 l = 12;

  #define useV

  void run() {
    #ifdef useV
      F64 xs[w + 1];
    #else
      F64 xs[(w + 1) * (w + 1)];
    #endif

    #ifdef useV
      oddsT<w>(l, xs);

      for(I32 index = 0; index < w + 1; index++) {
        sout << xs[index] << "\n";
      }
    #else
      odds(w, l, xs);

      for(I32 index = 0; index < w + 1; index++) {
        sout << xs[index] << "\n";
      }
    #endif

    //

    F64 tmin = 1000;

    for (U32 i = 0; i < 32; i++) {
      SystemTime::Time timeStart = SystemTime::Get();

      #ifdef useV
        oddsT<w>(l, xs);
      #else
        odds(w, l, xs);
      #endif

      SystemTime::Time timeStop = SystemTime::Get();

      F64 delta = floor(100000 * SystemTime::Seconds(timeStop - timeStart)) / 100;

      tmin = min(tmin, delta);
    }

    sout << "\n" << tmin << "ms\n";
  }
}
