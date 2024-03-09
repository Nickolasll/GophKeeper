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
		Name:      "register",
		Usage:     "register new user via username and password",
		ArgsUsage: "[username] [password]",
		Aliases:   []string{"r"},
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
		Name:      "login",
		Usage:     "sign in via username and password",
		ArgsUsage: "[username] [password]",
		Aliases:   []string{"l"},
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

			err = app.SyncAll.Do(session)
			if err != nil {
				log.Error(err)

				return cli.Exit(err, 1)
			}
			fmt.Println("all data was syncronized successfully")

			return nil
		},
	}
}

func createText() cli.Command {
	return cli.Command{
		Name:      "text",
		Usage:     "create new text content",
		ArgsUsage: "[content]",
		Aliases:   []string{"t"},
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

			fmt.Println("text created successfully")

			return nil
		},
	}
}

func updateText() cli.Command {
	return cli.Command{
		Name:      "text",
		Usage:     "update existing text via id and new content string",
		ArgsUsage: "[id] [content]",
		Aliases:   []string{"t"},
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
					fmt.Println("text not found, id: ", textID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("text updated successfully")

			return nil
		},
	}
}

func showText() cli.Command {
	return cli.Command{
		Name:    "texts",
		Usage:   "shows current user text data",
		Aliases: []string{"t"},
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
		Name:    "texts",
		Usage:   "override current user text data from remote",
		Aliases: []string{"t"},
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
			fmt.Println("text syncronized successfully")

			return nil
		},
	}
}

func createBinary() cli.Command {
	return cli.Command{
		Name:      "binary",
		Usage:     "create new binary content",
		ArgsUsage: "[path-to-file]",
		Aliases:   []string{"b"},
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

			fmt.Println("binary created successfully")

			return nil
		},
	}
}

func updateBinary() cli.Command {
	return cli.Command{
		Name:      "binary",
		Usage:     "update existing binary via id and data",
		ArgsUsage: "[id] [path-to-file]",
		Aliases:   []string{"b"},
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
					fmt.Println("binary not found, id: ", binID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("binary updated successfully")

			return nil
		},
	}
}

func showBinary() cli.Command {
	return cli.Command{
		Name:    "binaries",
		Usage:   "shows current user binary data",
		Aliases: []string{"b"},
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
		Name:    "binaries",
		Usage:   "override current user binary data from remote",
		Aliases: []string{"b"},
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
			fmt.Println("binary syncronized successfully")

			return nil
		},
	}
}

func createCredentials() cli.Command {
	return cli.Command{
		Name:      "credentials",
		Usage:     "create new credentials",
		ArgsUsage: "[name] [login] [password]",
		Aliases:   []string{"c"},
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

			fmt.Println("credentials created successfully")

			return nil
		},
	}
}

func updateCredentials() cli.Command {
	var name, login, password string

	return cli.Command{
		Name:      "credentials",
		Usage:     "update existing credentials via id and flags",
		ArgsUsage: "[id]",
		Aliases:   []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "name value to update",
				DefaultText: "",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "login",
				Aliases:     []string{"l"},
				Usage:       "login value to update",
				DefaultText: "",
				Destination: &login,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "password value to update",
				DefaultText: "",
				Destination: &password,
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

			err = app.UpdateCredentials.Do(*currentSession, credID, name, login, password)
			if err != nil {
				if errors.Is(err, domain.ErrEntityNotFound) {
					fmt.Println("credentials not found, id: ", credID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("credentials updated successfully")

			return nil
		},
	}
}

func showCredentials() cli.Command {
	return cli.Command{
		Name:    "credentials",
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
		Name:    "credentials",
		Usage:   "override current user credentials from remote",
		Aliases: []string{"c"},
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
			fmt.Println("credentials syncronized successfully")

			return nil
		},
	}
}

func createBankCard() cli.Command {
	return cli.Command{
		Name:      "bank-card",
		Usage:     "create new bank-card",
		ArgsUsage: "[number] [valid-thru] [cvv] [(optional) card-holder]",
		Aliases:   []string{"bc"},
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
				fmt.Println("invalid card-number: ", number)

				return nil
			}
			if !validValidThru.MatchString(validThru) {
				fmt.Println("invalid valid-thru value: ", validThru)

				return nil
			}
			if !validCVV.MatchString(cvv) {
				fmt.Println("invalid cvv value: ", cvv)

				return nil
			}
			if !validCardHolder.MatchString(cardHolder) {
				fmt.Println("invalid card-holder value: ", cardHolder)

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

			fmt.Println("bank-card created successfully")

			return nil
		},
	}
}

func updateBankCard() cli.Command {
	var number, validThru, cvv, cardHolder string

	return cli.Command{
		Name:      "bank-card",
		Usage:     "update existing bank-card via id and flags",
		ArgsUsage: "[id]",
		Aliases:   []string{"bc"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "number",
				Usage:       "number value to update",
				DefaultText: "",
				Destination: &number,
			},
			&cli.StringFlag{
				Name:        "valid-thru",
				Usage:       "valid-thru value to update",
				DefaultText: "",
				Destination: &validThru,
			},
			&cli.StringFlag{
				Name:        "cvv",
				Usage:       "cvv value to update",
				DefaultText: "",
				Destination: &cvv,
			},
			&cli.StringFlag{
				Name:        "card-holder",
				Usage:       "(optional) card-holder value to update",
				DefaultText: "",
				Destination: &cardHolder,
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
				fmt.Println(err, "invalid bank-card id: ", id)

				return nil
			}

			if cmd.NumFlags() == 0 {
				fmt.Println("invalid input: please pass at least one attribute (number, valid-thru, cvv, card-holder) to update")

				return nil
			}

			if number != "" && !validCardNumber.MatchString(number) {
				fmt.Println("invalid card-number: ", number)

				return nil
			}
			if validThru != "" && !validValidThru.MatchString(validThru) {
				fmt.Println("invalid valid-thru value: ", validThru)

				return nil
			}
			if cvv != "" && !validCVV.MatchString(cvv) {
				fmt.Println("invalid cvv value: ", cvv)

				return nil
			}
			if cardHolder != "" && !validCardHolder.MatchString(cardHolder) {
				fmt.Println("invalid card-holder value: ", cardHolder)

				return nil
			}

			err = app.UpdateBankCard.Do(*currentSession, cardID, number, validThru, cvv, cardHolder)
			if err != nil {
				if errors.Is(err, domain.ErrEntityNotFound) {
					fmt.Println("bank-card not found, id: ", cardID)

					return nil
				} else if errors.Is(err, domain.ErrBadRequest) {
					fmt.Println(err)

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("bank-card updated successfully")

			return nil
		},
	}
}

func showBankCards() cli.Command {
	return cli.Command{
		Name:    "bank-cards",
		Usage:   "shows current user bank-cards",
		Aliases: []string{"bc"},
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
		Name:    "bank-cards",
		Usage:   "override current user bank-cards from remote",
		Aliases: []string{"bc"},
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
			fmt.Println("bank-cards syncronized successfully")

			return nil
		},
	}
}

func syncAll() cli.Command {
	return cli.Command{
		Name:    "all",
		Usage:   "override all user data from remote",
		Aliases: []string{"a"},
		Action: func(_ context.Context, _ *cli.Command) error {
			if currentSession == nil {
				fmt.Println("unauthorized")

				return nil
			}

			err := app.SyncAll.Do(*currentSession)

			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) {
					fmt.Println("unauthorized")

					return nil
				} else {
					log.Error(err)

					return cli.Exit(err, 1)
				}
			}
			fmt.Println("all data was syncronized successfully")

			return nil
		},
	}
}
