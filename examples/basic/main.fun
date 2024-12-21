import lib from `./lib.fun`

# computes the sum of integers from 1..n
sum_range : Lam<Int, Int>
sum_range = fix(\rec -> \n ->
    when n == 0 is
        True t -> 0;
        False f -> n + rec(lib.dec(n))
)

sum_range(100) # result: 5050