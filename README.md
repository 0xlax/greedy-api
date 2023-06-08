# In memory key value database

This is a simple in-memory key-value datastore that allows performing operations based on certain commands. It exposes a REST API for communication, allowing clients to interact with the datastore using HTTP requests. The datastore stores data in memory and supports concurrent operations.



## Functionality

The Key-Value Store provides the following operations:

    SET: Set a key-value pair in the store.
    GET: Retrieve the value associated with a specific key.
    QPUSH: Push one or more values to a queue.
    QPOP: Pop a value from a queue.
    BQPOP: Block and pop a value from a queue, with an optional timeout.




## Handler Functions

The key-value store application includes the following handler functions:

    handleSET: Handles the SET command by setting a key-value pair in the store. Supports expiration time and optional conditions.
    handleGET: Handles the GET command by retrieving the value associated with a key from the store.
    handleQPUSH: Handles the QPUSH command by pushing one or more values to a queue.
    handleQPOP: Handles the QPOP command by popping a value from a queue.
    handleBQPOP: Handles the BQPOP command by blocking and popping a value from a queue, with an optional timeout.





Input
{
"command": "SET hello world"
}
Output
—--------------------
Input
{
"command": "123 SET COMMAND"
}
Output
{
"error": "invalid command"
}
—------------------
Input
{
"command": "GET hello"
}
Output
{
"value": "world"
}
—--------------------
Input
{
"command": "GET hello-123"
}
Output
{
"error": "key not found"
}
—---------------------
Input

{
"command": "QPUSH list_a a"
}
Output: BLANK
—------------------
Input
{
"command": "QPOP list_a a"
}
Output
{
"value": "a"
}
—-----------------------

Input
{
"command": "QPOP list_a"
}
Output
{
"error": "queue is empty"
}



1) Concatenation alternative
- Slice- Using a slice to store the values associated with a key allows for efficient append and removal operations, eliminating the need   for concatenation.
- Queue with Slice - A queue follows the FIFO (First-In-First-Out) principle, where the first value pushed is the first to be popped



2) mutex (read/write lock) - RWMutes for read/write locing (Lock(), UnRLock(), RUnlock(), lock())
you will allow multiple goroutines to concurrently read from the data store while ensuring exclusive write access when needed


3) Channels (queue validation)


4) Redis
5) tests
6) Optionals matter (BQPOP)
7) fix Expiry time
8) Must support Spaces in Input command
