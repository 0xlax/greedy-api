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







