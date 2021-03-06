package gateway

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockListUsersClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockListUsersClient) ListUsers(
	*provider.ListUsersInput) (*provider.ListUsersOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}

	return &provider.ListUsersOutput{
		Users: []*provider.UserType{
			{
				Attributes: []*provider.AttributeType{
					{
						Name:  aws.String("given_name"),
						Value: aws.String("Joe"),
					},
					{
						Name:  aws.String("family_name"),
						Value: aws.String("Blow"),
					},
					{
						Name:  aws.String("email"),
						Value: aws.String("joe@blow.com"),
					},
					{
						Name:  aws.String("phone_number"),
						Value: aws.String("+1234567890"),
					},
				},
				Enabled:              aws.Bool(true),
				UserCreateDate:       &time.Time{},
				UserLastModifiedDate: &time.Time{},
				Username:             aws.String("user123"),
				UserStatus:           aws.String("CONFIRMED"),
			},
		},
		PaginationToken: aws.String("token123"),
	}, nil
}

func TestListUsers(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListUsersClient{}}
	result, err := gw.ListUsers(
		aws.Int64(10),
		aws.String("paginationToken"),
		aws.String("userPoolId"),
	)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestListUsersFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListUsersClient{serviceErr: true}}
	result, err := gw.ListUsers(
		aws.Int64(10),
		aws.String("paginationToken"),
		aws.String("userPoolId"),
	)
	assert.Nil(t, result)
	assert.Error(t, err)
}
