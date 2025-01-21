package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

// Prompt.
type Prompt struct {
	Instruction string
	PullRequest PullRequest
}

// PullRequest.
type PullRequest struct {
	Title       string
	Description string
	Changes     string
}

func ghExecContext(ctx context.Context, args ...string) (string, error) {
	stdout, stderr, err := gh.ExecContext(ctx, args...)
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func generatePrompt(ctx context.Context, instruction, pr string) (string, error) {
	title, err := ghExecContext(ctx, "pr", "view", pr, "--json", "title", "--jq", ".title")
	if err != nil {
		return "", err
	}
	description, err := ghExecContext(ctx, "pr", "view", pr, "--json", "body", "--jq", ".body")
	if err != nil {
		return "", err
	}
	changes, err := ghExecContext(ctx, "pr", "diff", pr)
	if err != nil {
		return "", err
	}
	prompt := Prompt{
		Instruction: instruction,
		PullRequest: PullRequest{
			Title:       title,
			Description: description,
			Changes:     changes,
		},
	}
	b, err := xml.Marshal(prompt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func converse(ctx context.Context, modelID, content string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}
	client := bedrockruntime.NewFromConfig(cfg)

	output, err := client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: aws.String(modelID),
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: content,
					},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	b := []byte{}
	if msg, ok := output.Output.(*types.ConverseOutputMemberMessage); ok {
		for _, content := range msg.Value.Content {
			if text, ok := content.(*types.ContentBlockMemberText); ok {
				b = append(b, []byte(text.Value)...)
			}
		}
	}
	return string(b), nil
}

func main() {
	var (
		instruction string
		modelID     string
	)

	log.SetFlags(0)

	command := &cobra.Command{
		Use:   "gh-pr-review-ai [<url>]",
		Short: "Review a pull request using a generative AI.",
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			prompt, err := generatePrompt(ctx, instruction, args[0])
			if err != nil {
				return err
			}
			output, err := converse(ctx, modelID, prompt)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		},
		Args:          cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	command.Flags().StringVar(&instruction, "instruction", "You are an IT expert. Review the pull request below. If what should be fixed or improved are found, tell me them specifically. Write your response in Japanese.", "The prompt that provides instructions to the model about the task it should perform.")
	command.Flags().StringVar(&modelID, "model-id", "anthropic.claude-3-5-sonnet-20240620-v1:0", "The model with which to run inference.")

	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
