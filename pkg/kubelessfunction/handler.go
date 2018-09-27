package kubeless

import (
	"github.com/alexesDev/cas/pkg/cas"
	"github.com/kubeless/kubeless/pkg/functions"
)

func Handler(event functions.Event, context functions.Context) (string, error) {
	output, err := cas.ProcessJSON([]byte(event.Data))

	if err != nil {
		return "", err
	}

	return string(output), nil
}
