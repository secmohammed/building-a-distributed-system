package registry

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

//RegisterService is used to register services.
func RegisterService(r Registration) error {
    buf := new(bytes.Buffer)
    enc := json.NewEncoder(buf)
    err := enc.Encode(r)
    if err != nil {
        return err
    }
    res, err := http.Post(ServicesURL, "application/json", buf)
    if err != nil {
        return err
    }
    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to register service, registry service responded with code %v", res.StatusCode)
    }
    return nil
}

//ShutdownService is used to shutdown service.
func ShutdownService(serviceURL string) error {
    req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(serviceURL)))
    if err != nil {
        return err
    }
    req.Header.Add("Content-Type", "text/plain")
    res, err := http.DefaultClient.Do(req)
    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to deregister service, service responded with code %v", res.StatusCode)
    }
    return err
}
