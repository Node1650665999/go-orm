package session

import "geeorm/log"

// Begin a transaction
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	//调用 s.db.Begin() 得到 *sql.Tx 对象，赋值给 s.tx
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

// Commit a transaction
func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

// Rollback a transaction
func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
