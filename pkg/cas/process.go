package cas

import (
	"encoding/json"
	"log"
)

type Step struct {
	Type  string
	Input json.RawMessage
}

type StepAnswer struct {
	Type   string
	Error  string
	Output json.RawMessage
}

type Task struct {
	Addr string
	Plan []Step
}

type UploadPLUInput struct {
	ScaleId   uint32
	PLUNumber uint32
}

type DownloadPLUInput struct {
	ScaleId uint32
	Data    PLUData
}

type ErasePLUInput struct {
	ScaleId          uint32
	DepartmentNumber uint16
	PLUNumber        uint32
}

type TaskAnswer struct {
	Plan []StepAnswer
}

func ProcessJSON(buf []byte) ([]byte, error) {
	var task Task
	var answer TaskAnswer

	if err := json.Unmarshal(buf, &task); err != nil {
		return nil, err
	}

	scale, err := Connect(task.Addr)

	if err != nil {
		return nil, err
	}

	for _, step := range task.Plan {
		var as StepAnswer
		as.Type = step.Type

		switch step.Type {
		case "DownloadPLU":
			var input DownloadPLUInput

			if err := json.Unmarshal(step.Input, &input); err != nil {
				as.Error = err.Error()
				goto End
			}

			err = scale.DownloadPLU(input.ScaleId, input.Data)

			if err != nil {
				as.Error = err.Error()
				goto End
			}
		case "UploadPLU":
			var input UploadPLUInput

			if err := json.Unmarshal(step.Input, &input); err != nil {
				as.Error = err.Error()
				goto End
			}

			data := scale.UploadPLU(input.ScaleId, input.PLUNumber)
			b, err := json.Marshal(data)

			if err != nil {
				as.Error = err.Error()
				goto End
			}

			as.Output = b
		case "ErasePLU":
			var input ErasePLUInput

			if err := json.Unmarshal(step.Input, &input); err != nil {
				as.Error = err.Error()
				goto End
			}

			err = scale.ErasePLU(input.ScaleId, input.DepartmentNumber, input.PLUNumber)

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
