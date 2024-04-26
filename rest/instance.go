// SPDX-License-Identifier: MIT

package rest

import (
	"github.com/eludris-community/eludris-api-types.go/v2/models"
)

func (r *restImpl) GetInstanceInfo() (info models.InstanceInfo, err error) {
	_, err = r.Request(InstanceInfo.Compile(nil), nil, &info)
	return
}
