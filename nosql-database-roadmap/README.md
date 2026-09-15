# NoSQL Master Curriculum
### Zero → Production-Ready → Senior/Principal-Level Architecture & Interviews

I’d structure this as a **course/book rather than a MongoDB tutorial**. The goal is to understand the underlying distributed-database ideas first, then see how MongoDB, Redis, DynamoDB, Cassandra/ScyllaDB, and Neo4j make different trade-offs.

Each major stage includes **theory + implementation + architecture + interview preparation**.
Refer plan.md for details. 

---

## The Learning Sequence

I recommend we **do not simply march through the database products one after another**.

The dependency graph should be:

```text
                    NoSQL
                      │
             ┌────────┴────────┐
             │                 │
       Distributed Systems   Data Models
             │                 │
       CAP / BASE / PACELC     │
             │                 │
       Replication             │
       Partitioning            │
       Consistency             │
             └────────┬────────┘
                      │
                Data Modeling
                      │
        ┌─────────────┼─────────────┐
        │             │             │
     Document      Key-Value    Wide-Column
        │             │             │
     MongoDB       Redis       Cassandra
        │             │             │
        └─────────────┼─────────────┘
                      │
                    Graph
                      │
                    Neo4j
                      │
              Indexing & Performance
                      │
           Transactions & Consistency
                      │
            Multi-region / DR / Security
                      │
                Observability
                      │
             Advanced Architecture
                      │
               Capstone Systems
```

### Recommended major stages

**Stage 1:** Modules 0–3 — Foundations & distributed systems  
**Stage 2:** Module 4 — NoSQL data modeling  
**Stage 3:** Module 5 — MongoDB  
**Stage 4:** Module 6 — Redis + DynamoDB  
**Stage 5:** Module 7 — Cassandra + ScyllaDB  
**Stage 6:** Module 8 — Neo4j  
**Stage 7:** Module 9 — Indexing & query performance  
**Stage 8:** Modules 10–13 — consistency, transactions, reliability, DR  
**Stage 9:** Modules 14–17 — security, observability, performance, cloud  
**Stage 10:** Modules 18–19 — polyglot persistence & advanced architecture  
**Stage 11:** Module 20 — SQL vs. NoSQL decisions  
**Stage 12:** Modules 21–22 — senior/principal interviews + architecture

The key principle I'll follow if we work through this as lessons is: **we won't treat "knowing MongoDB commands" as knowing NoSQL.** You should be able to look at a workload, derive its access patterns, choose a data model, choose partition/replication/consistency strategies, predict failure modes, and defend the trade-offs in a senior-level interview.