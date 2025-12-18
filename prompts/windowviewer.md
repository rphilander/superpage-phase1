This is a new greenfield project.
We will build this in Go.
It is called the Window Viewer.
The executable will be called “windowviewer”.

Window Viewer is part of a larger system. The role of Window Viewer is to obtain snapshots from SnapshotDB which fall within a given time span, perform some computations on those snapshots, and return the results of those computations to its client.

Use curl to GET /doc from localhost:8082. That is the documentation for the API of SnapshotDB. You can also use curl to pull some data from SnapshotDB to see what it is like.

Window Viewer will have 2 required command line arguments. --api <port-no> determines which port number Window Viewer listens for HTTP requests to its REST API. --snapshotdb <port-no> tells Window Viewer which port on localhost SnapshotDB is listening for API requests.

Window Viewer allows API clients to determine which stories are the “top” stories for a given time window based. The client can choose as the criteria for “top” one of: highest ranking story achieved during window, highest number of points story achieved during window, largest number of total comments during window, largest number of incremental comments (so not counting comments made before the start of the window), largest number of incremental points. The client can also choose how many stories will be in the result set. Regardless of the criteria used to select the stories, there will be a complete and consistent set of data for each story in the API response.

The REST API has another endpoint GET /doc which returns detailed documentation of the entire REST API, including example requests and responses. The documentation does not need to concern itself with system internals or how to operate the system – it is only for clients of the REST API.

When you are done coding make sure the system builds and runs correctly, then write a detailed README in case another developer needs to debug or enhance this system in the future.

