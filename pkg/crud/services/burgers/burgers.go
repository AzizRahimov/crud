package burgers

import (
	"context"
	"crud/pkg/crud/models"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BurgersSvc struct {
	pool *pgxpool.Pool // dependency
}

func NewBurgersSvc(pool *pgxpool.Pool) *BurgersSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil"))
	}
	return &BurgersSvc{pool: pool}
}

func (service *BurgersSvc) InitDB() error {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewQueryError("can't init db: ", err)
	}
	_, err = conn.Query(context.Background(), BurgersDDL)
	if err != nil {
		return NewQueryError("can't init db: ", err)
	}
	return nil
}

func (service *BurgersSvc) BurgersList() (list []models.Burger, err error) {
	list = make([]models.Burger, 0) // TODO: for REST API
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return nil, NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), GetBurgers)
	if err != nil {
		return nil, NewQueryError(GetBurgers, err) // TODO: wrap to specific error
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Burger{}
		err := rows.Scan(&item.Id, &item.Name, &item.Price)
		if err != nil {
			return nil, NewDbError(err) // TODO: wrap to specific error
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, NewDbError(err)
	}

	return list, nil
}

func (service *BurgersSvc) Save(model models.Burger) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), SaveBurger, model.Name, model.Price)
	if err != nil {
		return NewQueryError(SaveBurger, err)
	}

	return nil
}

func (service *BurgersSvc) RemoveById(id int) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), RemoveBurger, id)
	if err != nil {
		return NewQueryError(RemoveBurger, err)
	}

	return nil
}
