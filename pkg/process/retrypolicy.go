package process

import (
	"encoding/json"
	"fmt"
)

type RetryPolicy struct {
	MethodConfig []MethodConfig `json:"methodConfig"`
}

type MethodConfig struct {
	Name         []Name         `json:"name"`
	WaitForReady bool           `json:"waitForReady"`
	RetryPolicy  RetryPolicyDef `json:"retryPolicy"`
}

type Name struct {
	Service string `json:"service"`
}

type RetryPolicyDef struct {
	MaxAttempts          int      `json:"MaxAttempts"`
	InitialBackoff       string   `json:"InitialBackoff"`
	MaxBackoff           string   `json:"MaxBackoff"`
	BackoffMultiplier    float64  `json:"BackoffMultiplier"`
	RetryableStatusCodes []string `json:"RetryableStatusCodes"`
}

func main() {
	policy := RetryPolicy{
		MethodConfig: []MethodConfig{
			{
				Name: []Name{
					{
						Service: "grpc.examples.echo.Echo",
					},
				},
				WaitForReady: true,
				RetryPolicy: RetryPolicyDef{
					MaxAttempts:          4,
					InitialBackoff:       ".01s",
					MaxBackoff:           ".01s",
					BackoffMultiplier:    1.0,
					RetryableStatusCodes: []string{"UNAVAILABLE"},
				},
			},
		},
	}

	jsonPolicy, err := json.MarshalIndent(policy, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%s\n", jsonPolicy)
}
