package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
)

func userParamToSelectQuery(p *common.UserParam) (q sqlQuery) {
	q.query = fmt.Sprintf("SELECT %s FROM %s ", userSelectColumnsFull, tableUsers)
	counter := 1
	if p.Max > 0 {
		q.query += fmt.Sprintf("LIMIT $%d ", counter)
		q.params = append(q.params, p.Max)
		counter++
	}
	if p.Offset > 0 {
		q.query += fmt.Sprintf("OFFSET $%d ", counter)
		q.params = append(q.params, p.Offset)
		counter++
	}
	return
}

func (db *Database) UserFollows(a, b uint32) (follows bool, err error) {
	err = db.queryWithPrepared(queryUserFollows, func(rows *sql.Rows) error {
		follows = true
		return nil
	}, a, b)
	return
}

func (db *Database) AddUserFollowing(a, b uint32) (err error) {
	_, err = db.getPrepared(queryUserFollowsUpsert).Exec(a, b)
	return
}

func (db *Database) DeleteUserFollowing(a, b uint32) (deleted bool, err error) {
	var affected uint32
	affected, err = db.execQuery(queryDeleteUserFollowing, a, b)
	deleted = affected > 0
	return
}

func (db *Database) getUserByQuery(name string, p interface{}) (user model.User, has bool, err error) {
	err = db.queryWithPrepared(name, func(rows *sql.Rows) error {
		rows.Next()
		scanUserColumnsFull(rows, &user)
		has = true
		return nil
	}, p)
	return
}

func (db *Database) GetUserByAPIToken(token string) (user model.User, has bool, err error) {
	user, has, err = db.getUserByQuery(queryGetUserByApiToken, token)
	return
}

func (db *Database) GetUsersByEmail(email string) (users []model.User, err error) {
	err = db.queryWithPrepared(queryGetUserByEmail, func(rows *sql.Rows) error {
		for rows.Next() {
			var user model.User
			scanUserColumnsFull(rows, &user)
			users = append(users, user)
		}
		return nil
	}, email)
	return
}

func (db *Database) GetUserByName(name string) (user model.User, has bool, err error) {
	user, has, err = db.getUserByQuery(queryGetUserByName, name)
	return
}

func (db *Database) GetUserByID(id uint32) (user model.User, has bool, err error) {
	user, has, err = db.getUserByQuery(queryGetUserByID, id)
	return
}

func (db *Database) InsertUser(u *model.User) (err error) {
	_, err = db.getPrepared(queryInsertUser).Exec(u.Username, u.Password, u.Email, u.Status, u.CreatedAt, u.UpdatedAt, u.LastLoginAt, u.LastLoginIP, u.Token, u.TokenExpiration, u.Language, u.MD5)
	return
}

func (db *Database) UpdateUser(u *model.User) (err error) {
	_, err = db.getPrepared(queryUpdateUser).Exec(u.ID, u.Username, u.Password, u.Email, u.Status, u.UpdatedAt, u.LastLoginAt, u.LastLoginIP, u.Token, u.TokenExpiration, u.Language, u.MD5)
	return
}

func (db *Database) GetUsersWhere(param *common.UserParam) (users []model.User, err error) {
	var user model.User
	var has bool
	if len(param.Email) > 0 {
		users, err = db.GetUsersByEmail(param.Email)
	} else if len(param.Name) > 0 {
		user, has, err = db.GetUserByName(param.Name)
		if has {
			users = append(users, user)
		}
	} else if len(param.ApiToken) > 0 {
		user, has, err = db.GetUserByAPIToken(param.ApiToken)
		if has {
			users = append(users, user)
		}
	} else if param.ID > 0 {
		user, has, err = db.GetUserByID(param.ID)
		if has {
			users = append(users, user)
		}
	} else {
		q := userParamToSelectQuery(param)
		if param.Max > 0 {
			users = make([]model.User, 0, param.Max)
		} else {
			users = make([]model.User, 0, 64)
		}
		err = q.Query(db.conn, func(rows *sql.Rows) error {

			for rows.Next() {
				var user model.User
				scanUserColumnsFull(rows, &user)
				users = append(users, user)
			}
			return nil
		})
	}
	return
}

func (db *Database) DeleteUsersWhere(param *common.UserParam) (deleted uint32, err error) {

	var queryName string
	var p interface{}
	if param.ID > 0 {
		queryName = queryDeleteUserByID
		p = param.ID
	} else if len(param.Email) > 0 {
		queryName = queryDeleteUserByEmail
		p = param.Email
	} else if len(param.ApiToken) > 0 {
		queryName = queryDeleteUserByToken
		p = param.ApiToken
	} else {
		// delete nothing
		return
	}
	deleted, err = db.execQuery(queryName, p)
	return
}
