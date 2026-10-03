# Cheatsheet

1. Boxing happens when a value type is converted to a reference type
1. Boxing: Value type represented as an object/reference
1. Unboxing: Extracting the value from a boxed value
1. Generics often help avoid boxing.
1. `int?` is shorthand for: `Nullable<int>`
1. Nullable value types allow you to represent absence of a value without inventing a special value
1. `if (age is int actualAge)`
1. string? = a compiler annotation saying the reference is allowed to be null.
1. constructors are not inherited by derived classes.
1. A property provides controlled access to data.
1. `readonly` field can only be assigned during initialization or construction
1. `readonly` applies to **fields**, not properties
1. `const` means the value is a **compile-time constant**.
1. `init` usefulness: lets you create an immutable-style object without requiring a constructor with many parameters.
1. Object-oriented programming (OOP) is a way of organizing software around objects that combine state and behavior.
1. obj is Person other
1. ReferenceEquals(a, b)
1. IEquatable<Person> provides a strongly typed equality method. It avoids needing to cast from object and is useful to generic collections and APIs that recognize it.
1. Don't mutate fields used for equality or hashing while an object is stored in a hash-based collection. The collection may have placed the object based on its old hash code and then fail to find it.
1. For strings, == compares string contents, not object identity
1. Arrays are reference types
1. `Array.Copy()`
1. `Clone()` creates a shallow copy of the array.
1. `CopyTo()` copies all elements into an existing destination array, starting at a specified destination index
1. `matrix.Rank`
1. `numbers.IsReadOnly`
1. `Length` gives the total element count, while `GetLength(dimension)` gives the length of an individual dimension.
1. `Array.IndexOf()` and `Array.LastIndexOf()`
1. `Comparer<T>.Create()`
1. stable sort
1. `new List<int>(n)`: creates an empty list with capacity for at least `n` elements. It does not create a list containing `n` elements.
1. `Capacity` is how many elements the internal storage array can hold before it needs to grow.
1. `TrimExcess()`: Reduces capacity when there is substantial unused space
1. Even if you create a list with enough capacity, you still can't access an index beyond its current `Count`
1. `EnsureCapacity(capacity)` doesn't shrink the list if the requested capacity is smaller than its current capacity.
1. `employees[101] = "Charlie";`: If `101` doesn't exist, the indexer creates a new entry. This is different from `Add()`, which throws if the key already exists.
1. The indexer adds a missing key or updates an existing key.
1. `TryAdd()` returns `false` if the key already exists and doesn't overwrite its value.
1. Adding duplicate key via Dictionary.Add() throws exception; Adding duplicate value via HashSet.Add() just ignores it. 