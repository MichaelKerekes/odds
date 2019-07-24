------------------------------------------------------------------------------
--
--  Odds
--
------------------------------------------------------------------------------

{-# language BangPatterns #-}

import Prelude
import GHC.Exts
import Control.DeepSeq        (NFData, force)
import Control.Exception      (evaluate)
import Control.Monad          (join)
import Control.Monad.IO.Class (MonadIO(..))
import Data.Int               (Int64  )
import Foreign                (alloca, peek)
import System.IO.Unsafe       (unsafePerformIO)

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

main = do
  print $ map (roundFraction 6) $ toStream $ oddsHalfStrict w l

  runTest "halfStrict" $ \p w l -> toStream $ oddsHalfStrict   w l
  runTest "strict    " $ \p w l -> toStream $ oddsStrict     p w l
  runTest "half      " $ \p w l ->            oddsHalf         w l
  runTest "odds      " odds

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

runTest :: String -> (Double -> Int -> Int -> [Double]) -> IO ()
runTest string f = do
  (dtime, ps) <- benchmark $ f half w l
  let ms = roundFraction 1 (asSeconds dtime * 1000)
  putStrLn (string ++ ": " ++ show ms ++ "ms ")

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

roundFraction :: Int -> Double -> Double
roundFraction digits x = (fromInteger (round (x * 10 ^ digits)) :: Double) / (10 ^ digits)

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

half = 0.5
w    = 12
l    = 12

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

odds :: Num a => a -> Int -> Int -> [a]
odds p 0 l = oneList 0
odds p w 0 = oneList w
odds p w l = zipWith (+) ws ls
  where
    ls =     fmap (* (1 - p)) (odds p w (l - 1))
    ws = 0 : fmap (* (p    )) (odds p (w - 1) l)

oneList w = 1 : replicate w 0

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

oddsStrict :: Double -> Int -> Int -> List Double
oddsStrict p 0 l = oneListList 0
oddsStrict p w 0 = oneListList w
oddsStrict p w l = zipWithList (+) ws ls
  where
    ls =          fmap (* (1 - p)) $ oddsStrict p w (l - 1)
    ws = Cons 0 $ fmap (* (p    )) $ oddsStrict p (w - 1) l

oneListList w = Cons 1 $ replicateList w 0

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

oddsHalf :: Int -> Int -> [Double]
oddsHalf w l = odds w l
  where
    odds :: Int -> Int -> [Double]
    odds 0 _ = oneList 0
    odds _ 0 = oneList w
    odds w l = average ls ws
      where
        ls =         oddsHalf w (l - 1)
        ws = (:) 0 $ oddsHalf (w - 1) l

average :: [Double] -> [Double] -> [Double]
average []       ys       = []
average xs       []       = []
average (x : xs) (y : ys) = 0.5 * (x + y) : average xs ys

------------------------------------------------------------------------------
--
------------------------------------------------------------------------------

oddsHalfStrict :: Int -> Int -> List Double
oddsHalfStrict 0 l = oneListList 0
oddsHalfStrict w 0 = oneListList w
oddsHalfStrict w l = averageList ls ws
  where
    ls =          oddsHalfStrict w (l - 1)
    ws = Cons 0 $ oddsHalfStrict (w - 1) l

averageList :: List Double -> List Double -> List Double
averageList Nil      ys             = Nil
averageList xs       Nil            = Nil
averageList (Cons x xs) (Cons y ys) = Cons (0.5 * (x + y)) (averageList xs ys)

------------------------------------------------------------------------------
--  List (strict)
------------------------------------------------------------------------------

data List a = Nil | Cons !a !(List a)

toStream :: List a -> [a]
toStream Nil         = []
toStream (Cons x xs) = x : toStream xs

instance Functor List where
  fmap f Nil         = Nil
  fmap f (Cons x xs) = Cons (f x) (fmap f xs)

zipWithList :: (a -> b -> c) -> List a -> List b -> List c
zipWithList f (Nil      ) (ys       ) = Nil
zipWithList f (xs       ) (Nil      ) = Nil
zipWithList f (Cons x xs) (Cons y ys) = Cons (f x y) (zipWithList f xs ys)

replicateList :: Int -> a -> List a
replicateList 0 x = Nil
replicateList n x = Cons x $ replicateList (n - 1) x

------------------------------------------------------------------------------
--  benchmark - example of using NFData
------------------------------------------------------------------------------

benchmark :: MonadIO m => (NFData a) => a -> m (Int64, a)
benchmark x = timeIO $ liftIO $ evaluate (force x)

------------------------------------------------------------------------------
--  timeIO
------------------------------------------------------------------------------

timeIO :: MonadIO m => m a -> m (Int64, a)
timeIO io = do
  !tstart <- getCPUTime
  !x      <- io
  !tend   <- getCPUTime

  return (tend - tstart, x)

------------------------------------------------------------------------------
--  CPUTime
------------------------------------------------------------------------------

asSeconds :: Int64 -> Double
asSeconds x = performancePeriod * realToFrac x

getCPUTime :: MonadIO m => m Int64
getCPUTime = liftIO queryPerformanceCounter

performancePeriod :: Double
performancePeriod = 1 / fromIntegral (unsafePerformIO queryPerformanceFrequency)

{-# NOINLINE performancePeriod #-}

------------------------------------------------------------------------------
--  QueryPerformanceCounter
------------------------------------------------------------------------------

foreign import ccall safe "QueryPerformanceCounter"   _queryPerformanceCounter   :: Ptr Int64 -> IO ()
foreign import ccall safe "QueryPerformanceFrequency" _queryPerformanceFrequency :: Ptr Int64 -> IO ()

queryPerformanceCounter   :: IO Int64
queryPerformanceFrequency :: IO Int64

queryPerformanceCounter   = alloca $ \p -> do { _queryPerformanceCounter   p; peek p }
queryPerformanceFrequency = alloca $ \p -> do { _queryPerformanceFrequency p; peek p }
