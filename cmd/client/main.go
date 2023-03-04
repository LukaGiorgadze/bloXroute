package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/LukaGiorgadze/bloXroute/configs"
	"github.com/LukaGiorgadze/bloXroute/internal/client"
	"github.com/LukaGiorgadze/bloXroute/internal/models"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/nats-io/nats.go"
)

func main() {

	cfg, err := configs.NewConfig()
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	// Initialize message client with Messaging System connection
	// See detailed comment in /cmd/server/main.go
	var msgClient client.IMessageClient = client.NewNatsClient(cfg.NatsURL, []nats.Option{nats.UserInfo(cfg.NatsUser, cfg.NatsPass)})
	err = msgClient.Connect()
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	defer func() {
		err = msgClient.Disconnect()
		if err != nil {
			color.Error.Println(err)
			os.Exit(1)
		}
	}()

	app := gcli.NewApp()
	app.Desc = "Client-Server application User"
	var key string
	var val string
	var random int
	var stress int = 1

	app.Add(&gcli.Command{
		Name: "get",
		Desc: "<info>get</> retrieves the list. <info>get -k {key}</> get specific item. <info>get -k {key} -stress {n}</> send <info>{n}</> amount of req.",
		Func: func(cmd *gcli.Command, args []string) (err error) {

			var data []byte
			var subj client.Subject = client.ItemGetListSubject

			if key != "" {
				data, err = json.Marshal(models.Item{
					Key: key,
				})
				subj = client.ItemGetSubject
			}
			if err != nil {
				return
			}

			for i := 0; i < stress; i++ {
				err = msgClient.Publish(subj, data)
			}
			return

		},
		Config: func(c *gcli.Command) {
			c.StrOpt(&key, "k", "", "", "")
			c.IntOpt(&stress, "stress", "", stress, "")
		},
	})

	app.Add(&gcli.Command{
		Name: "add",
		Desc: "<info>add -k {key} -v {value}</> or use <info>add random {N}</> to add random N items",
		Func: func(cmd *gcli.Command, args []string) error {
			if key == "" || val == "" {
				return errors.New("key and value should not be empty.")
			}

			data, err := json.Marshal(models.Item{
				Key:   key,
				Value: val,
			})
			if err != nil {
				return err
			}

			if err := msgClient.Publish(client.ItemMutateAddSubject, data); err != nil {
				return err
			}
			return nil
		},
		Config: func(c *gcli.Command) {
			c.StrOpt(&key, "k", "", "", "")
			c.StrOpt(&val, "v", "", "", "")
		},
		Subs: []*gcli.Command{
			{
				Name: "random",
				Desc: "<info>add random {N}</> items",
				Func: func(cmd *gcli.Command, args []string) error {

					for i := 1; i <= random; i++ {
						data, err := json.Marshal(models.Item{
							Key:   fmt.Sprintf("key_%d", i),
							Value: fmt.Sprintf("Value %d", i),
						})
						if err != nil {
							return err
						}

						if err := msgClient.Publish(client.ItemMutateAddSubject, data); err != nil {
							return err
						}
					}

					return nil
				},
				Config: func(c *gcli.Command) {
					c.IntOpt(&random, "n", "", random, "")
				},
			},
		},
	})

	app.Add(&gcli.Command{
		Name: "delete",
		Desc: "<info>delete -k {key}</>",
		Func: func(cmd *gcli.Command, args []string) error {
			if key == "" {
				return errors.New("key should not be empty.")
			}

			data, err := json.Marshal(models.Item{
				Key: key,
			})
			if err != nil {
				return err
			}

			if err := msgClient.Publish(client.ItemMutateDeleteSubject, data); err != nil {
				return err
			}
			return nil
		},
		Config: func(c *gcli.Command) {
			c.StrOpt(&key, "k", "", "", "")
		},
	})

	app.Run(nil)

}
