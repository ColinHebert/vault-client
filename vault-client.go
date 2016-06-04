package main
import (
    "time"
    "log"
    "net/http"
    "os"
    "flag"
    "fmt"
    "text/tabwriter"

    vault "github.com/hashicorp/vault/api"
)

func main() {
    vaultServer := os.Getenv("VAULT_ADDR")
    vaultToken  := os.Getenv("VAULT_TOKEN")
    entry       := flag.String("read", "", "Entry to read")
    field       := flag.String("field", "", "Field to return")

    flag.Parse()

    // Http Request
    client := &http.Client{
        Timeout: time.Second * 10,
    }
    req, _ := http.NewRequest("GET", vaultServer + "/v1/" + *entry, nil)
    req.Header.Set("X-Vault-Token", vaultToken)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    secret, err := vault.ParseSecret(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Printing
    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 5, 0, 2, ' ', 0)
    fmt.Fprintf(w, "Key\tValue\n")
    fmt.Fprintf(w, "lease_id\t%s\n", secret.LeaseID)
    fmt.Fprintf(w, "lease_duration\t%d\n", secret.LeaseDuration)
    for key, value := range secret.Data {
        fmt.Fprintf(w, "%s\t%v\n", key, value)
    }
    fmt.Fprintln(w)
	  w.Flush()

    if *field != "" {
        log.Println(secret.Data[*field])
    }
}
