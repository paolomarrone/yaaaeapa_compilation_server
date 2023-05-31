/*
 * Dynplug
 *
 * Copyright (C) 2022 Orastron Srl unipersonale
 *
 * Copyright is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3 of the License.
 *
 * Copyright is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Copyright.  If not, see <http://www.gnu.org/licenses/>.
 *
 * File authors: Paolo Marrone
 */

package main

import (
    "log"
    "net/http"
    "os"
    "path"
    "os/exec"
    "io/ioutil"
    "encoding/json"
)

const address = ":10002"

func main() {
 
    http.HandleFunc("/", handler)
    
    log.Println("Starting yaaaeapa compilation server")
    log.Println(http.ListenAndServeTLS(address, "./keys/localhost.crt", "./keys/localhost.key", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/uploadfiles":
        handleFilesInForm(w, r)
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Not Found"))
    }
}

func handleFilesInForm(w http.ResponseWriter, r *http.Request) {

    tmpdirpath, err := os.MkdirTemp("", "yaaaeapa_")
    if err != nil {
        return
    }
    defer os.RemoveAll(tmpdirpath)

    var files []map[string]interface{}

    respBody, _ := ioutil.ReadAll(r.Body)

    err = json.Unmarshal(respBody, &files)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for i := 0; i < len(files); i++ {
        name, okn := files[i]["name"].(string)
        str, oks := files[i]["str"].(string)
        if okn && oks {
            err := os.WriteFile(path.Join(tmpdirpath, name), []byte(str), 0644)
            if (err != nil) {
                log.Println("Error writing file " + name);
                return
            }
            log.Println("Wrote: " + path.Join(tmpdirpath, name))
        }
    }

    log.Println("Going to compile")
    cmd := exec.Command("./compile.sh", tmpdirpath)
    out, err := cmd.Output()
    if err != nil {
        log.Println("Could not run command: ", err)
        log.Println("Output: ", string(out))
        w.Header().Set("Compilation-log", string(err.Error()) + "\n" + string(out))
        unsuccess(w)
        return
    }
    log.Println("Compiled. Output: ", string(out))
    w.Header().Set("Compilation-log", string(out))

    fileBytes, err := ioutil.ReadFile(path.Join(tmpdirpath, "built.so"))
    if err != nil {
        log.Println("Could not read the output", err)
        w.Header().Set("Compilation-result", "failed")
        unsuccess(w)
        return
    }
    w.Header().Set("Compilation-result", "ok")
    //w.Header().Set("Content-Type", "application/octet-stream")
    w.WriteHeader(http.StatusOK)
    w.Write(fileBytes)
}

func unsuccess(w http.ResponseWriter) {
    w.WriteHeader(500)
    w.Write([]byte("Something went wrong\n"))
}
