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
RWMutex allows multiple readers to access the data simultaneously, improving performance in scenarios where there are frequent read operations.
Use Channels for Synchronization: Instead of using mutexes for synchronization, you can leverage channels to coordinate access to shared resources. Channels provide a more expressive and safer way to synchronize goroutines. For example, you can use a channel to signal when a particular operation is completed or to coordinate access to critical sections of code.

Use Channels for Queue Operations: Since your code already includes queue-related operations (QPUSH, QPOP, BQPOP), you can enhance them using channels. For instance, you can use a buffered channel as a queue data structure. When a value is pushed into the queue, it can be sent to the channel, and when a value is popped from the queue, it can be received from the channel. This way, you can utilize the built-in synchronization provided by channels.



4) Redis
Done


5) tests


7) fix Expiry time


8) Must support Spaces in Input command












func handleBQPOP(w http.ResponseWriter, parts []string) {
	if len(parts) != 3 {
		sendErrorResponse(w, "invalid command format")
		return
	}

	key := parts[1]
	timeout, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		sendErrorResponse(w, "invalid timeout")
		return
	}

	store.mutex.RLock()
	kv, ok := store.Data[key]
	store.mutex.RUnlock()

	if ok {
		if timeout == 0 {
			// A value of 0 immediately returns a value from the queue without blocking. same as QPOP
			values := kv.Value
			if len(values) > 0 {
				value := values[len(values)-1]
				values = values[:len(values)-1]
				store.Data[key].Value = values
				sendValueResponse(w, value)
				return
			}
		} else if timeout > 0 {
			// convert the timeout value from seconds (represented as a float64) to a time.Duration value.
			ticker := time.NewTicker(time.Duration(timeout) * time.Second)
			select {

			// If the ticker emitted a value, it means the specified timeout duration has elapsed.
			// Send a timeout error response and return.
			case <-ticker.C:
				sendErrorResponse(w, "timeout")
				return
			// If the ticker didn't emit a value before the timeout
			case <-time.After(1 * time.Second):
				store.mutex.RLock()
				kv, ok = store.Data[key]
				store.mutex.RUnlock()
				if ok {
					values := kv.Value
					if len(values) > 0 {
						value := values[len(values)-1]
						values = values[:len(values)-1]
						store.Data[key].Value = values
						sendValueResponse(w, value)
						return
					}
				}
			}
		}
	}

	sendErrorResponse(w, "queue is empty")
}
