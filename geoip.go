package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {

	http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			w.Write([]byte("Failed to dump request!"))
			w.Write([]byte("\n"))
			return
		}
		w.Write(b)
		w.Write([]byte("\n"))
	})

	http.HandleFunc("/", handler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	latLon := r.Header.Get("X-AppEngine-CityLatLong")
	parts := strings.Split(latLon, ",")
	lat := strings.TrimSpace(parts[0])
	lon := strings.TrimSpace(parts[1])

	m := map[string]string{
		"latLong": latLon,
		"city":    r.Header.Get("X-AppEngine-City"),
		"region":  r.Header.Get("X-AppEngine-Region"),
		"country": r.Header.Get("X-AppEngine-Country"),
		"lat":     lat,
		"lon":     lon,
	}

	js, err := json.Marshal(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(js))

}
