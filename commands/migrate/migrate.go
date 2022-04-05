// Package migrate Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package migrate

import (
	"container/list"
	"fmt"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"studyhut/models"
)

var (
	migrationList = &migrationCache{}
)

type MigrationDatabase interface {
	Version() int64                         //获取当前的版本
	ValidUpdate(version int64) error        //校验当前是否可更新
	ValidForBackupTableSchema() error       //校验并备份表结构
	ValidForUpdateTableSchema() error       //校验并更新表结构
	MigrationOldTableData() error           //恢复旧数据
	MigrationNewTableData() error           //插入新数据
	AddMigrationRecord(version int64) error //增加迁移记录
	MigrationCleanup() error                //最后的清理工作
	RollbackMigration() error               //回滚本次迁移
}

type migrationCache struct {
	items *list.List
}

func RunMigration() {

	if len(os.Args) >= 2 && os.Args[1] == "migrate" {

		migrate, err := models.NewMigration().FindFirst()

		if err != nil {
			//log.Fatalf("migrations table %s", err)
			migrate = models.NewMigration()
		}
		fmt.Println("Start migration databae... ")

		for el := migrationList.items.Front(); el != nil; el = el.Next() {

			//如果存在比当前版本大的版本，则依次升级
			if item, ok := el.Value.(MigrationDatabase); ok && item.Version() > migrate.Version {
				err := item.ValidUpdate(migrate.Version)
				if err != nil {
					log.Fatal(err)
				}
				err = item.ValidForBackupTableSchema()
				if err != nil {
					item.RollbackMigration()
					log.Fatal(err)
				}
				err = item.ValidForUpdateTableSchema()
				if err != nil {
					item.RollbackMigration()
					log.Fatal(err)
				}
				err = item.MigrationOldTableData()
				if err != nil {
					item.RollbackMigration()
					log.Fatal(err)
				}
				err = item.MigrationNewTableData()
				if err != nil {
					item.RollbackMigration()
					log.Fatal(err)
				}
				err = item.AddMigrationRecord(item.Version())
				if err != nil {
					item.RollbackMigration()
					log.Fatal(err)
				}
				err = item.MigrationCleanup()
				if err != nil {
					item.RollbackMigration()
					log.Fatal(err)
				}
			}
		}
		fmt.Println("Migration successful.")
		os.Exit(0)
	}
}

// ExportDatabaseTable 导出数据库的表结构
func ExportDatabaseTable() ([]string, error) {
	dbAdapter := beego.AppConfig.String("db_adapter")
	dbDatabase := beego.AppConfig.String("db_database")
	tables := make([]string, 0)

	o := orm.NewOrm()
	switch dbAdapter {
	case "mysql":
		{
			var lists []orm.Params
			_, err := o.Raw(fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'", dbDatabase)).Values(&lists)
			if err != nil {
				return tables, err
			}
			for _, table := range lists {
				var results []orm.Params

				_, err = o.Raw(fmt.Sprintf("show create table %s", table["TABLE_NAME"])).Values(&results)
				if err != nil {
					return tables, err
				}
				tables = append(tables, results[0]["Create Table"].(string))
			}
			break
		}
	case "sqlite3":
		{
			var results []orm.Params
			_, err := o.Raw("SELECT sql FROM sqlite_master WHERE sql IS NOT NULL ORDER BY rootpage ASC").Values(&results)
			if err != nil {
				return tables, err
			}
			for _, item := range results {
				if sql, ok := item["sql"]; ok {
					tables = append(tables, sql.(string))
				}
			}
			break
		}

	}
	return tables, nil
}

func RegisterMigration() {
	migrationList.items = list.New()

	migrationList.items.PushBack(NewMigrationVersion03())
}
