package main

import (
	"fmt"

	"github.com/alex-savin/go-meater"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./")     // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(err)
	}

	client, _ := meater.New(
		meater.Username(viper.GetString("credentials.username")),
		meater.Password(viper.GetString("credentials.password")),
		meater.LogLevel(viper.GetString("log")),
	)

	if len(client.GetProbes()) > 0 {
		for i, probe := range client.GetProbes() {
			fmt.Printf("PROBE #%d: %+v\n", i+1, probe)
			fmt.Printf("COOK  #%d: %+v\n", i+1, probe.Cook)
		}
	} else {
		fmt.Print("No active probes are detected\n")
	}
}
