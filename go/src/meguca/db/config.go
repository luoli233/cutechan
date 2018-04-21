package db

import (
	"database/sql"
	"encoding/json"

	"meguca/config"
	"meguca/util"
)

func loadServerConfig() error {
	conf, err := getServerConfig()
	if err != nil {
		return err
	}
	config.Set(conf)
	return listenFunc("config_updates", updateServerConfig)
}

func getServerConfig() (c config.ServerConfig, err error) {
	var data []byte
	err = db.QueryRow(`SELECT val FROM main WHERE id = 'config'`).Scan(&data)
	if err != nil {
		return
	}
	c, err = decodeServerConfig(data)
	return
}

func decodeServerConfig(data []byte) (c config.ServerConfig, err error) {
	err = json.Unmarshal(data, &c)
	return
}

func updateServerConfig(msg string) error {
	conf, err := decodeServerConfig([]byte(msg))
	if err != nil {
		return util.WrapError("reloading configuration", err)
	}
	config.Set(conf)
	return nil
}

func SetServerConfig(c config.ServerConfig) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return execPrepared("update_server_config", data)
}

func loadBoardConfigs() error {
	r, err := prepared["get_board_configs"].Query()
	if err != nil {
		return err
	}
	defer r.Close()

	for r.Next() {
		c, err := readBoardConfig(r)
		if err != nil {
			return err
		}
		if _, err := config.SetBoardConfig(c); err != nil {
			return err
		}
	}
	if err := r.Err(); err != nil {
		return err
	}

	return listenFunc("board_updated", updateBoardConfig)
}

func GetBoardConfig(board string) (config.BoardConfig, error) {
	return readBoardConfig(prepared["get_board_config"].QueryRow(board))
}

func readBoardConfig(r rowScanner) (c config.BoardConfig, err error) {
	var id string
	var modOnly bool
	var settings []byte
	if err = r.Scan(&id, &modOnly, &settings); err != nil {
		return
	}
	if err = json.Unmarshal(settings, &c); err != nil {
		return
	}
	c.ID = id
	c.ModOnly = modOnly
	return
}

func updateBoardConfig(board string) error {
	conf, err := GetBoardConfig(board)
	switch err {
	case nil:
		// Do nothing.
	case sql.ErrNoRows:
		config.RemoveBoard(board)
		return nil
	default:
		return err
	}
	if _, err = config.SetBoardConfig(conf); err != nil {
		return util.WrapError("reloading board configuration", err)
	}
	return nil
}

func WriteBoard(tx *sql.Tx, c config.BoardConfig) (err error) {
	settings, err := json.Marshal(c)
	if err != nil {
		return
	}
	_, err = getStatement(tx, "write_board").Exec(c.ID, c.ModOnly, settings)
	return
}

func UpdateBoard(c config.BoardConfig) (err error) {
	settings, err := json.Marshal(c)
	if err != nil {
		return
	}
	return execPrepared("update_board", c.ID, c.ModOnly, settings)
}
