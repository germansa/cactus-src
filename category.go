package src

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    . "./model"
)

type Category struct {
    ID      string `json:"id,omitempty"`
    Name    string `json:"name,ommitempty"`
    Image   string `json:"image,ommitempty"`
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
    var newCategory Category

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newCategory); err != nil {
        return
    }

    defer r.Body.Close()

    sqlStatement := `
            INSERT INTO categories (name, image)
            VALUES ($1, $2)`

    res, err := db.Exec(sqlStatement, newCategory.Name, newCategory.Image);
    defer db.Close()

    if err != nil {
        panic(err)
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"Category not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newCategory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM categories WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["category_uuid"]
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
        http.Error(w, "{\"error\":\"Category not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"Category was deleted\"}", http.StatusOK)
    }
}

func GetCategoryByUUID(w http.ResponseWriter, r *http.Request) {
    var newCategory Category

    sqlStatement := `SELECT id, name, image FROM categories WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["category_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newCategory.Name, &newCategory.Image);

    if resp != nil {
        http.Error(w, "{\"error\":\"Category not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newCategory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, name, image FROM categories")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    categories := make([]*Category, 0)

    for rows.Next() {
        category := new(Category)
        resp = rows.Scan(&category.ID, &category.Name, &category.Image)
        if err != nil {
            log.Panic(err)
        }
        categories = append(categories, category)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"Category not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(categories)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
    var newCategory Category

    sqlStatement := `UPDATE categories  SET name = $1, image = $2   WHERE id = $3;`
    vars := mux.Vars(r)
    var id = vars["category_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newCategory); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement,newCategory.Name, newCategory.Image, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }
    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"Category not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newCategory)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
