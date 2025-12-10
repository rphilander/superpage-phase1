This is a new greenfield project.
We will build this in Go.
It is called WebUI.
The executable will be called “webui”.

The WebUI is part of a larger system. It depends upon a component called the Querier. Use curl to GET /doc from localhost:8082. That is the documentation for the API of the Querier.

The purpose of the Web UI is to allow me (the user) to run queries against the Querier and see the results nicely formatted as a web page. The Web UI will also allow me to refresh the Querier’s data.

WebUI will have two required command line arguments. --browser <port-no> specifies the port number where WebUI listens for browser HTTP requests. --querier <port-no> specifies the localhost port where Querier is listening for REST API requests.

When a browser requests GET / from the browser port, Web UI will serve a web page that displays the data from the Querier with no filtering or sorting. The web page will allow the user to filter and sort the data as is supported by the Querier’s API. The web page will also allow the user to refresh the Querier’s data by clicking a button; the web page will then show the refreshed data with whatever filters and sorts the user has applied at that moment.

The web page will have simple styling intended to be highly readable.

The web page will not pull the Querier’s data directly from the Querier, to simplify dealing with browser's security measures. Rather, it will make requests to WebUI which will in turn make requests to the Querier.

When you are done coding make sure the system builds and runs correctly, then write a detailed README in case another developer needs to debug or enhance this system in the future.

