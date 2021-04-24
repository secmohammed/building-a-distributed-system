package main

import (
    "context"
    "fmt"
    "go-building-distributed-application/grades"
    "go-building-distributed-application/log"
    "go-building-distributed-application/registry"
    "go-building-distributed-application/service"
    stlog "log"
)

func main() {
    host, port := "localhost", "6000"
    serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
    var r registry.Registration
    r.ServiceName = registry.GradingService
    r.ServiceURL = serviceAddress
    r.RequiredServices = []registry.ServiceName{registry.LogService}
    r.ServiceUpdateURL = r.ServiceURL + "/services"
    ctx, err := service.Start(context.Background(),
        host,
        port,
        r,
        grades.RegisterHandlers)
    if err != nil {
        stlog.Fatal(err)
    }
    if logProvider, ok := registry.GetProvider(registry.LogService); ok == nil {
        fmt.Printf("Logging service found at: %v\n", logProvider)
        log.SetClientLogger(logProvider, r.ServiceName)
    }
    <-ctx.Done()
    fmt.Println("Shutting down grading service")
}
