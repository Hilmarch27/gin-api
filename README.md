Architecture:
|   .env
|   go.mod
|   go.sum
|   tree.txt
|   
+---cmd
|   \---api
|           main.go
|           
+---internal
|   +---delivery
|   |   \---http
|   |       +---handler
|   |       |       user-handler.go
|   |       |       
|   |       +---middleware
|   |       |       auth-middleware.go
|   |       |       
|   |       \---router
|   |               private-route.go
|   |               public-route.go
|   |               router.go
|   |               
|   +---domain
|   |       user.go
|   |       
|   +---repository
|   |       interfaces.go
|   |       user-repo.go
|   |       
|   \---usecase
|           auth-usecase.go
|           
\---pkg
    +---config
    |       config.go
    |       
    \---utils
