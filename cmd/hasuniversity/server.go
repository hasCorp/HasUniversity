package hasuniversity

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/hascorp/hasuniversity/route/flashcards"
	"github.com/hascorp/hasuniversity/route/ping"
)

const (
	keyspace               = "hasuniversity"
	allowedRetries         = 15 // total number of retries allowed
	allowedAttemptsAfterUp = 5  // allow for more retries if keyspace/table doesn't exist yet
	retryBase              = 1 * time.Second
)

func Start(port int, cassandra string) {
	if port < 1 {
		log.Fatalf("Invalid port: %d", port)
	}
	addr := fmt.Sprintf(":%d", port)

	ips := []string{}
	for _, s := range strings.Split(cassandra, ",") {
		ips = append(ips, strings.TrimSpace(s))
	}

	session, err := cqlClient(ips)
	if err != nil {
		log.Fatalf("Failed to create Cassandra client session: %v", err)
	}
	defer session.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", ping.PingHandler).Methods("GET")

	// flashcard APIs
	flashcardHandler := flashcards.FlashcardHandler{
		Session: session,
	}
	flashcardRoute := r.PathPrefix("/flashcard").Subrouter()

	flashcardRoute.HandleFunc("/", flashcardHandler.GetAllFlashcards).Methods("GET")
	flashcardRoute.HandleFunc("", flashcardHandler.GetAllFlashcards).Methods("GET")

	flashcardRoute.HandleFunc("/", flashcardHandler.AddFlashcardSet).Methods("POST")
	flashcardRoute.HandleFunc("", flashcardHandler.AddFlashcardSet).Methods("POST")

	flashcardRoute.HandleFunc("/{uuid:[0-9A-Za-z]+}", flashcardHandler.GetFlashcardSet).Methods("GET")

	http.Handle("/", r)

	srv := &http.Server{
		Handler:           r,
		Addr:              addr,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Printf("Starting server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func cqlClient(connections []string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(connections...)
	cluster.Keyspace = keyspace

	var err error
	var session *gocql.Session
	attemptsAfterUp := 0
	for attempts := 0; attempts < allowedRetries; attempts++ {
		if attemptsAfterUp >= allowedAttemptsAfterUp {
			log.Printf("Cassandra is available but db isn't set up properly after %d attempts\n", attemptsAfterUp)
			return nil, errors.New("db unavailable")
		}

		t := retryBase * (time.Duration(1 << attempts))
		session, err = cluster.CreateSession()
		if err == nil {
			return session, nil
		} else if strings.Contains(err.Error(), "no connections were made when creating the session") {
			attemptsAfterUp++
		} else if !strings.Contains(err.Error(), "connection refused") {
			log.Println("Fatal error connecting to Cassandra: ", err)
			return nil, err
		}
		// connection refused means that the Cassandra db is not ready yet,
		// so wait and retry
		log.Printf("After %d attempt(s), Cassandra not available yet. Sleeping for %v seconds\n", attempts+1, t.Seconds())
		time.Sleep(t)
	}

	if err == nil {
		err = errors.New("failed to connect to cassandra unexpectedly")
	}

	return nil, err
}
