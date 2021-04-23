package registry

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
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

func (r *registry) remove(url string) error {
    for i := range r.registrations {
        if r.registrations[i].ServiceURL == url {
            r.mutex.Lock()
            r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
            r.mutex.Unlock()
            return nil
        }
    }
    return fmt.Errorf("Service at URL %v not found", url)
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
    case http.MethodDelete:
        payload, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        url := string(payload)
        log.Printf("Removing service at URL: %v", url)
        err = reg.remove(url)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusInternalServerError)

            return
        }
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
}