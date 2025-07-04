# This is a comment
x = 42
y = `hello`
z = [1, 2, 3]
w = {name: `test`, value: 123}

add = \a, b -> a + b
factorial = \n -> 
  when n is
    0 -> 1;
    _ -> n * factorial (n - 1)

result = add(5, 3)
maybe = Just 42 