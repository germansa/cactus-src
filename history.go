package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type History struct {
    ID              string `json:"id,omitempty"`
    TableReference  string `json:"table_reference,omitempty"`
    Field           string `json:"field,omitempty"`
    Value           string `json:"value,omitempty"`
    CreatedDate     string `json:"created_date,omitempty"`
    UserID          string `json:"user_uuid, omitempty"`
    RowID           string `json:"row_uuid, omitempty"`
}

func CreateHistory(w http.ResponseWriter, r *http.Request) {
    var newHistory History

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newHistory); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatement := `
            INSERT INTO histories (table_reference, field, value, created_date, user_uuid, row_uuid)
            VALUES ($1, $2, $3, $4, $5, $6)`

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    res, err := db.Exec(sqlStatement, newHistory.TableReference, newHistory.Field, newHistory.Value, newHistory.CreatedDate, newHistory.UserID, newHistory.RowID);
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"History not created\"}", http.StatusNotFound)
        return
    }



    json.NewEncoder(w).Encode(newHistory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteHistory(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM histories WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["history_uuid"]
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
        http.Error(w, "{\"error\":\"History not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"History was deleted\"}", http.StatusOK)
    }
}

func GetHistoryByUUID(w http.ResponseWriter, r *http.Request) {
    var newHistory History

    sqlStatement := `SELECT id, table_reference, field, value, created_date, user_uuid, row_uuid FROM histories WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["history_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newHistory.TableReference, &newHistory.Field, &newHistory.Value, &newHistory.CreatedDate, &newHistory.UserID, &newHistory.RowID);

    if resp != nil {
        http.Error(w, "{\"error\":\"History not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newHistory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexHistory(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, table_reference, field, value, created_date, user_uuid, row_uuid FROM histories")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    histories := make([]*History, 0)

    for rows.Next() {
        history := new(History)
        err := rows.Scan(&history.ID, &history.TableReference, &history.Field, &history.Value, &history.CreatedDate, &history.UserID, &history.RowID)
        if err != nil {
            log.Panic(err)
        }
        histories = append(histories, history)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"History not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(histories)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateHistory(w http.ResponseWriter, r *http.Request) {
    var newHistory History

    sqlStatement := `UPDATE histories   SET table_reference = $1, field = $2, value = $3, created_date = $4, user_uuid = $5, row_uuid = $6  WHERE id = $7;`
    vars := mux.Vars(r)
    var id = vars["history_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newHistory); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newHistory.TableReference, newHistory.Field, newHistory.Value, newHistory.CreatedDate, newHistory.UserID, newHistory.RowID, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"Rate not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newHistory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
