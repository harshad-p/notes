We'll go through it systematically:

1. What an index actually is and how it works
2. Clustered vs nonclustered indexes
3. How SQL Server finds rows using an index
4. Single-column indexes
5. Composite indexes and **column order**
6. Included columns
7. Primary keys and unique indexes
8. Filtered indexes
9. How to decide **whether a column should be indexed**
10. How `WHERE`, `JOIN`, `ORDER BY`, and `GROUP BY` affect index choices
11. Selectivity/cardinality
12. How indexes affect `INSERT`/`UPDATE`/`DELETE`
13. How to create, modify, and remove indexes
14. How to inspect existing indexes
15. How to use the **execution plan** to determine whether an index is actually helping
16. Practical examples where we design an index from a query

And importantly, we'll work through **why** a particular index is appropriate rather than memorizing rules like "index columns in WHERE clauses."