package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type Client struct {
    ID        string `json:"id,omitempty"`
    Name 		  string `json:"name,omitempty"`
    LastName 	string `json:"last_name,omitempty"`
    Phone  		string `json:"phone,omitempty"`
    Email 		string `json:"email,omitempty"`
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
    var newClient Client

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newClient); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatement := `
            INSERT INTO clients (name, last_name, phone, email)
            VALUES ($1, $2, $3, $4)`

    res, err := db.Exec(sqlStatement, newClient.Name, newClient.LastName, newClient.Phone, newClient.Email);
    defer db.Close()

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"Client not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newClient)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM clients WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["client_uuid"]
    db, err := InitDB()
    res, err := db.Exec(sqlStatement, id);
    defer db.Close()

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err != nil {
        log.Panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"Client not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"Client was deleted\"}", http.StatusOK)
    }
}

func GetClientByUUID(w http.ResponseWriter, r *http.Request) {
    var newClient Client

    sqlStatement := `SELECT id, name, last_name, phone, email FROM clients WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["client_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newClient.Name, &newClient.LastName, &newClient.Phone, &newClient.Email);

    if resp != nil {
        http.Error(w, "{\"error\":\"Client not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newClient)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexClients(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, name, last_name, phone, email FROM clients")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    clients := make([]*Client, 0)

    for rows.Next() {
        client := new(Client)
        resp = rows.Scan(&client.ID, &client.Name, &client.LastName, &client.Phone, &client.Email)
        if err != nil {
            log.Panic(err)
        }
        clients = append(clients, client)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"Rate not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(clients)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
    var newClient Client

    sqlStatement := `UPDATE clients SET name = $1, last_name = $2, phone = $3, email = $4   WHERE id = $5;`
    vars := mux.Vars(r)
    var id = vars["client_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newClient); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newClient.Name, newClient.LastName, newClient.Phone, newClient.Email, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"Client not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newClient)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
