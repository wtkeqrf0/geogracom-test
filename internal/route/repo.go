package route

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO вынести моки в другую директорию

//go:generate ifacemaker -f *.go -o irepo.go -i Repository -s repo -p route
//go:generate mockgen -package route -source ./irepo.go -destination repo_mock.go
type repo struct {
	db     *pgxpool.Pool
	delete chan *Delete
}

func NewRepo(db *pgxpool.Pool) Repository {
	r := &repo{
		db:     db,
		delete: make(chan *Delete),
	}

	go r.listenDelete()
	return r
}

func (r *repo) Upsert(ctx context.Context, param *Upsert) (routeID int64, err error) {
	const query = `INSERT INTO route (id, name, load, cargo)
						VALUES($1, $2, $3, $4)
					ON CONFLICT (id) DO NOTHING`

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback(ctx)
		} else {
			txErr = tx.Commit(ctx)
		}

		if txErr != nil {
			log.Println(txErr.Error())
		}
	}()

	cmd, err := tx.Exec(ctx, query, param.ID, param.Name, param.Load, param.Cargo)
	if err != nil {
		return
	}

	routeID = param.ID
	if cmd.RowsAffected() == 0 {
		routeID, err = insert(ctx, &param.Insert, tx)
		if err != nil {
			return 0, err
		}
		err = unActual(ctx, []int64{param.ID}, tx)
	}
	return routeID, err
}

func insert(ctx context.Context, param *Insert, tx pgx.Tx) (int64, error) {
	const query = `INSERT INTO route(name, load, cargo)
						VALUES($1, $2, $3) RETURNING id`

	rows, err := tx.Query(ctx, query, param.Name, param.Load, param.Cargo)
	if err != nil {
		return 0, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowTo[int64])
}

func unActual(ctx context.Context, routeIDs []int64, tx pgx.Tx) error {
	const query = `UPDATE route SET
                 		actual = false
					WHERE id = any($1)`

	_, err := tx.Exec(ctx, query, routeIDs)
	return err
}

func (r *repo) Get(ctx context.Context, param *Get) (*Route, error) {
	const query = `SELECT 
    				id,
    				name,
    				load, 
    				cargo, 
    				actual
    			   FROM route
    			    WHERE id = $1 LIMIT 1`

	rows, err := r.db.Query(ctx, query, param.ID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Route])
}

func (r *repo) Delete(ctx context.Context, param *Delete) error {
	const query = `DELETE FROM route WHERE id = any($1)`

	_, err := r.db.Exec(ctx, query, param.Ids)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) DeleteBackground(param *Delete) {
	r.delete <- param
}

func (r *repo) listenDelete() {
	for param := range r.delete {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		if err := r.Delete(ctx, param); err != nil {
			// todo впихнуть defer
			cancel()
			log.Print(err.Error())
		}
		cancel()
	}
}
