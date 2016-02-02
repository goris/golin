# golin

Golin is an approach to create a microservice that allows multiple systems or apps to use authentication without needing to implement it by their own.

The project uses JWT as a way of securing the comunication bewteen the client and the login server. Please refer to the RFC7519 for details (it's included in the links below).

Interfaces and some abstractions were made to make it as generic as possible, this could be changed to use some configuration files or changing the login methods.

This microservice is designed to be scalable since for example, it can be set aside of any API and just coordinate the validation of the tokens and can grow accordingly to the needs of the stack. 

GinGonic was used to build the API since it's quick to develop and easy to maintain and debug. Also BoltDB is used to storage Tokens and is designed to have a Blacklist for tokens that are no longer valid because of logout. A library for JWT is also include to create verification and sign methods.

Here are some useful links we read for understanding JWT and some logins procedures:
  - http://blog.brainattica.com/restful-json-api-jwt-go/
  - https://github.com/gin-gonic/contrib/blob/master/jwt/example/example.go
  - http://golangtutorials.blogspot.mx/2011/06/interfaces-in-go.html
  - https://github.com/boltdb/bolt#resources
  - http://jwt.io/
  - https://tools.ietf.org/html/rfc7519

## TODO

- Read configuration using a file to determine where the DB is
- Configuration file for the data schema and how to look for it
- Create blacklist for logout
- Connect it using some UI to verify functionality
- Deployment method (consider using Docker and tool )
- Statistics about usage
- Check reliabilty and synchronization of BoltDB
  
