Updated project will be available at [sz47/docify2](https://github.com/sz47/docify2)



# Docify

A decentralized blockchain based platform to verify and validate user documents and prevent fraud. 

This project uses IPFS with [web3.storage](https://web3.storage)'s api.
This repo contains 3 different server and frontend codes which will be used in the project in the 3 directories.

### Before you start

Encryption to the data isn't added, which means there are no passwords (currently). For using Web3.storage's IPFS API you'll need to create an account on their site. Then get the token for storage and put it on `login/cmd/whatever.go`'s line 113 or you can just create an Enviroment Variable 'Token' and give it your token.

```go
w3s_context, err := w3s.NewClient(w3s.WithToken(os.Getenv("Token"))) //this line must be changed 
```

### Getting started

There are 3 servers, run all the go files in `login/cmd`, `org/cmd` and `portal` (Builds were not tested). Then open the `localhost:port` that portal is using. 

For admin panel just check the login source code, it's really small code.
