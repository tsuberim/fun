sum_range : Lam<Int, Int>
sum_range = fix(\rec -> \n ->
    when n == 0 is
        True h -> 0;
    else n + rec(n-1)
)
sum_range(100) # result: 5050