package data

import (
	"context"

	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	usersv1 "github.com/QuizWars-Ecosystem/load-tests/gen/external/users/v1"
	"github.com/brianvoe/gofakeit/v7"
)

type ProfileData struct {
	Profile  *usersv1.Profile
	Token    string
	Password string
}

var ProfilesData []*ProfileData

var jwtService = &jwt.Service{}

func GenerateRegisterRequest() *usersv1.RegisterRequest {
	return &usersv1.RegisterRequest{
		AvatarId: int32(gofakeit.IntN(10)),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 16),
	}
}

func GenerateLoginRequest() *usersv1.LoginRequest {
	if ProfilesData == nil {
		byUsername := gofakeit.Bool()
		pass := gofakeit.Password(true, true, true, true, false, 16)

		switch byUsername {
		case true:
			return &usersv1.LoginRequest{
				Identifier: &usersv1.LoginRequest_Username{
					Username: gofakeit.Username(),
				},
				Password: pass,
			}
		default:
			return &usersv1.LoginRequest{
				Identifier: &usersv1.LoginRequest_Email{
					Email: gofakeit.Email(),
				},
				Password: pass,
			}
		}
	} else {
		p := ProfilesData[gofakeit.IntN(len(ProfilesData))]
		byUsername := gofakeit.Bool()

		switch byUsername {
		case true:
			return &usersv1.LoginRequest{
				Identifier: &usersv1.LoginRequest_Username{
					Username: p.Profile.Username,
				},
				Password: p.Password,
			}
		default:
			return &usersv1.LoginRequest{
				Identifier: &usersv1.LoginRequest_Email{
					Email: p.Profile.Email,
				},
				Password: p.Password,
			}
		}
	}
}

func GenerateGetProfileRequest() (context.Context, *usersv1.GetProfileRequest) {
	var requestCtx context.Context

	if ProfilesData == nil {
		requestCtx = context.Background()

		return requestCtx, &usersv1.GetProfileRequest{
			Identifier: &usersv1.GetProfileRequest_Username{
				Username: gofakeit.Username(),
			},
		}
	} else {
		initiator := ProfilesData[gofakeit.IntN(len(ProfilesData))]
		self := gofakeit.Bool()

		requestCtx = jwtService.SetTokenInContext(context.Background(), initiator.Token)

		if self {
			return requestCtx, &usersv1.GetProfileRequest{
				Identifier: &usersv1.GetProfileRequest_UserId{
					UserId: initiator.Profile.Id,
				},
			}
		} else {
			byUsername := gofakeit.Bool()
			target := ProfilesData[gofakeit.IntN(len(ProfilesData))].Profile

			if byUsername {
				return requestCtx, &usersv1.GetProfileRequest{
					Identifier: &usersv1.GetProfileRequest_Username{
						Username: target.Username,
					},
				}
			} else {
				return requestCtx, &usersv1.GetProfileRequest{
					Identifier: &usersv1.GetProfileRequest_UserId{
						UserId: target.Id,
					},
				}
			}
		}
	}
}
