package service

import (
    "context"
    "fmt"
    "go-building-distributed-application/registry"
    "log"
    "net/http"
)

//Start web services.
func Start(ctx context.Context, host, port string, reg registry.Registration, registerHandlersFunc func()) (context.Context, error) {
    registerHandlersFunc()
    ctx = startService(ctx, host, port, reg.ServiceName)
    err := registry.RegisterService(reg)
    if err != nil {
        return ctx, err
    }
    return ctx, nil
}
func startService(ctx context.Context, host, port string, serviceName registry.ServiceName) context.Context {
    ctx, cancel := context.WithCancel(ctx)
    var srv http.Server
    srv.Addr = ":" + port
    go func() {
        log.Println(srv.ListenAndServe())
        cancel()
    }()
    go func() {
        fmt.Printf("%v started. press any key to stop. \n", serviceName)
        var s string
        fmt.Scanln(&s)
        err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", host, port))
        if err != nil {
            log.Println(err)
        }
        srv.Shutdown(ctx)
        cancel()
    }()
    return ctx
}
