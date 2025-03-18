// main.go
package main

import (
	"bufio"
	"database/sql"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq" 
)


type Email struct {
	ID         int      
	Username   string    
	Recipient  string   
	RawData    string    
	ReceivedAt time.Time 
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	
	db := connectDB()
	// Start the cleanup job to purge emails older than 7 days.
	go startCleanupJob(db)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	// Start listening for incoming SMTP 
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	log.Println("Temporary Mail Service SMTP Server listening on port "+port)

	
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		
		go handleConnection(conn, db)
	}
}


func handleConnection(conn net.Conn, db *sql.DB) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	
	writer.WriteString("220 Welcome to Temporary Mail Service\r\n")
	writer.Flush()

	var rcptTo string 

	// Process commands in a continuous loop.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			} else {
				log.Println("Error reading command:", err)
			}
			return
		}

		
		input := strings.ToUpper(strings.TrimSpace(line))
		log.Printf("Received: %s", input)

		// Use a switch statement to handle various SMTP commands.
		switch {
		case input == "VRFY":
			writer.WriteString("250 VRFY not supported\r\n")
		case input == "NOOP":
			writer.WriteString("250 OK\r\n")
		case input == "EXPN":
			writer.WriteString("250 EXPN not supported\r\n")
		case input == "HELP":
			writer.WriteString("214-Commands supported:\r\n")
			writer.WriteString("214 HELO, EHLO, MAIL FROM, RCPT TO, DATA, NOOP, VRFY, HELP, EXPN, QUIT\r\n")
		case strings.HasPrefix(input, "HELO") || strings.HasPrefix(input, "EHLO"):
			writer.WriteString("250 Hello\r\n")
		case strings.HasPrefix(input, "MAIL FROM:"):
			writer.WriteString("250 Sender OK\r\n")
		case strings.HasPrefix(input, "RCPT TO:"):
			// Save the recipient address for later extraction.
			rcptTo = strings.TrimSpace(line[len("RCPT TO:"):])
			writer.WriteString("250 Recipient OK\r\n")
		case input == "DATA":
			writer.WriteString("354 End data with <CR><LF>.<CR><LF>\r\n")
			writer.Flush()

			
			var dataLines []string
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println("Error reading email data:", err)
					return
				}
				if strings.TrimSpace(dataLine) == "." {
					break
				}
				dataLines = append(dataLines, dataLine)
			}
			rawEmail := strings.Join(dataLines, "")
			log.Println("Received raw email data:")
			log.Println(rawEmail)

			
			username := extractUsername(rcptTo)
			
			emailRecord := Email{
				Username:   username,
				Recipient:  rcptTo,
				RawData:    rawEmail,
				ReceivedAt: time.Now(),
			}
			// Store the email in the database.
			if err := storeEmail(db, emailRecord); err != nil {
				log.Println("Error storing email:", err)
				writer.WriteString("550 Error storing email\r\n")
			} else {
				writer.WriteString("250 OK: Message accepted\r\n")
			}
		case input == "QUIT":
			writer.WriteString("221 Bye\r\n")
			writer.Flush()
			return
		default:
			writer.WriteString("500 Unrecognized command\r\n")
		}
		writer.Flush()
	}
}


func extractUsername(rcpt string) string {
	rcpt = strings.Trim(rcpt, "<>")
	parts := strings.Split(rcpt, "@")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// connectDB connects to the Supabase/PostgreSQL database.
func connectDB() *sql.DB {
	connStr := os.Getenv("SUPABASE_CONN_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}
	return db
}


func storeEmail(db *sql.DB, email Email) error {
	query := `
        INSERT INTO emails (username, recipient, raw_data, received_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, email.Username, email.Recipient, email.RawData, email.ReceivedAt)
	return err
}


func startCleanupJob(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour) // Adjust cleanup frequency as needed.
	go func() {
		for range ticker.C {
			_, err := db.Exec(`DELETE FROM emails WHERE received_at < NOW() - INTERVAL '7 days'`)
			if err != nil {
				log.Println("Error cleaning up old emails:", err)
			} else {
				log.Println("Cleanup job: Old emails removed.")
			}
		}
	}()
}
