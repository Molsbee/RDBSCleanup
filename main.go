package main

import (
	"fmt"
	"github.com/Molsbee/RDBSCleanup/clc"
	"github.com/Molsbee/RDBSCleanup/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

var (
	databaseUsername string
	databasePassword string
	clcUsername      string
	clcPassword      string
)

func init() {
	if databaseUsername = os.Getenv("DB_User"); len(databaseUsername) == 0 {
		log.Panic("environment variable 'DB_User' needs to be set")
	}
	if databasePassword = os.Getenv("DB_Pass"); len(databasePassword) == 0 {
		log.Panic("environment variable 'DB_Pass' needs to be set")
	}
	if clcUsername = os.Getenv("CLC_Username"); len(clcUsername) == 0 {
		log.Panic("environment variable 'CLC_Username' needs to be set")
	}
	if clcPassword = os.Getenv("CLC_Password"); len(clcPassword) == 0 {
		log.Panic("environment variable 'CLC_Password' needs to be set")
	}
}

func main() {
	databaseConnectionString := fmt.Sprintf("%s:%s@tcp(10.90.85.95)/dbaas?charset=utf8mb4&parseTime=True&loc=Local", databaseUsername, databasePassword)
	db, err := gorm.Open(mysql.Open(databaseConnectionString), &gorm.Config{})
	if err != nil {
		log.Panicf("unable to connect to database %s", err)
	}

	var activeSubscriptions []model.Subscription
	err = db.Preload("Customer").Find(&activeSubscriptions, "subscription_status = ?", "Active").Error
	if err != nil || len(activeSubscriptions) == 0 {
		log.Panicf("no activeSubscriptions found or error occurred %s", err)
	}

	var serviceAccount model.ServiceAccount
	if err = db.Find(&serviceAccount, "username = ?", "ctl_mysql").Error; err != nil {
		log.Panicf("unable to read service account for deleting appfog subscriptions")
	}

	api, err := clc.NewAPI(clc.Config{
		CLCUsername:    os.Getenv("CLC_Username"),
		CLCPassword:    os.Getenv("CLC_Password"),
		RDBSAppfogUser: serviceAccount.Username,
		RDBSAppfogPass: serviceAccount.Password,
	})
	if err != nil {
		log.Panicf("unable to create api service %s", err)
	}

	for _, activeSubscription := range activeSubscriptions {
		accountAlias := activeSubscription.Customer.Alias
		accountDetails, _ := api.GetAccount(accountAlias)
		if accountDetails["status"] == "Deleted" {
			fmt.Printf("attempting to delete account (%s) subscription (%d)\n", accountAlias, activeSubscription.ID)
			if strings.Contains(activeSubscription.InstanceType, "APPFOG") {
				if err := api.DeleteAppfogSubscriptions(activeSubscription.ExternalID); err != nil {
					fmt.Printf("error occurred delete rdbs subscription id %d for account alias %s (%s)\n", activeSubscription.ID, accountAlias, err)
				}
			} else {
				if err := api.DeleteRDBSSubscription(accountAlias, activeSubscription.ID); err != nil {
					fmt.Printf("error occurred delete rdbs subscription id %d for account alias %s (%s)\n", activeSubscription.ID, accountAlias, err)
				}
			}
		}
	}
}
