package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"smeacs/util"
	"strconv"
	"strings"
)

func calcHash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func main() {
	addr := ":" + os.Getenv("PORT")

	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.Write([]byte("ok"))
			return
		} else {
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			SWLng, err := strconv.ParseFloat(r.FormValue("swlng"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			NELat, err := strconv.ParseFloat(r.FormValue("nelat"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			NELng, err := strconv.ParseFloat(r.FormValue("nelng"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			SWLat, err := strconv.ParseFloat(r.FormValue("swlat"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			bb := util.BoundingBox{
				SWLng: SWLng,
				NELat: NELat,
				NELng: NELng,
				SWLat: SWLat,
			}

			mapArgs := fmt.Sprintf(`-projwin %s -projwin_srs EPSG:4283 -of PDF sixmaps_etopo_lores.xml`, bb.ToArgs())

			mapName := fmt.Sprintf("maps/%s.pdf", calcHash(mapArgs))

			cmdArgs := fmt.Sprintf(`%s %s`, mapArgs, mapName)
			args := strings.Split(cmdArgs, " ")
			cmd := exec.Command("gdal_translate", args...)
			//.gdal_translate", args...)

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			log.Printf("INFO cmdArgs: %+v \n", cmdArgs)

			if err := cmd.Run(); err != nil {
				log.Printf("ERROR out.String(): %+v \n", out.String())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error generating map"))
				w.Write([]byte(err.Error()))
				return
			}

			http.ServeFile(w, r, mapName)
		}
	})

	s := &http.Server{
		Handler: router,
		Addr:    addr,
	}
	log.Println("INFO Starting plain http server on", os.Getenv("PORT"))

	log.Fatalf("ERROR %+v", s.ListenAndServe())
}
