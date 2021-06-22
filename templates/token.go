package templates

import (
	"fmt"
	"strings"

	"github.com/eqlabs/flow-wallet-service/flow_helpers"
	"github.com/eqlabs/flow-wallet-service/templates/template_strings"
	"github.com/iancoleman/strcase"
	"github.com/onflow/flow-go-sdk"
)

type Token struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func EnabledTokens() []Token {
	return parseConfig().enabledTokens
}

func EnabledTokenAddresses() map[string]string {
	return parseConfig().enabledTokenAddresses
}

func EnabledTokenNames() []string {
	enabled := EnabledTokens()
	keys := make([]string, len(enabled))
	for i, k := range enabled {
		keys[i] = k.CanonName()
	}
	return keys
}

func NewToken(name string) (Token, error) {
	t := Token{Name: name} // So we can compare the CanonName
	address, ok := EnabledTokenAddresses()[t.CanonName()]
	if !ok {
		return Token{}, fmt.Errorf("token %s not enabled", t.CanonName())
	}
	return Token{Name: name, Address: address}, nil
}

func (t *Token) CanonName() string {
	return t.ParseName()[0]
}

func TokenFromEvent(e flow.Event, chainId flow.ChainID) (Token, error) {
	// Example event:
	// A.0ae53cb6e3f42a79.FlowToken.TokensDeposited
	ss := strings.Split(e.Type, ".")
	a, err := flow_helpers.ValidateAddress(ss[1], chainId)
	if err != nil {
		return Token{}, err
	}
	return Token{Name: ss[2], Address: a}, nil
}

func (t *Token) IsEnabled() bool {
	for _, e := range EnabledTokens() {
		if t.Name == e.Name && t.Address == e.Address {
			return true
		}
	}
	return false
}

func (t *Token) ParseName() [3]string {
	// TODO: how to handle these kind of cases?
	if strings.ToLower(t.Name) == "fusd" {
		return [3]string{
			"FUSD", "FUSD", "fusd",
		}
	}

	return [3]string{
		strcase.ToCamel(t.Name),
		strcase.ToScreamingSnake(t.Name),
		strcase.ToLowerCamel(t.Name),
	}
}

func fungibleTemplateCode(tmpl_str string, token Token) string {
	p := token.ParseName()
	camel := p[0]
	snake := p[1]
	lower := p[2]

	r := strings.NewReplacer(
		"TokenName", camel,
		"TOKEN_NAME", snake,
		"tokenName", lower,
	)

	tmpl_str = r.Replace(tmpl_str)

	// Replace token address
	if token.Address != "" {
		r := strings.NewReplacer(
			fmt.Sprintf("%s_ADDRESS", snake), token.Address,
		)
		tmpl_str = r.Replace(tmpl_str)
	}

	return Code(&Template{Source: tmpl_str})
}

func FungibleTransferCode(token Token) string {
	return fungibleTemplateCode(
		template_strings.GenericFungibleTransfer,
		token,
	)
}

func FungibleSetupCode(token Token) string {
	return fungibleTemplateCode(
		template_strings.GenericFungibleSetup,
		token,
	)
}

func FungibleBalanceCode(token Token) string {
	return fungibleTemplateCode(
		template_strings.GenericFungibleBalance,
		token,
	)
}