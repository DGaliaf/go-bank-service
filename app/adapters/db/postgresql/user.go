package postgresql

import (
	"avito-tech/app/internal/domain/entity"
	custom_error "avito-tech/app/internal/errors"
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type UserStorage struct {
	db *pgx.Conn
}

func NewUserStorage(db *pgx.Conn) *UserStorage {
	return &UserStorage{db: db}
}

func (u UserStorage) Create(ctx context.Context) (int, error) {
	var userId int

	sql := `INSERT INTO public.user DEFAULT VALUES RETURNING "id"`

	if err := u.db.QueryRow(ctx, sql).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, custom_error.ErrNotFound
		}

		return -1, custom_error.ErrScanFail
	}

	return userId, nil
}

func (u UserStorage) FindOne(ctx context.Context, id int) (*entity.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Select("*").From("public.user").Where(sq.Eq{"id": id}).ToSql()

	user := entity.User{}

	if err := u.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, custom_error.ErrNotFound
		}
		return nil, custom_error.ErrScanFail
	}

	return &user, nil
}

func (u UserStorage) ChargeBalance(ctx context.Context, user *entity.User) error {
	oldUser, err := u.FindOne(ctx, user.Id)
	if err != nil {
		return err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Update("public.user").Where(sq.Eq{"id": user.Id}).Set("balance", oldUser.Balance+user.Balance).ToSql()

	res, err := u.db.Exec(ctx, sql, args...)
	if err != nil {
		return custom_error.ErrExecuteSQL
	}

	if res.RowsAffected() == 0 {
		return custom_error.ErrNotFound
	}

	return nil
}

func (u UserStorage) RemoveBalance(ctx context.Context, user *entity.User) error {
	oldUser, err := u.FindOne(ctx, user.Id)
	if err != nil {
		return err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	updatedBalance := oldUser.Balance - user.Balance
	if updatedBalance < 0 {
		return custom_error.ErrNegativeBalance
	}

	sql, args, _ := psql.Update("public.user").Where(sq.Eq{"id": user.Id}).Set("balance", updatedBalance).ToSql()

	res, err := u.db.Exec(ctx, sql, args...)
	if err != nil {
		return custom_error.ErrExecuteSQL
	}

	if res.RowsAffected() == 0 {
		return custom_error.ErrNotFound
	}

	return nil
}

func (u UserStorage) TransferMoney(ctx context.Context, data *entity.TransferMoney) error {
	tx, err := u.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return custom_error.ErrTransaction
	}
	defer tx.Rollback(ctx)

	fromUser := new(entity.User)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Select("*").From("public.user").Where(sq.Eq{"id": data.From}).ToSql()
	if err := tx.QueryRow(ctx, sql, args...).Scan(&fromUser.Id, &fromUser.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return custom_error.ErrNotFound
		}

		return custom_error.ErrScanFail
	}

	updatedBalance := fromUser.Balance - data.Amount
	if updatedBalance < 0 {
		return custom_error.ErrNegativeBalance
	}

	sql, args, _ = psql.Update("public.user").Where(sq.Eq{"id": data.From}).Set("balance", updatedBalance).ToSql()
	res, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return custom_error.ErrExecuteSQL
	}

	if res.RowsAffected() == 0 {
		return custom_error.ErrNotFound
	}

	toUser := new(entity.User)
	sql, args, _ = psql.Select("*").From("public.user").Where(sq.Eq{"id": data.To}).ToSql()
	if err := tx.QueryRow(ctx, sql, args...).Scan(&toUser.Id, &toUser.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return custom_error.ErrNotFound
		}

		return custom_error.ErrScanFail
	}

	updatedBalance = toUser.Balance + data.Amount
	sql, args, _ = psql.Update("public.user").Where(sq.Eq{"id": data.To}).Set("balance", updatedBalance).ToSql()
	res, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return custom_error.ErrExecuteSQL
	}

	if res.RowsAffected() == 0 {
		return custom_error.ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return custom_error.ErrTransaction
	}

	return nil
}
