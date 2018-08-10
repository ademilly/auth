# auth

## setup

```bash
    $ go get github.com/ademilly/auth
```

## usage

### Login

Login function builds an anonymous function parsing request to obtain credentials from HTTP POST request JSON body, uses the given `func getRegistered() (User, error)` function to retrieve user to be compared with and return a string token if comparison checks out.

```go
    func login(w http.ResponseWriter, r *http.Request) {
        // obtain jwtKey by some means
        token, err := auth.Login(auth.Tokenizer("mydomain.com", jwtKey))(w, r)
        // handler error
        w.Write([]byte(token))
    }

    func main() {
        // setup server...
        // retrieve jwtKey from env for example
        handler.HandleFunc("/login", login)
        // serve on localhost for example
    }
```

```bash
    $ curl "Content-Type: application/json" -d '{"username":"johndoe","password":"please_subscribe"}' localhost/login
    abcd
```

### Protect

Protect function wraps any http.Handler you want to be behind login

```go
    func hello(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    }

    func main() {
        // setup server...
        // retrieve jwtKey from env for example
        tokenMiddleware := auth.TokenMiddleware(jwtKey)
        handler.HandleFunc("/hello", auth.Protect(tokenMiddleware, hello))
        // serve on localhost for example
    }
```

```bash
    $ curl localhost/hello
    Required authorization token not found
    # get some token 'abcd'
    $ curl -H "Authorization: Bearer abcd" localhost/hello
    Hello, World!
```
