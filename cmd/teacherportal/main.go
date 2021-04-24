package main

import (
    "context"
    "fmt"
    "go-building-distributed-application/log"
    "go-building-distributed-application/registry"
    "go-building-distributed-application/service"
    "go-building-distributed-application/teacherportal"
    stlog "log"
)

func main() {
    err := teacherportal.ImportTemplates()
    if err != nil {
        stlog.Fatal(err)
    }
    host, port := "localhost", "5000"
    serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
    var r registry.Registration
    r.ServiceName = registry.TeacherPortal
    r.ServiceURL = serviceAddress
    r.RequiredServices = []registry.ServiceName{
        registry.LogService,
        registry.GradingService,
    }
    r.ServiceUpdateURL = r.ServiceURL + "/services"
    ctx, err := service.Start(
        context.Background(),
        host,
        port,
        r,
        teacherportal.RegisterHandlers,
    )
    if err != nil {
        stlog.Fatal(err)
    }
    if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
        log.SetClientLogger(logProvider, r.ServiceName)

    }
    <-ctx.Done()
    fmt.Println("TeacherPortal is shutting down...")
}
