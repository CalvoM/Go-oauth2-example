# Go oauth2 server example
This is a simple example for implementing oauth2 on Golang apis.

## Usage
We use a js client on the browser, the *UI* folder contains the client code.

To access the UI, run the server
```bash
go run server.go
```
Then, open the url **http://localhost:3001/index.html**

Generate the client credentials then login. Check the logs in the console and, the session storage.
## Credit
Thanks to the library and example from [Go oauth2 library](https://github.com/go-oauth2/oauth2).