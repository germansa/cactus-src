package src

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "time"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    "strings"
    "golang.org/x/crypto/bcrypt"
    "fmt"
)

type User struct {
    ID          string `json:"id,omitempty"`
    Name        string `json:"name,omitempty"`
    LastName    string `json:"last_name,omitempty"`
    Password    string `json:"password,omitempty"`
    Phone       string `json:"phone,omitempty"`
    Email       string `json:"email,omitempty"`
}

type Token struct {
  Token string `json:"token,omitempty"`
  UserID string `json:"user_uuid,omitempty"`
}

type Err struct {
  Error string `json:"error,omitempty"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LogIn(w http.ResponseWriter, r *http.Request) {
    var newUser User

    fmt.Println("11")

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newUser); err != nil {
        return
    }

    defer r.Body.Close()
    var password = newUser.Password

    sqlStatement := `SELECT id, name, last_name, phone, email, password FROM users WHERE email = $1`

    row := db.QueryRow(sqlStatement, newUser.Email)

    defer db.Close()

    resp := row.Scan(&newUser.ID, &newUser.Name, &newUser.LastName, &newUser.Phone, &newUser.Email, &newUser.Password);

    if resp != nil {
        http.Error(w, "{\"error\":\"User not found\"}", http.StatusNotFound)

        return
    }

    match := CheckPasswordHash(password, newUser.Password)

    if !match {
      http.Error(w, "{\"error\":\"Wrong email or password\"}", http.StatusNotFound)
      return
    }

    token := jwt.New(jwt.SigningMethodHS256)

    /* Set token claims */
    claims := token.Claims.(jwt.MapClaims)
    claims["admin"] = true
    claims["name"] = newUser.Name
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    /* Sign the token with our secret */
    tokenString, _ := token.SignedString(mySigningKey)

    var response Token
    response.Token = tokenString
    response.UserID = newUser.ID

    json.NewEncoder(w).Encode(response)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var newUser User

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newUser); err != nil {
        return
    }

    defer r.Body.Close()

    hashedPassword, _ := HashPassword(newUser.Password)

    sqlStatement := `
            INSERT INTO users (name, last_name, password, phone, email)
            VALUES ($1, $2, $3, $4, $5)`

    res, err := db.Exec(sqlStatement, newUser.Name, newUser.LastName, hashedPassword, newUser.Phone, newUser.Email);
    defer db.Close()

    if err != nil {
        var e = strings.Replace(err.Error(), "\"", "", -1)
        http.Error(w, "{\"error\":\"" + e +"\"}", http.StatusNotFound)
        return
    }

    count, err := res.RowsAffected()

    if (count < 1) {
        http.Error(w, "{\"error\":\"User not created\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newUser)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    sqlStatement := `DELETE FROM users WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["user_uuid"]
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
        http.Error(w, "{\"error\":\"User not found\"}", http.StatusNotFound)
    } else {
        http.Error(w, "{\"message\":\"User was deleted\"}", http.StatusOK)
    }
}

func GetUserByUUID(w http.ResponseWriter, r *http.Request) {
    var newUser User

    sqlStatement := `SELECT id, name, last_name, password, phone, email FROM users WHERE id = $1`
    vars := mux.Vars(r)
    var id = vars["user_uuid"]
    db, err := InitDB()
    row := db.QueryRow(sqlStatement, id)
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    resp := row.Scan(&id, &newUser.Name, &newUser.LastName, &newUser.Password, &newUser.Phone, &newUser.Email);

    if resp != nil {
        http.Error(w, "{\"error\":\"User not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newUser)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}

func IndexUser(w http.ResponseWriter, r *http.Request) {
    var resp error

    db, err := InitDB()

    if err != nil {
        log.Panic(err)
    }

    rows, err := db.Query("SELECT id, name, last_name, password, phone, email FROM users")
    defer db.Close()

    if err != nil {
        log.Panic(err)
    }

    defer rows.Close()

    users := make([]*User, 0)

    for rows.Next() {
        user := new(User)
        err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Password, &user.Phone, &user.Email)
        if err != nil {
            log.Panic(err)
        }
        users = append(users, user)
    }

    if err = rows.Err(); err != nil {
        log.Panic(err)
    }

    if resp != nil {
        http.Error(w, "{\"error\":\"User not found\"}", http.StatusNotFound)

        return
    } else {
        json.NewEncoder(w).Encode(users)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    }
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    var newUser User

    sqlStatement := `UPDATE users   SET name = $1, last_name = $2, password = $3, phone = $4, email = $5    WHERE id = $6;`
    vars := mux.Vars(r)
    var id = vars["user_uuid"]
    db, err := InitDB()

    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&newUser); err != nil {
        return
    }

    defer r.Body.Close()

    res, err := db.Exec(sqlStatement, newUser.Name, newUser.LastName, newUser.Password, newUser.Phone, newUser.Email, id)
    defer db.Close()

    if err != nil {
        panic(err)
    }
    count, err := res.RowsAffected()

    if err != nil {
        panic(err)
    }

    if (count < 1) {
        http.Error(w, "{\"error\":\"User not found\"}", http.StatusNotFound)

        return
    }

    json.NewEncoder(w).Encode(newUser)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
