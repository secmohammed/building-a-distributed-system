package registry

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
)

//ServerPort is use to define server port
const ServerPort = ":3000"

//ServicesURL is used to define the endpoint of services.
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
    registrations []Registration
    mutex         *sync.Mutex
}

var reg = registry{registrations: make([]Registration, 0), mutex: new(sync.Mutex)}

func (r *registry) add(reg Registration) error {
    r.mutex.Lock()
    r.registrations = append(r.registrations, reg)
    r.mutex.Unlock()
    return nil
}

//Service is used as an empty struct to attach serveHttp to the registered service only.
type Service struct{}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Println("Request recieved")
    switch r.Method {
    case http.MethodPost:
        dec := json.NewDecoder(r.Body)
        var r Registration
        err := dec.Decode(&r)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        log.Printf("Adding service : %v with url %v\n", r.ServiceName, r.ServiceURL)
        err = reg.add(r)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
}
