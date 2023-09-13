package repository

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"service/config"
	"service/internal/domain"
)

type db struct {
	conn *sql.DB
}

func NewDB(config config.Config) (*db, error) {
	conn, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		return nil, err
	}

	return &db{
		conn: conn,
	}, nil
}

// Метод "Стоп" нужен для того, чтобы, если найден объект подключения, он сразу закрывался.
// Это нужно для работы с зависимостями
func (db *db) Stop() {
	if db.conn != nil {
		db.conn.Close()
	}
}

// Если в натс-стриминг поступают данные, WriteOrder отправляет их в кэш и базу данных
func (db *db) WriteOrder(order *domain.Order) error {
	// Формируем транзакцию
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	deliveryID := uuid.New()

	_, err = tx.Exec(`INSERT INTO delivery VALUES($1, $2, $3, $4, $5, $6, $7, $8)`, deliveryID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		// Если возникла ошибка, откатываем назад
		tx.Rollback()
		return err
	}

	paymentID := uuid.New()

	_, err = tx.Exec(`INSERT INTO payment VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, paymentID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`INSERT INTO "order" VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, order.UID, order.TrackNumber, order.Entry, deliveryID, paymentID, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`INSERT INTO item VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, item.ChrtID, order.UID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// Если всё хорошо, коммитим
	tx.Commit()

	return nil
}

// Этот метод вызывается только при старте сервера и всё, что есть в базе данных, помещает в кэш
func (db *db) GetAllOrders() ([]domain.Order, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, err
	}

	orders := []domain.Order{}
	/* Запускаем запрос SELECT к таблице "ордер", а также
	ко всем полям таблиц деливери и пеймент, а затем
	джойним две последние с ордером по внешним ключам.
	Поскольку связь "one-to-one", можно делать такие запросы */
	rows, err := tx.Query(`SELECT o.order_uid,
		o.track_number,
		o.entry,
		o.locale,
		o.internal_signature,
		o.customer_id,
		o.delivery_service,
		o.shardkey,
		o.sm_id,
		o.date_created,
		o.oof_shard,
		d.*,
		p.*
	FROM "order" o
		JOIN delivery d ON d.id = o.delivery_id
		JOIN payment p ON p.id = o.payment_id`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close() // Откладываем закрытие запроса выше для корректной работы базы данных

	// Для каждой следующей строки результата создаём объект order, чтобы от указателя получить данные по Delivery и Payment
	for rows.Next() {
		order := domain.Order{
			Delivery: &domain.Delivery{},
			Payment:  &domain.Payment{},
		}

		// Теперь сканируем все 30 аргументов в структуру
		err := rows.Scan(&order.UID, &order.TrackNumber, &order.Entry,
			&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.ID, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.ID, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Всё, что отсканировали, добавляем к списку заказов
		orders = append(orders, order)
	}
	// Для каждого списка заказов находим все айтемы, связанные с ним
	for i, order := range orders {
		rows, err := tx.Query(`SELECT * FROM item WHERE order_uid = $1`, order.UID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			item := domain.Item{}
			// Для каждого айтема так же сканируем всё, что получили
			err := rows.Scan(&item.ChrtID, &item.OrderID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// И добавляем к ордеру
			orders[i].Items = append(orders[i].Items, item)
		}
	}

	tx.Commit()
	// После завершения операций и подтверждения транзакции возвращаем список заказов
	return orders, nil
}
