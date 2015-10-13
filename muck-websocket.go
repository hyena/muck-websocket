package main

import (
    "flag"
    "log"
    "net"
    "net/http"
    "sync"

    "github.com/Cristofori/kmud/telnet"
    "github.com/gorilla/websocket"
)

// Flags
var addr = flag.String("addr", "localhost:8000", "http service address")
var muckHost = flag.String("muck", "localhost:4021",
        "host and port for proxied muck")


var upgrader = websocket.Upgrader{} // use default options

func openTelnet() (t *telnet.Telnet, err error) {
    conn, err := net.Dial("tcp", *muckHost)
    if err != nil {
        return nil, err
    }

    return telnet.NewTelnet(conn), nil
}

func telnetProxy(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "Not found", 404)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", 405)
        return
    }
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Error creating websocket", 500)
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()

    log.Printf("Opening a proxy for '%s'", r.RemoteAddr)
    t, err := openTelnet();
    if err != nil {
        log.Println("Error opening telnet proxy: ", err)
        return
    }
    defer t.Close()

    // TODO: Handle Telnet codes and send across the ip to FB.
    log.Printf("Connection open for '%s'. Proxying.", r.RemoteAddr)

    var wg sync.WaitGroup
    var once sync.Once
    wg.Add(1)  // Exit when either goroutine stops.

    // Send messages from the websocket to the MUCK.
    go func() {
        defer once.Do(func () { wg.Done() })
        for {
            _, bytes, err := c.ReadMessage()
            if err != nil {
                log.Printf("Error reading from ws(%s): %v", r.RemoteAddr, err)
                break
            }
            // TODO: Partial writes.
            if _, err := t.Write(bytes); err != nil {
                log.Printf("Error sending message to Muck for %s: %v",
                        r.RemoteAddr, err)
                break
            }
        }
    }()

    // Send messages from the MUCK to the websocket.
    go func() {
        defer once.Do(func () { wg.Done() })
        for {
            bytes := make([]byte, 1024)
            if _, err := t.Read(bytes); err != nil {
                log.Printf("Error reading from muck for %s: %v",
                        r.Host, err)
                break
            }
            if err := c.WriteMessage(websocket.TextMessage, bytes); err != nil {
                log.Printf("Error sending to ws(%s): %v", r.RemoteAddr, err)
                break
            }
        }
    }()

    // Wait until either go routine exits and then close both connections.
    wg.Wait()
    log.Printf("Proxying completed for %s", r.RemoteAddr)
}

func main() {
    flag.Parse()
    log.SetFlags(0)

    http.HandleFunc("/", telnetProxy)
    err := http.ListenAndServe(*addr, nil)
    // Use this instead if you want to do SSL. You'll need to use `openssl`
    // to generate "cert.pem" and "key.pem" files.
    // err := http.ListenAndServeTLS(*addr, "cert.pem", "key.pem", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
