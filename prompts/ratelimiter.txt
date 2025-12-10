This is a new greenfield project.
We will build this in Go.
It is called Rate Limiter.
The executable will be called “ratelimiter”.

The purpose of the Rate Limiter is to fetch HTML docs from URLs, but to rate limit itself so as to not be a burden upon the remote website.

The Rate Limiter will have two required command line arguments. --rate <num-sec> specifies a number of seconds as a positive integer. The semantics are that the Rate Limiter will make at most one HTTP request every <num-sec> seconds. --api <port-no> specifies the port number where the Rate Limiter’s REST API can be accessed.

The REST API will have an endpoint POST /fetch which causes the Rate Limiter to fetch a URL. The URL to be fetched is specified in the request body. The response body contains the HTML document retrieved from that URL. The request and response bodies are both JSON objects. If a request arrives too soon vis-a-vis the rate limit, then the request will block until the Rate Limiter is able to provide a response.

The REST API has another endpoint GET /doc which returns detailed documentation of the entire REST API, including example requests and responses. The documentation does not need to concern itself with system internals or how to operate the system – it is only for clients of the REST API.

When you are done coding make sure the system builds and runs correctly, then write a detailed README in case another developer needs to debug or enhance this system in the future.

