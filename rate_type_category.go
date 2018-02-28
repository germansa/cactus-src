package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)
type RateTypeCategory struct {
    ID              string  `json:"id,omitempty"`
    Name            string  `json:"name"`
    RateTypeID      string  `json:"rate_type_uuid"`
    RateTypeName    string  `json:"rate_type_name"`
    CategoryID      string  `json:"category_uuid"`
    CategoryName    string  `json:"category_name"`
    Amount          float32 `json:"amount"`
    Percent         float32 `json:"percent"`
    Mandatory       bool    `json:"mandatory"`
    Active          bool    `json:"active"`
}

func CreateRateTypeCategory(w http.ResponseWriter, r *http.Request) {
    var newRateTypeCategory RateTypeCategory

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newRateTypeCategory); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatment := `INSERT INTO rate_type_categories (name, rate_type_uuid, category_uuid, amount, percent, mandatory, active)
                    VALUES ($1, $2, $3, $4, $5, $6, $7)`

    res, err := db.Exec(sqlStatment, newRateTypeCategory.Name, newRateTypeCategory.RateTypeID,
        newRateTypeCategory.CategoryID, newRateTypeCategory.Amount, newRateTypeCategory.Percent,
        newRateTypeCategory.Mandatory, newRateTypeCategory.Active)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"RateTypeCategory not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRateTypeCategory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteRateTypeCategory(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM rate_type_categories WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["rateTypeCategory_uuid"]
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
        http.Error(w, "{\"error\":\"RateTypeCategory not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"RateTypeCategory was deleted\"}", http.StatusOK)
    }
}

func GetRateTypeCategoryByUUID(w http.ResponseWriter, r *http.Request) {
    var newRateTypeCategory RateTypeCategory

    sqlStatement := `SELECT id, name, rate_type_uuid, category_uuid, amount, percent, mandatory, active
                    FROM rate_type_categories WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["rateTypeCategory_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newRateTypeCategory.Name, &newRateTypeCategory.RateTypeID, &newRateTypeCategory.CategoryID,
        &newRateTypeCategory.Amount, &newRateTypeCategory.Percent, &newRateTypeCategory.Mandatory,
        &newRateTypeCategory.Active);

    if resp != nil {
        http.Error(w, "{\"error\":\"RateTypeCategory not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRateTypeCategory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexRateTypeCategory(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT rate_type_categories.id, " +
            "rate_type_categories.name, " +
            "rate_type_categories.rate_type_uuid, " +
            "rate_type_categories.category_uuid, " +
            "rate_type_categories.amount, " +
            "rate_type_categories.percent, " +
            "rate_type_categories.mandatory, " +
            "rate_type_categories.active, " +
            "categories.name as categoryName, " +
            "rate_types.name as rateTypeName " +
            "FROM rate_type_categories " +
            "INNER JOIN categories ON rate_type_categories.category_uuid = categories.id " +
            "INNER JOIN rate_types ON rate_type_categories.rate_type_uuid = rate_types.id;")

    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    rate_type_categories := make([]*RateTypeCategory, 0)

    for rows.Next() {
        rate_type_categorie := new(RateTypeCategory)
        resp = rows.Scan(&rate_type_categorie.ID, &rate_type_categorie.Name, &rate_type_categorie.RateTypeID, &rate_type_categorie.CategoryID,
            &rate_type_categorie.Amount, &rate_type_categorie.Percent, &rate_type_categorie.Mandatory,
            &rate_type_categorie.Active, &rate_type_categorie.CategoryName, &rate_type_categorie.RateTypeName)
        if err != nil {
            log.Panic(err)
        }
        rate_type_categories = append(rate_type_categories, rate_type_categorie)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"RateTypeCategory not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(rate_type_categories)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateRateTypeCategory(w http.ResponseWriter, r *http.Request) {
    var newRateTypeCategory RateTypeCategory

    sqlStatement := `UPDATE rate_type_categories SET name = $1, rate_type_uuid = $2, category_uuid = $3, amount = $4, percent = $5, mandatory = $6, active = $7 WHERE id = $8;`

    vars := mux.Vars(r)
    var id = vars["rateTypeCategory_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newRateTypeCategory); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newRateTypeCategory.Name, newRateTypeCategory.RateTypeID,
        newRateTypeCategory.CategoryID, newRateTypeCategory.Amount, newRateTypeCategory.Percent,
        newRateTypeCategory.Mandatory, newRateTypeCategory.Active, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"RateTypeCategory not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newRateTypeCategory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
