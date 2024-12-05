# ready-vue-go
simple tcp chat app built with a Vue frontend and a Go backend...
## structure:
```bash
$ ls
\__vue-frontend/ # The frontend for the application
\____App.vue # The single-page vue frontend
\____...
\__go-backend/ # The backend for the application, that works on a single main.go file
\____main.go: # The whole backend for the application
\____...
```

## how to run:
1. run the go backend first with `go run main.go`
2. then, you have to run the vue frontend with `npm run serve`
3. you can open it in multiple windows to simulate multiple users
