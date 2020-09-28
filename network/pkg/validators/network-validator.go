package validators

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"gopkg.in/go-playground/validator.v9"
)

type NetworkValidator struct {
	validate *validator.Validate
}

func NewNetworkValidator() *NetworkValidator {
	return &NetworkValidator{validate: validator.New()}
}

func (v *NetworkValidator) ValidateStruct(value interface{}) error {
	err := v.validate.Struct(value)
	// fmt.Println("validate result:", err)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e.ActualTag())       // alpha
			log.Println(e.Field())           // FirstName
			log.Println(e.Kind())            // string
			log.Println(e.Namespace())       // RegisterRequest.FirstName
			log.Println(e.Param())           //
			log.Println(e.StructField())     // FirstName
			log.Println(e.StructNamespace()) // RegisterRequest.FirstName
			log.Println(e.Type())            // string
			log.Println(e.Value())           // 123
		}
	}
	return err
}

func (v *NetworkValidator) ValidateUserId(id string) error {
	var s = struct {
		UserId string `validate:"len=24,alphanum"`
	}{id}
	return v.ValidateStruct(s)
}

func (v *NetworkValidator) ValidateId(id string) error {
	parts := strings.Split(id, "/")
	if len(parts) == 2 {
		err := v.ValidateKey(parts[1])
		if err == nil {
			return nil
		}
	}
	return errors.New("Invalid id")
}

func (v *NetworkValidator) ValidateKey(key string) error {
	if len(key) == 24 {
		return nil
	}
	return errors.New("Invalid key")
}

func (v *NetworkValidator) ValidateAddExperienceRequest(req *model.AddExperienceRequest) error {
	err := v.validate.Struct(req)
	if err != nil {
		return err
	}
	if req.WorkingCurrently {
		return nil
	}

	err = errors.New("Requestor Date should be greater then Requested Date")
	if req.ToYear < req.FromYear {
		return err
	}
	if req.ToYear == req.FromYear && req.ToMonth < req.FromMonth {
		return err
	}

	return nil
}
