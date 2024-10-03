package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type User struct {
	Id   int    `mapstructure:"Id" yaml:"Id,flow"`
	Name string `mapstructure:"Name" yaml:"Name,flow"`
}

var (
	user  = newUserValue(0, "")
	users = newUserSliceValue([]User{})
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("user(pflag): %+v\n", user)   // debug
		cmd.Printf("users(pflag): %+v\n", users) // debug

		viper.SetConfigFile("config.yaml")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if cmd.Flags().Changed("user") {
			viper.Set("user", user)
		} else {
			if err := viper.UnmarshalKey("user", &user); err != nil {
				panic(err)
			}
		}

		if cmd.Flags().Changed("users") {
			tmp := []User{}
			if err := viper.UnmarshalKey("users", &tmp); err != nil {
				panic(err)
			}
			cmd.Printf("users(file): %+v\n", tmp) // debug

			*users.value = Merge(*users.value, tmp)

			viper.Set("users", *users.value)
		} else {
			if err := viper.UnmarshalKey("users", &users.value); err != nil {
				panic(err)
			}
		}

		cmd.Printf("user: %+v\n", user)   // debug
		cmd.Printf("users: %+v\n", users) // debug

		if cmd.Flags().Changed("user") || cmd.Flags().Changed("users") {
			if err := viper.WriteConfig(); err != nil {
				panic(err)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().Var(user, "user", "{Id: 1, Name: Foo}")
	rootCmd.Flags().Var(users, "users", `[{Id: 1, Name: Foo}, {Id: 2, Name: Bar}]`)

	// don't bind it because result of String() function is merged as string.
	//viper.BindPFlag("user", rootCmd.Flags().Lookup("user"))   // e.g.) user: '{Id: 1, Name: Foo}'
	//viper.BindPFlag("users", rootCmd.Flags().Lookup("users")) // e.g.) users: '[{Id: 1, Name: Foo}, {Id: 2, Name: Bar}]'
}
