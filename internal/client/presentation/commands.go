// Package presentation содержит фабрику cli приложения и имплементацию команд
package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func registration() cli.Command {
	return cli.Command{
		Name:    "register",
		Usage:   "register new user via username and password",
		Aliases: []string{"r"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			login := cmd.Args().Get(0)
			password := cmd.Args().Get(1)
			session, err := app.Registration.Do(login, password)
			if err != nil {
				if errors.Is(err, domain.ErrLoginConflict) {
					fmt.Println("user with this login already exists: ", login)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			currentSession = &session
			fmt.Println("registration successful")

			return nil
		},
	}
}

func login() cli.Command {
	return cli.Command{
		Name:    "login",
		Usage:   "sign in via username and password",
		Aliases: []string{"l"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			login := cmd.Args().Get(0)
			password := cmd.Args().Get(1)
			session, err := app.Login.Do(login, password)
			if err != nil {
				if errors.Is(err, domain.ErrUnauthorized) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			currentSession = &session
			fmt.Println("login successful")

			return nil
		},
	}
}

func createText() cli.Command {
	return cli.Command{
		Name:    "create",
		Usage:   "create new text content",
		Aliases: []string{"c"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			content := cmd.Args().First()
			err := app.CreateText.Do(*currentSession, content)
			if err != nil {
				if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}

			fmt.Println("Text created successfully")

			return nil
		},
	}
}

func updateText() cli.Command {
	return cli.Command{
		Name:    "update",
		Usage:   "update existing text via id and new content string",
		Aliases: []string{"c"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			textID := cmd.Args().Get(0)
			content := cmd.Args().Get(1)
			err := app.UpdateText.Do(*currentSession, textID, content)
			if err != nil {
				if errors.Is(err, domain.ErrEntityNotFound) {
					fmt.Println("Text not found, id: ", textID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Text updated successfully")

			return nil
		},
	}
}

func showText() cli.Command { //nolint: dupl
	return cli.Command{
		Name:    "show",
		Usage:   "shows current user text data",
		Aliases: []string{"c"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			text, err := app.ShowText.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}

			s, err := json.MarshalIndent(text, "", "\t")
			if err != nil {
				log.Error(err)

				return cli.Exit(err, 1)
			}
			fmt.Print(string(s))

			return nil
		},
	}
}

func createBinary() cli.Command {
	return cli.Command{
		Name:    "create",
		Usage:   "create new binary content",
		Aliases: []string{"c"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			contentPath := cmd.Args().First()
			content, err := os.ReadFile(contentPath) //nolint: gosec
			if err != nil {
				fmt.Println(err)

				return nil
			}

			err = app.CreateBinary.Do(*currentSession, content)
			if err != nil {
				if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}

			fmt.Println("Binary created successfully")

			return nil
		},
	}
}

func updateBinary() cli.Command {
	return cli.Command{
		Name:    "update",
		Usage:   "update existing binary via id and data",
		Aliases: []string{"c"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			binID := cmd.Args().Get(0)
			contentPath := cmd.Args().Get(1)
			content, err := os.ReadFile(contentPath) //nolint: gosec
			if err != nil {
				fmt.Println(err)

				return nil
			}

			err = app.UpdateBinary.Do(*currentSession, binID, content)
			if err != nil {
				if errors.Is(err, domain.ErrEntityNotFound) {
					fmt.Println("Binary not found, id: ", binID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Binary updated successfully")

			return nil
		},
	}
}

func showBinary() cli.Command { //nolint: dupl
	return cli.Command{
		Name:    "show",
		Usage:   "shows current user binary data",
		Aliases: []string{"c"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			text, err := app.ShowBinary.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}

			s, err := json.MarshalIndent(text, "", "\t")
			if err != nil {
				log.Error(err)

				return cli.Exit(err, 1)
			}
			fmt.Print(string(s))

			return nil
		},
	}
}
