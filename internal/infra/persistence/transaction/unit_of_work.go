package transaction

import (
	"context"
	"database/sql"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

type UnitOfWorkInterface interface {
	Executer(ctx context.Context, fn func(exec persistence.Executer) error) error
}

type UnitOfWork struct {
	DB *sql.DB
}

func (u *UnitOfWork) Execute(ctx context.Context, fn func(exec persistence.Executer) error) (err error) {
	tx, err := u.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
