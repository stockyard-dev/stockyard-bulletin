package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ *sql.DB }
type Subscriber struct {
	ID          int64      `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Tags        string     `json:"tags"`
	Status      string     `json:"status"`
	Token       string     `json:"token,omitempty"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
type Campaign struct {
	ID        int64      `json:"id"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	Status    string     `json:"status"`
	SentCount int        `json:"sent_count"`
	OpenCount int        `json:"open_count"`
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	dsn := filepath.Join(dataDir, "bulletin.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(resource TEXT NOT NULL,record_id TEXT NOT NULL,data TEXT NOT NULL DEFAULT '{}',PRIMARY KEY(resource, record_id))`)
	return &DB{db}, nil
}
func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS subscribers(id INTEGER PRIMARY KEY AUTOINCREMENT,email TEXT NOT NULL UNIQUE,name TEXT DEFAULT '',tags TEXT DEFAULT '',status TEXT DEFAULT 'active',token TEXT DEFAULT '',confirmed_at DATETIME,created_at DATETIME DEFAULT CURRENT_TIMESTAMP);CREATE TABLE IF NOT EXISTS campaigns(id INTEGER PRIMARY KEY AUTOINCREMENT,subject TEXT NOT NULL,body TEXT NOT NULL,status TEXT DEFAULT 'draft',sent_count INTEGER DEFAULT 0,open_count INTEGER DEFAULT 0,sent_at DATETIME,created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`)
	return err
}
func (db *DB) ListSubscribers(status string) ([]Subscriber, error) {
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = db.Query(`SELECT id,email,name,tags,status,confirmed_at,created_at FROM subscribers WHERE status=? ORDER BY created_at DESC`, status)
	} else {
		rows, err = db.Query(`SELECT id,email,name,tags,status,confirmed_at,created_at FROM subscribers ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subscriber
	for rows.Next() {
		var s Subscriber
		rows.Scan(&s.ID, &s.Email, &s.Name, &s.Tags, &s.Status, &s.ConfirmedAt, &s.CreatedAt)
		out = append(out, s)
	}
	return out, nil
}
func (db *DB) Subscribe(s *Subscriber) error {
	if s.Status == "" {
		s.Status = "active"
	}
	if s.Token == "" {
		s.Token = fmt.Sprintf("%x", time.Now().UnixNano())
	}
	res, err := db.Exec(`INSERT INTO subscribers(email,name,tags,status,token)VALUES(?,?,?,?,?) ON CONFLICT(email) DO UPDATE SET status='active',name=excluded.name`, s.Email, s.Name, s.Tags, s.Status, s.Token)
	if err != nil {
		return err
	}
	s.ID, _ = res.LastInsertId()
	return nil
}
func (db *DB) UpdateSubscriberStatus(id int64, status string) error {
	_, err := db.Exec(`UPDATE subscribers SET status=? WHERE id=?`, status, id)
	return err
}
func (db *DB) DeleteSubscriber(id int64) error {
	_, err := db.Exec(`DELETE FROM subscribers WHERE id=?`, id)
	return err
}
func (db *DB) ListCampaigns() ([]Campaign, error) {
	rows, err := db.Query(`SELECT id,subject,body,status,sent_count,open_count,sent_at,created_at FROM campaigns ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Campaign
	for rows.Next() {
		var c Campaign
		rows.Scan(&c.ID, &c.Subject, &c.Body, &c.Status, &c.SentCount, &c.OpenCount, &c.SentAt, &c.CreatedAt)
		out = append(out, c)
	}
	return out, nil
}
func (db *DB) CreateCampaign(c *Campaign) error {
	if c.Status == "" {
		c.Status = "draft"
	}
	res, err := db.Exec(`INSERT INTO campaigns(subject,body,status)VALUES(?,?,?)`, c.Subject, c.Body, c.Status)
	if err != nil {
		return err
	}
	c.ID, _ = res.LastInsertId()
	return nil
}
func (db *DB) SendCampaign(id int64) (int, error) {
	var activeCount int
	db.QueryRow(`SELECT COUNT(*) FROM subscribers WHERE status='active'`).Scan(&activeCount)
	now := time.Now()
	_, err := db.Exec(`UPDATE campaigns SET status='sent',sent_count=?,sent_at=? WHERE id=?`, activeCount, now, id)
	return activeCount, err
}
func (db *DB) DeleteCampaign(id int64) error {
	_, err := db.Exec(`DELETE FROM campaigns WHERE id=?`, id)
	return err
}
func (db *DB) Stats() (map[string]int, error) {
	var subs, active, camps, sent int
	db.QueryRow(`SELECT COUNT(*) FROM subscribers`).Scan(&subs)
	db.QueryRow(`SELECT COUNT(*) FROM subscribers WHERE status='active'`).Scan(&active)
	db.QueryRow(`SELECT COUNT(*) FROM campaigns`).Scan(&camps)
	db.QueryRow(`SELECT COUNT(*) FROM campaigns WHERE status='sent'`).Scan(&sent)
	return map[string]int{"total_subscribers": subs, "active_subscribers": active, "campaigns": camps, "sent_campaigns": sent}, nil
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
