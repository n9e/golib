// this is a sample echo rest api module
package user

import (
	"strings"

	"github.com/yubo/golib/orm"
	"github.com/yubo/golib/util"
)

const (
	CREATE_TABLE_SQLITE = "CREATE TABLE IF NOT EXISTS `user` (" +
		"  `id`    integer      PRIMARY    KEY AUTOINCREMENT," +
		"  `name`  varchar(128) DEFAULT '' NOT NULL," +
		"  `phone` varchar(16)  DEFAULT '' NOT NULL" +
		");" +
		" CREATE UNIQUE INDEX `user_index_name` on `user` (`name`);" +
		" CREATE UNIQUE INDEX `user_index_phone` on `user` (`phone`);"
)

func createUser(db *orm.Db, in *CreateUserInput) (int64, error) {
	return db.InsertLastId("user", in.Body)
}

func genUserSql(in *GetUsersInput) (where string, args []interface{}) {
	a := []string{}
	b := []interface{}{}
	if query := util.StringValue(in.Query); query != "" {
		a = append(a, "name like ?")
		b = append(b, "%"+query+"%")
	}
	if len(a) > 0 {
		where = " where " + strings.Join(a, " and ")
		args = b
	}
	return
}

func getUsers(db *orm.Db, in *GetUsersInput) (total int, list []*User, err error) {
	sql, args := genUserSql(in)

	err = db.Query("select count(*) from user "+sql, args...).Row(&total)
	if in.Count {
		return
	}

	err = db.Query("select * from user"+sql+in.SqlExtra("id desc"), args...).Rows(&list)
	return
}

func getUser(db *orm.Db, name string) (ret *User, err error) {
	err = db.Query("select * from user where name = ?", name).Row(&ret)
	return
}

func updateUser(db *orm.Db, in *UpdateUserInput) (*User, error) {
	if err := db.Update("user", in); err != nil {
		return nil, err
	}
	return getUser(db, in.Name)
}

func deleteUser(db *orm.Db, name string) (ret *User, err error) {
	if ret, err = getUser(db, name); err != nil {
		return
	}
	err = db.ExecNumErr("delete from user where name = ?", name)
	return

}
