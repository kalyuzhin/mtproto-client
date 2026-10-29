package postgresql

import (
	"context"
	"fmt"

	"github.com/kalyuzhin/mtproto-client/internal/pkg/domain"
	"github.com/samber/lo"
)

type userDB struct {
	accountID   int64  `db:"account_id"`
	userID      int64  `db:"user_id"`
	firstName   string `db:"first_name"`
	lastName    string `db:"last_name"`
	username    string `db:"username"`
	phoneNumber string `db:"phone_number"`
}

// GetUsersByID – получить всех пользователей, которые записаны у определенного аккаунта.
func (d *Database) GetUsersByID(ctx context.Context, accountID int64) {

}

// InsertUsers – внести пользователей в базу данных для конкретного аккаунта (accountID)
func (d *Database) InsertUsers(ctx context.Context, accountID int64, users []domain.User) (int64, error) {
	const query = `
	INSERT INTO tg_users(account_id,
            user_id,
            first_name,
            last_name,
            username,
            phone_number)
	VALUES(...)
	ON CONFLICT (account_id, user_id)
		DO UPDATE SET
    		first_name = EXCLUDED.first_name,
    		last_name = EXCLUDED.last_name,
    		username = EXCLUDED.username,
    		phone_number = EXCLUDED.phone_number;`

	usersDTO := mapDomainUserToDTO(users, accountID)

	rows, err := d.Exec(ctx, query, usersDTO)
	if err != nil {
		return 0, fmt.Errorf("d.InsertUsers.Exec: %w", err)
	}

	return rows.RowsAffected(), nil
}

func mapDomainUserToDTO(users []domain.User, accountID int64) []userDB {
	return lo.Map(users, func(user domain.User, idx int) userDB {
		return userDB{
			accountID:   accountID,
			userID:      user.GetUserID(),
			firstName:   user.GetFirstName(),
			lastName:    user.GetLastName(),
			phoneNumber: user.GetPhoneNumber(),
		}
	})
}
