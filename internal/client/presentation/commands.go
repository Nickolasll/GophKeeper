// Package presentation содержит фабрику cli приложения и имплементацию команд
package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func parseID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

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

			id := cmd.Args().Get(0)
			content := cmd.Args().Get(1)

			textID, err := parseID(id)
			if err != nil {
				fmt.Println(err, "invalid text id: ", id)

				return nil
			}

			err = app.UpdateText.Do(*currentSession, textID, content)
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

func showText() cli.Command {
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

func syncText() cli.Command {
	return cli.Command{
		Name:    "sync",
		Usage:   "override current user text data from remote",
		Aliases: []string{"s"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			err := app.SyncText.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Text syncronized successfully")

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

			id := cmd.Args().Get(0)
			contentPath := cmd.Args().Get(1)

			binID, err := parseID(id)
			if err != nil {
				fmt.Println(err, "invalid binary id: ", id)

				return nil
			}

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

func showBinary() cli.Command {
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

func syncBinary() cli.Command {
	return cli.Command{
		Name:    "sync",
		Usage:   "override current user binary data from remote",
		Aliases: []string{"s"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			err := app.SyncBinary.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Binary syncronized successfully")

			return nil
		},
	}
}

func createCredentials() cli.Command {
	return cli.Command{
		Name:    "create",
		Usage:   "create new credentials",
		Aliases: []string{"c"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			name := cmd.Args().Get(0)
			login := cmd.Args().Get(1)
			password := cmd.Args().Get(2)

			err := app.CreateCredentials.Do(*currentSession, name, login, password)
			if err != nil {
				if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}

			fmt.Println("Credentials created successfully")

			return nil
		},
	}
}

func updateCredentials() cli.Command {
	return cli.Command{
		Name:    "update",
		Usage:   "update existing credentials via id and flags",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "name to update",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:        "login",
				Aliases:     []string{"l"},
				Usage:       "login to update",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "password to update",
				DefaultText: "",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			id := cmd.Args().First()
			credID, err := parseID(id)
			if err != nil {
				fmt.Println(err, "invalid credentials id: ", id)

				return nil
			}

			if cmd.NumFlags() == 0 {
				fmt.Println("invalid input: please pass at least one attribute (name, login, password) to update")

				return nil
			}

			name := cmd.String("name")
			login := cmd.String("login")
			password := cmd.String("password")

			err = app.UpdateCredentials.Do(*currentSession, credID, name, login, password)
			if err != nil {
				if errors.Is(err, domain.ErrEntityNotFound) {
					fmt.Println("Credentials not found, id: ", credID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Credentials updated successfully")

			return nil
		},
	}
}

func showCredentials() cli.Command {
	return cli.Command{
		Name:    "show",
		Usage:   "shows current user credentials data",
		Aliases: []string{"c"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			text, err := app.ShowCredentials.Do(*currentSession)

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

func syncCredentials() cli.Command {
	return cli.Command{
		Name:    "sync",
		Usage:   "override current user credentials from remote",
		Aliases: []string{"s"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			err := app.SyncCredentials.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Credentials syncronized successfully")

			return nil
		},
	}
}

func createBankCard() cli.Command {
	return cli.Command{
		Name:    "create",
		Usage:   "create new bank card",
		Aliases: []string{"b"},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			number := cmd.Args().Get(0)
			validThru := cmd.Args().Get(1)
			cvv := cmd.Args().Get(2)
			cardHolder := ""
			if cmd.Args().Len() == 4 { //nolint: gomnd
				cardHolder = cmd.Args().Get(3)
			}

			if !validCardNumber.MatchString(number) {
				fmt.Println("invalid card number: ", number)

				return nil
			}
			if !validValidThru.MatchString(validThru) {
				fmt.Println("invalid valid thru value: ", validThru)

				return nil
			}
			if !validCVV.MatchString(cvv) {
				fmt.Println("invalid cvv value: ", cvv)

				return nil
			}
			if !validCardHolder.MatchString(cardHolder) {
				fmt.Println("invalid card holder value: ", cardHolder)

				return nil
			}

			err := app.CreateBankCard.Do(*currentSession, number, validThru, cvv, cardHolder)
			if err != nil {
				if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}

			fmt.Println("Bank card created successfully")

			return nil
		},
	}
}

func updateBankCard() cli.Command {
	return cli.Command{
		Name:    "update",
		Usage:   "update existing bank card via id and flags",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "number",
				Aliases:     []string{"n"},
				Usage:       "number to update",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:        "valid_thru",
				Aliases:     []string{"vt"},
				Usage:       "valid thru value to update",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:        "cvv",
				Aliases:     []string{"c"},
				Usage:       "cvv value to update",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:        "card_holder",
				Aliases:     []string{"ch"},
				Usage:       "card holder value to update",
				DefaultText: "",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}
			id := cmd.Args().First()

			cardID, err := parseID(id)
			if err != nil {
				fmt.Println(err, "invalid bank card id: ", id)

				return nil
			}

			if cmd.NumFlags() == 0 {
				fmt.Println("invalid input: please pass at least one attribute (number, valid_thru, cvv, card_holder) to update")

				return nil
			}

			number := cmd.String("number")
			validThru := cmd.String("valid_thru")
			cvv := cmd.String("cvv")
			cardHolder := cmd.String("card_holder")

			if number != "" && !validCardNumber.MatchString(number) {
				fmt.Println("invalid card number: ", number)

				return nil
			}
			if validThru != "" && !validValidThru.MatchString(validThru) {
				fmt.Println("invalid valid thru value: ", validThru)

				return nil
			}
			if cvv != "" && !validCVV.MatchString(cvv) {
				fmt.Println("invalid cvv value: ", cvv)

				return nil
			}
			if cardHolder != "" && !validCardHolder.MatchString(cardHolder) {
				fmt.Println("invalid card holder value: ", cardHolder)

				return nil
			}

			err = app.UpdateBankCard.Do(*currentSession, cardID, number, validThru, cvv, cardHolder)
			if err != nil {
				if errors.Is(err, domain.ErrEntityNotFound) {
					fmt.Println("Bank card not found, id: ", cardID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Credentials updated successfully")

			return nil
		},
	}
}

func showBankCards() cli.Command {
	return cli.Command{
		Name:    "show",
		Usage:   "shows current user bank cards",
		Aliases: []string{"c"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			text, err := app.ShowBankCards.Do(*currentSession)

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

func syncBankCards() cli.Command {
	return cli.Command{
		Name:    "sync",
		Usage:   "override current user bank cards from remote",
		Aliases: []string{"s"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			err := app.SyncBankCards.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("Bank cards syncronized successfully")

			return nil
		},
	}
}
