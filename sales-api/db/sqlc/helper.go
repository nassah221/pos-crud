package db

import (
	"context"
	"fmt"
)

func (q *Queries) CreateCashierWithReturn(ctx context.Context, arg CreateCashierParams) (*Cashier, error) {
	dbResult, err := q.db.ExecContext(ctx, createCashier, arg.Name, arg.Password)
	if err != nil {
		return nil, fmt.Errorf("cannot create cashier: %v", err)
	}

	insertID, err := dbResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch last insert id: %v", err)
	}

	cashier, err := q.DetailCashier(context.Background(), int32(insertID))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch cashier with id %d: %v", insertID, err)
	}

	return &cashier, nil
}

func (q *Queries) CreateProductWithReturn(ctx context.Context, arg CreateProductParams) (*Product, error) {
	dbResult, err := q.db.ExecContext(ctx, createProduct,
		arg.Name,
		arg.ImageUrl,
		arg.Price,
		arg.Stock,
		arg.CategoryID,
		arg.Sku,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create product: %v", err)
	}

	insertID, err := dbResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch last insert id: %v", err)
	}

	product, err := q.GetProduct(context.Background(), int32(insertID))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch product with id %d: %v", insertID, err)
	}

	return &product, nil
}

func (q *Queries) CreateDiscountWithReturn(ctx context.Context, arg CreateDiscountParams) (*Discount, error) {
	dbResult, err := q.db.ExecContext(ctx, createDiscount,
		arg.Qty,
		arg.Type,
		arg.Result,
		arg.ExpiredAt,
		arg.ExpiredAtFormat,
		arg.StringFormat,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create discount: %v", err)
	}

	insertID, err := dbResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch last insert id: %v", err)
	}

	discount, err := q.GetDiscount(context.Background(), int32(insertID))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch product with id %d: %v", insertID, err)
	}

	return &discount, nil
}

func (q *Queries) CreateCategoryWithReturn(ctx context.Context, name string) (*Category, error) {
	dbResult, err := q.db.ExecContext(ctx, createCategory, name)
	if err != nil {
		return nil, fmt.Errorf("cannot create category: %v", err)
	}

	insertID, err := dbResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch last insert id: %v", err)
	}

	category, err := q.GetCategory(context.Background(), int32(insertID))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch category with id %d: %v", insertID, err)
	}

	return &category, nil
}

func (q *Queries) CreatePaymentWithReturn(ctx context.Context, arg CreatePaymentParams) (*Payment, error) {
	dbResult, err := q.db.ExecContext(ctx, createPayment, arg.Name, arg.Type, arg.Logo)
	if err != nil {
		return nil, fmt.Errorf("cannot create payment: %v", err)
	}

	insertID, err := dbResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch last insert id: %v", err)
	}

	payment, err := q.GetPayment(context.Background(), int32(insertID))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch payment with id %d: %v", insertID, err)
	}

	return &payment, nil
}

const (
	sqlMode = `SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));`
)

func (q *Queries) SetSQLMode(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, sqlMode)
	return err
}
