```
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
```

```
run docker compose build --no-cache
```

```
to starting the contaner
run docker compose up -d 
```

```
to delete container
run docker compose down
run docker system prune 
```