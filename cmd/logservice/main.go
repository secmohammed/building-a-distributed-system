package main

import (
    "context"
    "fmt"
    "go-building-distributed-application/log"
    "go-building-distributed-application/registry"
    "go-building-distributed-application/service"
    stlog "log"
)

func main() {
    log.Run("./app.log")
    host, port := "localhost", "4000"
    var r registry.Registration
    serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
    r.ServiceName = registry.LogService
    r.ServiceURL = serviceAddress
    r.HeartbeatURL = r.ServiceURL + "/heartbeat"
    r.RequiredServices = make([]registry.ServiceName, 0)
    r.ServiceUpdateURL = r.ServiceURL + "/services"

    ctx, err := service.Start(context.Background(), host, port, r, log.RegisterHandlers)
    if err != nil {
        stlog.Fatal(err)
    }
    <-ctx.Done()
    fmt.Println("shutting down logservice")
}
