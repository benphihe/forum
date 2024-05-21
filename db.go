package forum
import (
	"database/sql"
	"log"
	"github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "path/to/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Rest of your code here
}
