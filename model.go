package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    . "./model"
)

type Model struct {
    ID              string `json:"id,omitempty"`
    Name            string `json:"name,omitempty"`
    Make            string `json:"make,omitempty"`
    CategoryUUID    string `json:"category_uuid,omitempty"`
}

func CreateModel(w http.ResponseWriter, r *http.Request) {
    var newModel Model

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newModel); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatement := `
            INSERT INTO models (name, make, category_uuid)
            VALUES ($1, $2, $3)`

    res, err := db.Exec(sqlStatement, newModel.Name, newModel.Make, newModel.CategoryUUID);
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"Model not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newModel)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteModel(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM models WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["model_uuid"]
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
        http.Error(w, "{\"error\":\"Model not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"Model was deleted\"}", http.StatusOK)
    }
}

func GetModelByUUID(w http.ResponseWriter, r *http.Request) {
    var newModel Model

    sqlStatement := `SELECT id, name, make, category_uuid FROM models WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["model_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newModel.Name, &newModel.Make, &newModel.CategoryUUID);

    if resp != nil {
        http.Error(w, "{\"error\":\"Model not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newModel)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexModels(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, name, make, category_uuid FROM models")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    models := make([]*Model, 0)

    for rows.Next() {
        model := new(Model)
        err := rows.Scan(&model.ID, &model.Name, &model.Make, &model.CategoryUUID)
        if err != nil {
            log.Panic(err)
        }
        models = append(models, model)
    }
    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"Model not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(models)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateModel(w http.ResponseWriter, r *http.Request) {
    var newModel Model

    sqlStatement := `UPDATE models  SET name = $1, make = $2, category_uuid = $3    WHERE id = $4;`
    vars := mux.Vars(r)
    var id = vars["model_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newModel); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newModel.Name, newModel.Make, newModel.CategoryUUID, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"Model not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newModel)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
