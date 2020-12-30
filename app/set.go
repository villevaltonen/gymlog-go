package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// set is a basic entity for holding sets
type set struct {
	ID          int     `json:"id"`
	UserID      string  `json:"userId"`
	Weight      float64 `json:"weight"`
	Exercise    string  `json:"exercise"`
	Repetitions int     `json:"repetitions"`
}

func (s *Server) handleGetSet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusBadRequest, "Invalid set ID")
			return
		}

		set := set{ID: id}
		if err := set.getSet(s.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
				log.Println(err.Error())
				respondWithError(w, http.StatusNotFound, "Set not found")
			default:
				log.Println(err.Error())
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		respondWithJSON(w, http.StatusOK, s)
	}
}

func (s *Server) handleGetSets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var set set
		count, _ := strconv.Atoi(r.FormValue("count"))
		start, _ := strconv.Atoi(r.FormValue("start"))

		if count > 10 || count < 1 {
			count = 10
		}
		if start < 0 {
			start = 0
		}

		sets, err := set.getSets(s.DB, start, count)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, sets)
	}
}

func (s *Server) handleCreateSet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var set set
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&set); err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		if err := set.createSet(s.DB); err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, set)
	}
}

func (s *Server) handleUpdateSet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusBadRequest, "Invalid set ID")
			return
		}

		var set set
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&set); err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		set.ID = id

		if err := set.updateSet(s.DB); err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, set)
	}

}

func (s *Server) handleDeleteSet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusBadRequest, "Invalid Set ID")
			return
		}

		set := set{ID: id}
		if err := set.deleteSet(s.DB); err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	}
}

// getSet fetches a set from database with id
func (s *set) getSet(db *sql.DB) error {
	return db.QueryRow("SELECT user_id, weight, exercise, repetitions FROM sets WHERE id=$1",
		s.ID).Scan(&s.UserID, &s.Weight, &s.Exercise, &s.Repetitions)
}

// getSets fetches multiple sets from database with user id
func (s *set) getSets(db *sql.DB, start, count int) ([]set, error) {
	rows, err := db.Query(
		"SELECT id, user_id, weight, exercise, repetitions FROM sets LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sets := []set{}

	for rows.Next() {
		var s set
		if err := rows.Scan(&s.ID, &s.UserID, &s.Weight, s.Exercise, s.Repetitions); err != nil {
			return nil, err
		}
		sets = append(sets, s)
	}

	return sets, nil
}

// updateSet executes update query to database
func (s *set) updateSet(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE sets SET user_id=$2, weight=$3, exercise=$4, repetitions=$5 WHERE id=$1",
			s.ID, s.UserID, s.Weight, s.Exercise, s.Repetitions)

	return err
}

// deleteSet deletes a set from database with
func (s *set) deleteSet(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM sets WHERE id=$1", s.ID)

	return err
}

// createSet creates a set into database with given JSON
func (s *set) createSet(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO sets(user_id, weight, exercise, repetitions) VALUES($1, $2, $3, $4) RETURNING id",
		s.UserID, s.Weight, s.Exercise, s.Repetitions).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}