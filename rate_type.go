package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    . "./model"
)

type RateType struct {
    ID       string  `json:"id,omitempty"`
    Name     string `json:"name"`
    Priority int    `json:"priority"`
}

func CreateRateType(w http.ResponseWriter, r *http.Request) {
    var newRateType RateType

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newRateType); err != nil {
        return
    }

    defer r.Body.Close()


    sqlStatment := `INSERT INTO rate_types (name, priority)
                    VALUES ($1, $2)`

    res, err := db.Exec(sqlStatment, newRateType.Name, newRateType.Priority)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"RateType not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRateType)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteRateType(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM rate_types WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["rateType_uuid"]
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
        http.Error(w, "{\"error\":\"RateType not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"RateType was deleted\"}", http.StatusOK)
    }
}

func GetRateTypeByUUID(w http.ResponseWriter, r *http.Request) {
    var newRateType RateType

    sqlStatement := `SELECT id, name, priority FROM rate_types WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["rateType_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newRateType.Name, &newRateType.Priority);

    if resp != nil {
        http.Error(w, "{\"error\":\"RateType not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRateType)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexRateType(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, name, priority FROM rate_types")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    rateTypes := make([]*RateType, 0)

    for rows.Next() {
        rateType := new(RateType)
        resp = rows.Scan(&rateType.ID, &rateType.Name, &rateType.Priority)
        if err != nil {
            log.Panic(err)
        }
        rateTypes = append(rateTypes, rateType)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"RateType not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(rateTypes)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateRateType(w http.ResponseWriter, r *http.Request) {
    var newRateType RateType

    sqlStatement := `UPDATE rate_types SET name = $1, priority = $2 WHERE id = $3;`
    vars := mux.Vars(r)
    var id = vars["rateType_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newRateType); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newRateType.Name, newRateType.Priority, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"RateType not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRateType)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
