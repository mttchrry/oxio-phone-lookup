package phoneNumbers

import (
	"context"
	"fmt"
	"log"

	"regexp"

	"github.com/biter777/countries"
	"github.com/mttchrry/oxio-phone-lookup/pkg/phoneNumbers/model"
	"github.com/mttchrry/oxio-phone-lookup/pkg/utils/errors"
)

const (
	ErrInvalidNumberValue = errors.Error("invalid formatted phone number")
	ErrDifferentCountryCodes = errors.Error("different filter country codes type")
	ErrInvalidCountryCode = errors.Error("invalid country codes")
)

type PhoneNumbers struct {
	numberToCountryCode map[string][]string
	countryCodeNums map[string]string
	nonNumericRegex *regexp.Regexp
	phoneRegex *regexp.Regexp
}

func New() (*PhoneNumbers) {
	cs := countries.AllInfo()
	p := PhoneNumbers{}
	p.numberToCountryCode = map[string][]string{}
	p.countryCodeNums = map[string]string{}

	for _, s := range cs {
		cc := fmt.Sprintf("%v", s.Code.CallCodes()[0])
		p.numberToCountryCode[cc] = append(p.numberToCountryCode[cc], s.Alpha2)
		p.countryCodeNums[s.Alpha2] = cc
	}

	p.nonNumericRegex = regexp.MustCompile(`[^0-9]+`)
	p.phoneRegex = regexp.MustCompile(`^[\+]?[(]?[0-9]{0,4}[)]?[-\s]?[(]?[0-9]{3}[)]?[-\s]?[0-9]{3}[-\s]?[0-9]{4}$`)
	return &p

}

func (p *PhoneNumbers) Parse(ctx context.Context, number string, countryCode string) (*model.PhoneNumber, error) {
	simpleN, err := p.simplifyString(ctx, number)
	if err != nil {
		return nil, err
	}

	countryNumber := simpleN[:len(simpleN)-10]
	areaCode := simpleN[len(simpleN)-10:len(simpleN)-7]
	local := simpleN[len(simpleN)-7:]

	if countryNumber != ""{
		countryNumber = "+"+countryNumber
	}
	if countryCode != "" {
		if _, ok := p.countryCodeNums[countryCode]; !ok {
			return nil, fmt.Errorf("error: country code %v doesnt exist", countryCode)
		}
		if len(simpleN) != 10 {
			// verify number matches CC 
			fmt.Printf("\n ^^^ converted number to CC : %v\n", p.numberToCountryCode[countryNumber])
			found := false
			for _, cc := range p.numberToCountryCode[countryNumber] {
				fmt.Printf("\n *** %v = %v ?", countryCode, cc)
				if countryCode == cc {
					found = true
				}
			}
			if !found {
				return nil, fmt.Errorf("error: country number %v doesnt match given code: %v", countryNumber, countryCode)
			}
		}
	} else {
		if len(simpleN) == 10 {
			return nil, fmt.Errorf("error: need country identifier in number %v or a country code", number)
		}
		countryCode = p.numberToCountryCode[countryNumber][0] // pick one at random
	}

	if len(simpleN) == 10{
		countryNumber = p.countryCodeNums[countryCode]
		simpleN = countryNumber + simpleN
	} else {
		simpleN = "+"+simpleN
	}

	pN := model.PhoneNumber{
		PhoneNumber: simpleN,
		CountryCode: countryCode,
		AreaCode: areaCode,
		Local: local,
	}

	return &pN, nil
}

func (p *PhoneNumbers) simplifyString(ctx context.Context, number string) (string, error) {
	if !p.phoneRegex.MatchString(number){
		err := fmt.Errorf("error: invalid formatted phone number: %v", number)
		log.Printf(err.Error())
		return "", err
	}

	simpleN := p.nonNumericRegex.ReplaceAllString(number, "")
	if len(simpleN) < 10 {
		// this shouldn't happen, just checking
		return "", fmt.Errorf("internal error on regex, resulting in %v becoming sn %v", number, simpleN)
	}
	return simpleN, nil
}