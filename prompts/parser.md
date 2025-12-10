This is a new greenfield project.
We will build this in Go.
It is called Parser.
The executable will be called “parser”.

The Parser is part of a larger system. It depends upon a component called the Rate Limiter. Use curl to GET /doc from localhost:8080. That is the documentation for the API of the Rate Limiter.

The role of the Parser is to obtain the current top stories from Hacker News and parse them into structured data. It will do this in response to client requests to its REST API. The Parser will not interact with Hacker News directly, rather it will use the Rate Limiter to obtain the HTML documents for Hacker News URLs: https://news.ycombinator.com/, https://news.ycombinator.com/?p=2, https://news.ycombinator.com/?p=3, and so forth. The number of pages it pulls from Hacker News is determined by a required command line parameter. It will then parse the content from those pages and return that information to the client as a single JSON object.

The Parser will have three required command line arguments. --api <port-no> determines which port number the Parser listens for HTTP requests to its REST API. --ratelimiter tells the Parser which localhost port the Rate Limiter is listening on. --num-pages <N> tells the Parser how many Hacker News Pages to pull (must be a positive integer).

The REST API will have an endpoint POST /fetch through which clients will obtain the Hacker News content. When a request arrives from the client, the Parser will request the N URLs from the Rate Limiter sequentially – no need for concurrency. When the Rate Limiter has the HTML documents it will parse them and construct the JSON object to return to the client. The object will have some top-level metadata about the N documents (when they were fetched, and so forth), as well as an array of stories. For each story there is a JSON object with fields for the headline, the URL of the article, the username of the submitter, the number of points, the number of comments, the URL of the discussion page, the story id (HN’s identifier – so that we can track stories over time if we like), the story’s current rank, and the story's page as two fields: the units (hours, days, etc.) and the age measured in those units (this mirrors the information HN displays in its web pages).

The REST API has another endpoint GET /doc which returns detailed documentation of the entire REST API, including example requests and responses. The documentation does not need to concern itself with system internals or how to operate the system – it is only for clients of the REST API.

When you are done coding make sure the system builds and runs correctly, then write a detailed README in case another developer needs to debug or enhance this system in the future. To test the system you will use curl to obtain the structured data from the Parser and then also use curl to pull the actual HTML from Hacker News and make sure the two align.

