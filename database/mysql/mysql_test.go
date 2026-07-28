package mysql

//import (
//	"context"
//	"testing"
//	"time"
//
//	"github.com/lsjhtang/kid/database/repos"
//)
//
//func NewMySQL() *MySQL {
//	return New(
//		WithHost("127.0.01"),
//		WithPort(3306),
//		WithUser("test"),
//		WithPassword("123456"),
//		WithDb("test"),
//		WithMaxLife(24),
//		WithMaxIdle(10),
//		WithMaxOpen(50),
//		WithLogLevel(3),
//	)
//}
//
//func NewRepository() repos.IBasicRepository {
//	return &Repository{DB: NewMySQL()}
//}
//
//type User struct {
//	Id        int64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:用户 ID"`
//	Uid       int64  `gorm:"column:uid;type:bigint;uniqueIndex:idx_uid;comment:用户 uid"`
//	Phone     string `gorm:"column:phone;type:varchar(18);uniqueIndex:idx_phone;comment:手机号"`
//	Gender    int8   `gorm:"column:gender;type:tinyint;comment:性别 1=男 2=女"`
//	Nickname  string `gorm:"column:nickname;type:varchar(24);comment:昵称"`
//	Avatar    string `gorm:"column:avatar;type:varchar(256);comment:头像"`
//	Password  string `gorm:"column:password;type:varchar(64);comment:密码"`
//	CreatedAt int64  `gorm:"column:created_at;type:bigint;comment:创建时间"`
//	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;comment:更新时间"`
//}
//
//func TestGet(t *testing.T) {
//	repo := NewRepository()
//	var user User
//	err := repo.GetOne(&repos.Query{
//		Recipient: &user,
//		Wheres:    repos.NewWheres().And(repos.Eq("uid", 9)),
//	})
//	if err != nil {
//		t.Error(err)
//	} else {
//		t.Log(user)
//	}
//}
//
//func TestFind(t *testing.T) {
//	repo := NewRepository()
//	var (
//		users []*User
//		total int64
//	)
//	if err := repo.Find(&repos.Query{
//		Wheres:    repos.NewWheres().And(repos.Lte("created_at", 1683363168392)),
//		Recipient: &users,
//		OrderBy:   "created_at ASC",
//		Num:       1,
//		Size:      10,
//		Total:     &total,
//	}); err != nil {
//		t.Error(err)
//	}
//	for _, user := range users {
//		t.Log(user)
//	}
//	t.Logf("总数: %d", total)
//}
//
//func TestCreate(t *testing.T) {
//	repo := NewRepository()
//	user := User{
//		Uid:       9,
//		Phone:     "12345",
//		Gender:    1,
//		Nickname:  "1",
//		Avatar:    "1",
//		Password:  "1",
//		CreatedAt: time.Now().UnixMilli(),
//		UpdatedAt: time.Now().UnixMilli(),
//	}
//	if err := repo.Create(&user); err != nil {
//		t.Error(err)
//	}
//	t.Log("ok")
//}
//
//func TestCreateInTx(t *testing.T) {
//	repo := NewRepository()
//	if err := repo.NewTx(context.Background(), func(tx context.Context) error {
//		if err := repo.CreateInTx(tx, &User{
//			Uid:       9,
//			Phone:     "12345",
//			Gender:    1,
//			Nickname:  "1",
//			Avatar:    "1",
//			Password:  "1",
//			CreatedAt: time.Now().UnixMilli(),
//			UpdatedAt: time.Now().UnixMilli(),
//		}); err != nil {
//			return err
//		}
//		if err := repo.UpdateColumnInTx(tx, &repos.Update{
//			Model:  new(User),
//			Wheres: repos.NewWheres().And(repos.Eq("uid", 9)),
//			Fields: map[string]interface{}{
//				"chip": 18,
//			},
//		}); err != nil {
//			return err
//		}
//		return nil
//	}); err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log("ok")
//}
//
//func TestUpdate(t *testing.T) {
//	repo := NewRepository()
//	user := User{
//		Phone:     "12345",
//		Gender:    2,
//		Nickname:  "2",
//		Avatar:    "2",
//		Password:  "2",
//		CreatedAt: 0,
//		UpdatedAt: 2,
//	}
//	if err := repo.Update(&repos.Update{
//		Wheres:  repos.NewWheres().And(repos.Eq("uid", 9)),
//		Update:  &user,
//		SaveNil: true,
//		Omits:   []string{"phone", "updated_at"},
//	}); err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log("ok")
//}
//
//func TestUpdateInTx(t *testing.T) {
//	repo := NewRepository()
//	var user = &User{
//		Phone: "12345",
//	}
//	if err := repo.NewTx(context.Background(), func(tx context.Context) error {
//		if err := repo.UpdateInTx(tx, &repos.Update{
//			Model:  new(User),
//			Wheres: repos.NewWheres().And(repos.Eq("uid", 9)),
//			Update: &user,
//		}); err != nil {
//			return err
//		}
//		return nil
//	}); err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log("ok")
//}
//
//func TestDelete(t *testing.T) {
//	repo := NewRepository()
//	if err := repo.Delete(&repos.Delete{
//		Model:  new(User),
//		Wheres: repos.NewWheres().And(repos.Eq("uid", 0)),
//	}); err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log("ok")
//}
//
//func TestDeleteInTx(t *testing.T) {
//	repo := NewRepository()
//	if err := repo.NewTx(context.Background(), func(tx context.Context) error {
//		if err := repo.DeleteInTx(tx, &repos.Delete{
//			Model:  new(User),
//			Wheres: repos.NewWheres().And(repos.Eq("uid", 0)),
//		}); err != nil {
//			return err
//		}
//		//if err := repo.UpdateColumnInTx(tx, &repos.Update{
//		//	Model:  new(User),
//		//	Wheres: repos.NewWheres().Add(repos.Eq("uid", 160263)),
//		//	Fields: map[string]interface{}{
//		//		"age": 18,
//		//	},
//		//}); err != nil {
//		//	return err
//		//}
//		return nil
//	}); err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log("ok")
//}
