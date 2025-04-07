# danger log

## Cache

### Orders

- for one symbol, we use **two heap** to maintain some most possible matched order in **memory** instead of quering database everytime
- for one symbol, we use a **LRU pool** to main some hot symbol in **memory**(cache)

#### danger


1. > **danger**: If the system receives a high volume of unmatched orders, they will accumulate in memory, potentially leading to an Out-Of-Memory (OOM) error.

    > **solution**: Implement a mechanism to periodically flush unmatched orders when there are too many orders

2. > **danger**: if there are too many symbol in our pool, we will meet OOM

    > **solution**: we use a **LRU pool** to keep a reasonable size of trading room for symbol

3. > **danger**: About Multithread safety for in server memory, and database, if multithread access same recourse, they will have conflicts

    > **solution**: we add mutex for this situations, so only one thread can access lru pool, or heap of a symbol


## Database

1. > **danger**: when we try to update some field like the remaining amount of an order, multithreads may have conflicts

    > **solution**: we use transcation of database and add "for update" after our sql query. so we can make sure only one thread can access this data. if failed, we rollback it and retry


1. > **danger**: account and other field should exist

    > **solution**: we set foreign key for this, so we can make sure they have related data