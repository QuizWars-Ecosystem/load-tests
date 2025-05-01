package samples

import (
	"context"

	"github.com/QuizWars-Ecosystem/load-tests/units/clients"
	"github.com/QuizWars-Ecosystem/load-tests/units/data"
	"github.com/QuizWars-Ecosystem/load-tests/units/scenario"
)

var client *clients.UsersService

func init() {
	scenario.Register("users_auth_register", &scenario.Scenario{
		Name: "register_user",
		Call: func(ctx context.Context) error {
			req := data.GenerateRegisterRequest()
			res, err := client.Register(ctx, req)
			if err != nil {
				return err
			}

			data.ProfilesData = append(data.ProfilesData, &data.ProfileData{
				Profile:  res.Profile,
				Token:    res.Token,
				Password: req.Password,
			})

			return nil
		},
		Init: func() error {
			var err error
			client, err = clients.GetUsersClient(scenario.GlobalAddr)
			if err != nil {
				return err
			}

			return nil
		},
	})

	scenario.Register("users_auth_login", &scenario.Scenario{
		Name: "login_user",
		Call: func(ctx context.Context) error {
			_, err := client.Login(ctx, data.GenerateLoginRequest())
			return err
		},
		Init: func() error {
			var err error
			client, err = clients.GetUsersClient(scenario.GlobalAddr)
			if err != nil {
				return err
			}

			return nil
		},
		DependsOn: []string{"users_auth_register"},
	})

	scenario.Register("users_profile_get_profile", &scenario.Scenario{
		Name: "get_profile",
		Call: func(ctx context.Context) error {
			requesterCtx, req := data.GenerateGetProfileRequest()
			_, err := client.GetProfile(requesterCtx, req)
			return err
		},
		Init: func() error {
			var err error
			client, err = clients.GetUsersClient(scenario.GlobalAddr)
			if err != nil {
				return err
			}

			return nil
		},
		DependsOn: []string{"users_auth_register"},
	})
}
