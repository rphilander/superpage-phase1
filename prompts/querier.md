This is a new greenfield project.
We will build this in Go.
It is called Querier.
The executable will be called “querier”.

The Querier is part of a larger system. It depends upon a component called the Parser. Use curl to GET /doc from localhost:8081. That is the documentation for the API of the Parser.

The role of the Querier is to obtain the Hacker News data from the Parser and filter and sort that data in a manner specified by the client of its API. The Querier does the filtering and sorting on a copy of underlying data which it stores locally in memory (no persistence to disk). The Querier has separate endpoints in its API for refreshing the underlying data by pulling a fresh set of data from the Parser, and for querying the underlying data.

The Querier will have two required command line arguments. --api <port-no> determines which port number the Querier listens for HTTP requests to its REST API. --parser tells the Querier which localhost port the Parser is listening on.

The REST API will have an endpoint POST /query through which clients will obtain the Hacker News content with various filters and sorts applied. The body of the request will be a JSON object specifying the logic to be applied. The client can optionally specify a filter for each field in the Hacker News data. If no filter is specified for a field then the data is not filtered upon that field. If no filters are specified at all then all of the data is returned. If more than one filter is specified then the criteria are applied like a Boolean “and” – only rows meeting all of the filter criteria are returned (if any).

The filters supported for each field are appropriate to the field itself. For strings (e.g. headlines, usernames) the client provides a string and the Querier applies a fuzzy match against that field (use the widely used Go fuzzymatch module). For integer values (e.g. number of points, number of comments, story rank) the filter is a one-sided or two-sided range. For the age of the post, the client similarly specifies a range in the same manner as that field is defined: a time unit (hours, days, etc.) and a number relative to those units; the range can be one-sided or two-sided.

If specifying a sort order, the client specifies an array of fields and a direction (ascending or descending) for each field. The array can be empty to indicate no sorting is to be performed.

The REST API also has an endpoint POST /refresh which causes the Querier to fetch a new set of data from the Parser.

If the Querier receives a request to POST /query before any request to POST /refresh, so that there is no data in memory to query, then the Querier will first pull a set of data from the Parser and the proceed with running the query.

The REST API has another endpoint GET /doc which returns detailed documentation of the entire REST API, including example requests and responses. The documentation does not need to concern itself with system internals or how to operate the system – it is only for clients of the REST API.

When you are done coding make sure the system builds and runs correctly, then write a detailed README in case another developer needs to debug or enhance this system in the future. To test the system you will use curl to obtain the transformed data from the Querier and then also use curl to pull the raw data from the Parser and make sure the two align.

