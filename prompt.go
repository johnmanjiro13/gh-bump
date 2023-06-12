package bump

import "github.com/AlecAivazis/survey/v2"

type prompter struct{}

func newPrompter() Prompter {
	return &prompter{}
}

func (p *prompter) Input(question string, validator survey.Validator) (string, error) {
	var result string
	if err := survey.AskOne(&survey.Input{Message: question}, &result, survey.WithValidator(validator)); err != nil {
		return "", err
	}
	return result, nil
}

func (p *prompter) Select(question string, options []string) (string, error) {
	var result string
	if err := survey.AskOne(&survey.Select{Message: question, Options: options}, &result); err != nil {
		return "", err
	}
	return result, nil
}

func (p *prompter) Confirm(question string) (bool, error) {
	var result bool
	if err := survey.AskOne(&survey.Confirm{Message: question}, &result); err != nil {
		return false, err
	}
	return result, nil
}
