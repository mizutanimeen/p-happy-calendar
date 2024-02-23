package mysql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mizutanimeen/P-happiness-100-strikes/internal/db/model"
)

var (
	machineTable  = os.Getenv("MYSQL_MACHINE_TABLE")
	machineID     = os.Getenv("MYSQL_MACHINE_ID")
	machineUserID = os.Getenv("MYSQL_USERS_ID")
	machineName   = os.Getenv("MYSQL_MACHINE_NAME")
	machineRate   = os.Getenv("MYSQL_MACHINE_RATE")
)

func (s *Mysql) MachinesGet(userID string) ([]*model.Machine, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", machineTable, machineUserID)
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error query: %w", err)
	}
	defer rows.Close()

	var machines []*model.Machine
	for rows.Next() {
		var machine model.Machine
		if err := rows.Scan(&machine.ID, &machine.UserID, &machine.Name, &machine.Rate, &machine.Create_at, &machine.Update_at); err != nil {
			return nil, fmt.Errorf("error scan: %w", err)
		}
		machines = append(machines, &machine)
	}
	return machines, nil
}

func (s *Mysql) MachineGetByID(id string, userID string) (*model.Machine, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? AND %s = ?", machineTable, machineID, machineUserID)
	row := s.DB.QueryRow(query, id, userID)

	var machine model.Machine
	if err := row.Scan(&machine.ID, &machine.UserID, &machine.Name, &machine.Rate, &machine.Create_at, &machine.Update_at); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scan: %w", err)
	}
	return &machine, nil
}

func (s *Mysql) MachineCreate(userID string, name string, rate int) error {
	query := fmt.Sprintf("INSERT INTO %s(%s, %s, %s) VALUES(?,?,?)", machineTable, machineUserID, machineName, machineRate)
	insert, err := s.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error prepare: %w", err)
	}

	if _, err := insert.Exec(userID, name, rate); err != nil {
		return fmt.Errorf("error exec: %w", err)
	}
	return nil
}

func (s *Mysql) MachineUpdate(id string, name string, rate int) error {
	query := fmt.Sprintf("UPDATE %s SET %s=?, %s=? WHERE %s=?", machineTable, machineName, machineRate, machineID)
	update, err := s.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error prepare: %w", err)
	}

	if _, err := update.Exec(name, rate, id); err != nil {
		return fmt.Errorf("error exec: %w", err)
	}
	return nil
}

func (s *Mysql) MachineDelete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", machineTable, machineID)
	delete, err := s.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error prepare: %w", err)
	}

	if _, err := delete.Exec(id); err != nil {
		return fmt.Errorf("error exec: %w", err)
	}
	return nil
}
