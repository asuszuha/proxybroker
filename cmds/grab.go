package cmds

import (
  "embed"
  "fmt"
  "github.com/Ziloka/ProxyBroker/services"
  "github.com/oschwald/geoip2-golang"
  "github.com/urfave/cli/v2"
  "os"
  "strings"
)

func Grab(c *cli.Context, assetFS embed.FS) (err error) {

  // Set default values for flags
  types := c.StringSlice("types")
  if len(types) == 0 {
    types = []string{"http", "https"}
  }
  timeout := c.Int("timeout")
  if timeout == 0 {
    timeout = 5000
  }
  countries := c.StringSlice("countries")
  ports := c.StringSlice("ports")

  bytes, readFileError := assetFS.ReadFile("assets/GeoLite2-Country.mmdb")

  if readFileError != nil {
    return readFileError
  }

  db, dbErr := geoip2.FromBytes(bytes)
  if err != nil {
    return dbErr
  }
  defer db.Close()

  proxies := make(chan []string)
  go services.Collect(assetFS, db, proxies, types, countries, ports)

  displayedProxies := []string{}
  for _, proxy := range <-proxies {
    displayedProxies = append(displayedProxies, proxy)
    fmt.Println("[+] " + proxy)
  }

  if true {
    data := []byte(strings.Join(displayedProxies, "\n"))
    f, fileCreateErr := os.Create("proxies.txt")
    if fileCreateErr != nil {
      panic(fileCreateErr)
    }
    fileWriteErr := os.WriteFile("proxies.txt", data, 0644)
    if fileWriteErr != nil {
      panic(fileWriteErr)
    }
    defer f.Close()

  }
  return nil
}
