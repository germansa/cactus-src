package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    . "./model"
)

type Rates struct {
    ID              string  `json:"id,omitempty"`
    Name            string  `json:"name"`
    Quantity        int     `json:"quantity"`
    Unit_price      float32 `json:"unit_price"`
    RateTypeID      string  `json:"rate_type_uuid"`
    ReservationID   string  `json:"reservation_uuid"`
}

func CreateRates(w http.ResponseWriter, r *http.Request) {
    var newRate Rates

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newRate); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatment := `INSERT INTO rates (name, quantity, unit_price, rate_type_uuid, reservation_uuid)
                    VALUES ($1, $2, $3, $4, $5)`

    res, err := db.Exec(sqlStatment, newRate.Name, newRate.Quantity, newRate.Unit_price, newRate.RateTypeID, newRate.ReservationID)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"Rate not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRate)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteRates(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM rates WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["rates_uuid"]
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
        http.Error(w, "{\"error\":\"Rate not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"Rate was deleted\"}", http.StatusOK)
    }
}

func GetRatesByUUID(w http.ResponseWriter, r *http.Request) {
    var newRate Rates

    sqlStatement := `SELECT id, name, quantity, unit_price, rate_type_uuid, reservation_uuid FROM rates WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["rates_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newRate.Name, &newRate.Quantity, &newRate.Unit_price, &newRate.RateTypeID, &newRate.ReservationID);

    if resp != nil {
        http.Error(w, "{\"error\":\"Rate not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRate)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexRates(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, name, quantity, unit_price, rate_type_uuid, reservation_uuid FROM rates")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    rates := make([]*Rates, 0)

    for rows.Next() {
        rate := new(Rates)
        resp = rows.Scan(&rate.ID, &rate.Name, &rate.Quantity, &rate.Unit_price, &rate.RateTypeID, &rate.ReservationID)
        if err != nil {
            log.Panic(err)
        }
        rates = append(rates, rate)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"Rate not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(rates)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateRates(w http.ResponseWriter, r *http.Request) {
    var newRate Rates

    sqlStatement := `UPDATE rates SET name = $1, quantity = $2, unit_price = $3, rate_type_uuid = $4, reservation_uuid = $5 WHERE id = $6;`
    vars := mux.Vars(r)
    var id = vars["rates_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newRate); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newRate.Name, newRate.Quantity, newRate.Unit_price, newRate.RateTypeID, newRate.ReservationID, id)
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

    json.NewEncoder(w).Encode(newRate)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
