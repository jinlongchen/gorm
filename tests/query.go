package tests

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestFind(t *testing.T, db *gorm.DB) {
	db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})

	t.Run("Find", func(t *testing.T) {
		var users = []User{{
			Name:     "find",
			Age:      1,
			Birthday: Now(),
		}, {
			Name:     "find",
			Age:      2,
			Birthday: Now(),
		}, {
			Name:     "find",
			Age:      3,
			Birthday: Now(),
		}}

		if err := db.Create(&users).Error; err != nil {
			t.Fatalf("errors happened when create users: %v", err)
		}

		t.Run("First", func(t *testing.T) {
			var first User
			if err := db.Where("name = ?", "find").First(&first).Error; err != nil {
				t.Errorf("errors happened when query first: %v", err)
			} else {
				AssertObjEqual(t, first, users[0], "Name", "Age", "Birthday")
			}
		})

		t.Run("Last", func(t *testing.T) {
			var last User
			if err := db.Where("name = ?", "find").Last(&last).Error; err != nil {
				t.Errorf("errors happened when query last: %v", err)
			} else {
				AssertObjEqual(t, last, users[2], "Name", "Age", "Birthday")
			}
		})

		var all []User
		if err := db.Where("name = ?", "find").Find(&all).Error; err != nil || len(all) != 3 {
			t.Errorf("errors happened when query find: %v, length: %v", err, len(all))
		} else {
			for idx, user := range users {
				t.Run("FindAll#"+strconv.Itoa(idx+1), func(t *testing.T) {
					AssertObjEqual(t, all[idx], user, "Name", "Age", "Birthday")
				})
			}
		}

		t.Run("FirstMap", func(t *testing.T) {
			var first = map[string]interface{}{}
			if err := db.Model(&User{}).Where("name = ?", "find").First(first).Error; err != nil {
				t.Errorf("errors happened when query first: %v", err)
			} else {
				for _, name := range []string{"Name", "Age", "Birthday"} {
					t.Run(name, func(t *testing.T) {
						dbName := db.NamingStrategy.ColumnName("", name)
						reflectValue := reflect.Indirect(reflect.ValueOf(users[0]))
						AssertEqual(t, first[dbName], reflectValue.FieldByName(name).Interface())
					})
				}
			}
		})

		var allMap = []map[string]interface{}{}
		if err := db.Model(&User{}).Where("name = ?", "find").Find(&allMap).Error; err != nil {
			t.Errorf("errors happened when query first: %v", err)
		} else {
			for idx, user := range users {
				t.Run("FindAllMap#"+strconv.Itoa(idx+1), func(t *testing.T) {
					for _, name := range []string{"Name", "Age", "Birthday"} {
						t.Run(name, func(t *testing.T) {
							dbName := db.NamingStrategy.ColumnName("", name)
							reflectValue := reflect.Indirect(reflect.ValueOf(user))
							AssertEqual(t, allMap[idx][dbName], reflectValue.FieldByName(name).Interface())
						})
					}
				})
			}
		}
	})
}