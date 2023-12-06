package postgresql

import (
	mymodel "Subscriber/internal/model"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	connStr = "user=Timka password=root dbname=data_model sslmode=disable"
)

func New() map[string]mymodel.MyModel {
	mapa := make(map[string]mymodel.MyModel)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	rows, err := db.Query(
		`select public.order.order_uid from public.order;`,
	)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var uid string
		rows.Scan(&uid)
		mapa[uid] = getDataByUid(uid)
	}

	return mapa
}

// TODO: надо дописать всякие злоебучие проверки на то, что такая запись в базе существует и тд и тп
func getDataByUid(uid string) mymodel.MyModel {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	row := db.QueryRow(
		`select * from public.order
		join public.delivery on public.delivery.order_uid = public.order.order_uid
		join public.payment on public.payment.transaction = public.order.order_uid
		where public.order.order_uid = $1;`, uid)

	model := mymodel.MyModel{}
	err = row.Scan(
		&model.Order_uid, 
		&model.Track_number, 
		&model.Entry,
		&model.Locale,
		&model.Internal_signature,
		&model.Customer_id,
		&model.Delivery_service,
		&model.Shardkey,
		&model.Sm_id,
		&model.Date_created,
		&model.Oof_shard,
		&model.Delivery.Delivery_uuid,
		&model.Delivery.Name,
		&model.Delivery.Phone,
		&model.Delivery.Zip,
		&model.Delivery.City,
		&model.Delivery.Address,
		&model.Delivery.Region,
		&model.Delivery.Email,
		&model.Payment.Transaction,
		&model.Payment.Request_id,
		&model.Payment.Currency,
		&model.Payment.Provider,
		&model.Payment.Amount,
		&model.Payment.Payment_dt,
		&model.Payment.Bank,
		&model.Payment.Delivery_cost,
		&model.Payment.Goods_total,
		&model.Payment.Custom_fee,
	)
	if err != nil {
		log.Println(err)
	}

	track_number := model.Track_number

	rows, err := db.Query(`select * from public.item where public.item.track_number = $1;`, track_number)
	if err != nil {
		log.Println(err)
	}

	for rows.Next(){
        item := mymodel.Item{}
        err := rows.Scan(&item.Chrt_id, 
						&item.Track_number, 
						&item.Price,
						&item.Rid,
						&item.Name,
						&item.Sale,
						&item.Size,
						&item.Total_price,
						&item.Nm_id,
						&item.Brand,
						&item.Status,
		)
        if err != nil{
            log.Println(err)
            continue
        }
        model.Items = append(model.Items, item)
    }
	
	return model
}

// TODO: тут надо подумать приходят ли валидные данные или стоит сделать проверку
func SaveData(m mymodel.MyModel) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`insert into public.order values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, 
			m.Order_uid, m.Track_number, m.Entry, m.Locale, m.Internal_signature, m.Customer_id,
			m.Delivery_service, m.Shardkey, m.Sm_id, m.Date_created, m.Oof_shard,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into public.delivery values ($1, $2, $3, $4, $5, $6, $7, $8)`, 
			m.Delivery.Delivery_uuid, m.Delivery.Name, m.Delivery.Phone, m.Delivery.Zip,
			m.Delivery.City, m.Delivery.Address, m.Delivery.Region, m.Delivery.Email,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into public.payment values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, 
			m.Payment.Transaction, m.Payment.Request_id, m.Payment.Currency, m.Payment.Provider,
			m.Payment.Amount, m.Payment.Payment_dt, m.Payment.Bank, m.Payment.Delivery_cost,
			m.Payment.Goods_total, m.Payment.Custom_fee,
	)
	if err != nil {
		return err
	}

	for _, item := range m.Items {
		_, err = db.Exec(`insert into public.item values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, 
			item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.Total_price, item.Nm_id, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return nil
}