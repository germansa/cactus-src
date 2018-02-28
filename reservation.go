package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type Reservation struct {
    ID               string  `json:"id"`
    ConfirmationCode string  `json:"confirmation_code"`
    ClientID         string  `json:"client_uuid"`
    Status           int     `json:"status"`
    PickupModelID    string  `json:"pickup_model_uuid"`
    CategoryID       string  `json:"category_uuid"`
    CreatedDate      string  `json:"created_date"`
    Total            float32 `json:"total"`
    PickupDate       string  `json:"pickup_date"`
    PickupTime       string  `json:"pickup_time"`
    ReturnDate       string  `json:"return_date"`
    ReturnTime       string  `json:"return_time"`
}

type ReservationReturnFilter struct {
  FirstDate        string  `json:"first_date"`
  SecondDate  string  `json:"second_date"`
}

type ReservationPickupFilter struct {
  FirstDate        string  `json:"first_date"`
  SecondDate  string  `json:"second_date"`
}

func filterDateReturn(w http.ResponseWriter, r *http.Request) {
  var filteredReservation ReservationReturnFilter
  var resp error

  sqlStatement := `SELECT SELECT id, confirmation_code, client_uuid, status, pickup_model_uuid, category_uuid, " +
            "created_date, total, pickup_date, pickup_time, return_date, return_time FROM reservations
            WHERE return_date >= $1 AND return_date <= $2`

  db, err := InitDB()

  defer db.Close()

  if err != nil {
      log.Panic(err)
  }

  decoder := json.NewDecoder(r.Body)

  if err := decoder.Decode(&filteredReservation); err != nil {
      return
  }

  rows, err := db.Query(sqlStatement, filteredReservation.FirstDate, filteredReservation.SecondDate)

  reservations := make([]*Reservation, 0)

  for rows.Next() {
      reservation := new(Reservation)
      resp = rows.Scan(&reservation.ID, &reservation.ConfirmationCode, &reservation.ClientID, &reservation.Status,
          &reservation.PickupModelID, &reservation.CategoryID, &reservation.CreatedDate, &reservation.Total,
          &reservation.PickupDate, &reservation.PickupTime, &reservation.ReturnDate, &reservation.ReturnTime)
      if err != nil {
          log.Panic(err)
      }
      reservations = append(reservations, reservation)
  }

  if resp != nil {
      http.Error(w, "{\"error\":\"Reservation not found\"}", http.StatusNotFound)

      return
  }

  json.NewEncoder(w).Encode(reservations)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
}

func filterDatePickup(w http.ResponseWriter, r *http.Request) {
  var filteredReservation ReservationPickupFilter
  var resp error

  sqlStatement := `SELECT id, confirmation_code, client_uuid, status, pickup_model_uuid, category_uuid, " +
            "created_date, total, pickup_date, pickup_time, return_date, return_time FROM reservations
            WHERE pickup_date >= $1 AND pickup_date <= $2`

  db, err := InitDB()

  defer db.Close()

  if err != nil {
      log.Panic(err)
  }

  decoder := json.NewDecoder(r.Body)

  if err := decoder.Decode(&filteredReservation); err != nil {
      return
  }

  rows, err := db.Query(sqlStatement, filteredReservation.FirstDate, filteredReservation.SecondDate)

  reservations := make([]*Reservation, 0)

  for rows.Next() {
      reservation := new(Reservation)
      resp = rows.Scan(&reservation.ID, &reservation.ConfirmationCode, &reservation.ClientID, &reservation.Status,
          &reservation.PickupModelID, &reservation.CategoryID, &reservation.CreatedDate, &reservation.Total,
          &reservation.PickupDate, &reservation.PickupTime, &reservation.ReturnDate, &reservation.ReturnTime)
      if err != nil {
          log.Panic(err)
      }
      reservations = append(reservations, reservation)
  }

  if resp != nil {
      http.Error(w, "{\"error\":\"Reservation not found\"}", http.StatusNotFound)

      return
  }

  json.NewEncoder(w).Encode(reservations)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
}

func CreateReservation(w http.ResponseWriter, r *http.Request) {
    var newReservation Reservation

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newReservation); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatment := `INSERT INTO reservations (confirmation_code, client_uuid, status, pickup_model_uuid,
                    category_uuid, created_date, total, pickup_date, pickup_time, return_date, return_time)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

    res, err := db.Exec(sqlStatment, newReservation.ConfirmationCode, newReservation.ClientID, newReservation.Status,
        newReservation.PickupModelID, newReservation.CategoryID, newReservation.CreatedDate, newReservation.Total,
        newReservation.PickupDate, newReservation.PickupTime, newReservation.ReturnDate, newReservation.ReturnTime)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"Reservation not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newReservation)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM reservations WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["reservation_uuid"]
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
        http.Error(w, "{\"error\":\"Reservation not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"Reservation was deleted\"}", http.StatusOK)
    }
}

func GetReservationByUUID(w http.ResponseWriter, r *http.Request) {
    var newReservation Reservation

    sqlStatement := `SELECT id, confirmation_code, client_uuid, status, pickup_model_uuid, category_uuid, " +
            "created_date, total, pickup_date, pickup_time, return_date, return_time FROM reservations WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["reservation_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newReservation.ConfirmationCode, &newReservation.ClientID, &newReservation.Status,
        &newReservation.PickupModelID, &newReservation.CategoryID, &newReservation.CreatedDate, &newReservation.Total,
        &newReservation.PickupDate, &newReservation.PickupTime, &newReservation.ReturnDate, &newReservation.ReturnTime);

    if resp != nil {
        http.Error(w, "{\"error\":\"Reservation not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newReservation)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexReservations(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, confirmation_code, client_uuid, status, pickup_model_uuid, category_uuid, " +
            "created_date, total, pickup_date, pickup_time, return_date, return_time FROM reservations")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    reservations := make([]*Reservation, 0)

    for rows.Next() {
        reservation := new(Reservation)
        resp = rows.Scan(&reservation.ID, &reservation.ConfirmationCode, &reservation.ClientID, &reservation.Status,
            &reservation.PickupModelID, &reservation.CategoryID, &reservation.CreatedDate, &reservation.Total,
            &reservation.PickupDate, &reservation.PickupTime, &reservation.ReturnDate, &reservation.ReturnTime)
        if err != nil {
            log.Panic(err)
        }
        reservations = append(reservations, reservation)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"Reservation not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(reservations)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateReservation(w http.ResponseWriter, r *http.Request) {
    var newReservation Reservation

    sqlStatement := `UPDATE reservations SET confirmation_code = $1, client_uuid = $2, status = $3,
        pickup_model_uuid = $4, category_uuid = $5, created_date = $6, total = $7, pickup_date = $8, pickup_time = $9,
        return_date = $10, return_time = $11 WHERE id = $12;`
    vars := mux.Vars(r)
    var id = vars["reservation_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newReservation); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newReservation.ConfirmationCode, newReservation.ClientID, newReservation.Status,
        newReservation.PickupModelID, newReservation.CategoryID, newReservation.CreatedDate, newReservation.Total,
        newReservation.PickupDate, newReservation.PickupTime, newReservation.ReturnDate, newReservation.ReturnTime, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"Reservation not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newReservation)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
