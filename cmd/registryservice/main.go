package main

import (
    "context"
    "fmt"
    "go-building-distributed-application/registry"
    "log"
    "net/http"
)

func main() {
    registry.SetupRegistryService()
    http.Handle("/services", &registry.Service{})
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    var srv http.Server
    srv.Addr = registry.ServerPort
    go func() {
        log.Println(srv.ListenAndServe())
        cancel()
    }()
    go func() {
        fmt.Println("Registry service started, Press any key to stop")
        var s string
        fmt.Scanln(&s)
        srv.Shutdown(ctx)
        cancel()
    }()
    <-ctx.Done()
    fmt.Println("shutting down registry service")
}
