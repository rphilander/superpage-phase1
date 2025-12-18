This is a new greenfield project.
We will build this in Go.
It is called SnapshotDB.
The executable will be called “snapshotdb”.

SnapshotDB is part of a larger system. The role of SnapshotDB is to obtain snapshots of Hacker New content from another component called the Parser and store data using SQLite. It also has a REST API with which clients can obtain snapshots.

Use curl to GET /doc from localhost:8081. That is the documentation for the API of the Parser. You can also use curl to pull data from the Parser to see what it is like.

SnapshotDB will have 4 required command line arguments. --api <port-no> determines which port number SnapshotDB listens for HTTP requests to its REST API. --db <path> tells SnapshotDB where the SQLite database will be or is located. --parser <port-no> tells SnapshotDB which port on localhost the Parser is listening for API requests. --freq <num-seconds> determines how often (in seconds) SnapshotDB pulls a snapshot from the Parser, e.g. 60 means once every 60 seconds.

SnapshotDB REST API will allow clients to obtain all of the snapshots within a time window, or obtain all of the story ids (deduped) within a time window, or obtain all of the data specific to one particular story within a time window.

SnapshotDB REST API will also have an endpoint GET /status for obtaining operational data: time running, number of snapshots obtained, number of errors, etc.

When SnapshotDB is started, if there is no previous snapshot in the database then it should immediately obtain one from the Parser. Note that in this scenario the Parser might not have started yet. SnapshotDB should do intelligent retries with exponential backoff for a few seconds before considering it an error. If SnapshotDB is started and there is at least one snapshot in the database, SnapshotDB will wait until the appropriate amount of time has passed before obtaining the next snapshot.

The REST API has another endpoint GET /doc which returns detailed documentation of the entire REST API, including example requests and responses. The documentation does not need to concern itself with system internals or how to operate the system – it is only for clients of the REST API.

When you are done coding make sure the system builds and runs correctly, then write a detailed README in case another developer needs to debug or enhance this system in the future.


