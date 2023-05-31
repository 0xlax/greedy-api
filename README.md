# In memory key value database

This is a simple in-memory key-value datastore that allows performing operations based on certain commands. It exposes a REST API for communication, allowing clients to interact with the datastore using HTTP requests. The datastore stores data in memory and supports concurrent operations.

## Sample Inputs and Output

| Input                                                                                   | Output                                                                                       |
|-----------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------|
| {"command": "SET hello world"}                                                          | -                                                                                            |
| {"command": "123 SET COMMAND"}                                                          | {"error": "invalid command"}                                                                 |
| {"command": "GET hello"}                                                                 | {"value": "world"}                                                                           |
| {"command": "GET hello-123"}                                                             | {"error": "key not found"}                                                                   |
| {"command": "QPUSH list_a a"}                                                            | -                                                                                            |
| {"command": "QPOP list_a"}                                                               | {"value": "a"}                                                                               |
| {"command": "QPOP list_a"}                                                               | {"error": "queue is empty"}                                                                  |
| {"command": "BQPOP list_a 10"}                                                           | {"error": "queue is empty"}                                                                  |
| {"command": "BQPOP list_1 0"}                                                            | {"value": "a"}                                                                               |
| {"command": "BQPOP list_1 0"}                                                            | {"error": "null"}                                                                             |
| {"command": "BQPOP list_1 10"}                                                           | {"error": "null after 10 seconds"}                                                           |
| {"command": "BQPOP list_1 10"}                                                           | {"error": "null after 10 seconds"}                                                           |
| {"command": "QPUSH list_1 a"}                                                            | -                                                                                            |
| {"command": "QPOP list_1"}                                                               | {"value": "a"}                                                                               |
| {"command": "QPOP list_1"}                                                               | {"value": "2"}                                                                               |
| {"command": "QPOP list_x"}                                                               | {"error": "null"}                                                                             |
