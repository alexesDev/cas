package cas

import (
	"encoding/json"
	"log"
)

// Step is
type step struct {
	Type  string
	Input json.RawMessage
}

type stepAnswer struct {
	Type   string
	Error  string
	Output json.RawMessage
}

type task struct {
	Addr string
	Plan []step
}

type uploadPLUInput struct {
	ScaleID   uint32
	PLUNumber uint32
}

type downloadPLUInput struct {
	ScaleID uint32
	Data    PLUData
}

type erasePLUInput struct {
	ScaleID          uint32
	DepartmentNumber uint16
	PLUNumber        uint32
}

type taskAnswer struct {
	Plan []stepAnswer
}

// ProcessJSON starts execution of the steps described in JSON
func ProcessJSON(buf []byte) ([]byte, error) {
	var task task
	var answer taskAnswer

	if err := json.Unmarshal(buf, &task); err != nil {
		return nil, err
	}

	scale, err := Connect(task.Addr)
	defer scale.Disconnect()

	if err != nil {
		return nil, err
	}

	for _, step := range task.Plan {
		var as stepAnswer
		as.Type = step.Type

		switch step.Type {
		case "DownloadPLU":
			var input downloadPLUInput

			if err := json.Unmarshal(step.Input, &input); err != nil {
				as.Error = err.Error()
				goto End
			}

			err = scale.DownloadPLU(input.ScaleID, input.Data)

			if err != nil {
				as.Error = err.Error()
				goto End
			}
		case "UploadPLU":
			var input uploadPLUInput

			if err := json.Unmarshal(step.Input, &input); err != nil {
				as.Error = err.Error()
				goto End
			}

			data, err := scale.UploadPLU(input.ScaleID, input.PLUNumber)

			if err != nil {
				as.Error = err.Error()
				goto End
			}

			b, err := json.Marshal(data)

			if err != nil {
				as.Error = err.Error()
				goto End
			}

			as.Output = b
		case "ErasePLU":
			var input erasePLUInput

			if err := json.Unmarshal(step.Input, &input); err != nil {
				as.Error = err.Error()
				goto End
			}

			err = scale.ErasePLU(input.ScaleID, input.DepartmentNumber, input.PLUNumber)

			if err != nil {
				as.Error = err.Error()
			}
		default:
			log.Fatalf("unknown step type: %q", step.Type)
		}

	End:
		answer.Plan = append(answer.Plan, as)
	}

	b, err := json.Marshal(answer)

	if err != nil {
		return nil, err
	}

	return b, nil
}
