package validator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trusch/jwtd/jwt"
)

//Validator validates requests
type Validator struct {
	cert interface{}
}

// New returns a new validator
func New(cert interface{}) (*Validator, error) {
	return &Validator{cert}, nil
}

// Validate validates a requests
func (validator *Validator) Validate(r *http.Request, service string, required map[string]string) error {
	if len(required) == 0 {
		log.Print("no labels required for this reques, now forwarding...")
		return nil
	}
	claims, err := jwt.GetClaimsFromRequest(r, validator.cert)

	if err != nil {
		log.Printf("can not get claims from request (%v), return 401", err)
		return err
	}

	if srv, ok := claims["service"].(string); ok {
		if srv != service {
			log.Printf("service in claim doesn't match, return 401")
			return errors.New("service mismatch")
		}
	} else {
		log.Printf("service in claim not valid, return 401")
		return errors.New("no valid service field in token")
	}

	if err = validator.validateNbf(claims); err != nil {
		log.Printf("NBF check failed: %v, return 401", err)
		return err
	}

	if err = validator.validateExp(claims); err != nil {
		log.Printf("EXP check failed: %v, return 401", err)
		return err
	}

	if err = validator.validateLabels(claims, validator.resolveVariables(required, mux.Vars(r))); err != nil {
		log.Printf("claims do not have the required labels: %v, return 401", err)
		return err
	}
	log.Printf("all checks passed, forwarding...")
	return nil
}

func (validator *Validator) resolveVariables(reqs map[string]string, vars map[string]string) map[string]string {
	res := make(map[string]string)
	for key, value := range reqs {
		if len(key) > 0 && key[0] == '$' {
			varName := key[1:]
			if val, ok := vars[varName]; ok {
				key = val
			}
		}
		if len(value) > 0 && value[0] == '$' {
			varName := value[1:]
			if val, ok := vars[varName]; ok {
				value = val
			}
		}
		res[key] = value
	}
	return res
}

func (validator *Validator) validateNbf(claims map[string]interface{}) error {
	if nbfStr, ok := claims["nbf"].(string); ok {
		nbf := time.Time{}
		err := nbf.UnmarshalText([]byte(nbfStr))
		if err != nil {
			return fmt.Errorf("failed parsing nbf string: %v", nbfStr)
		}
		if time.Now().Before(nbf) {
			return errors.New("nbf is in the future")
		}
		return nil
	}
	return errors.New("no nbf given")
}

func (validator *Validator) validateExp(claims map[string]interface{}) error {
	if expStr, ok := claims["exp"].(string); ok {
		exp := time.Time{}
		err := exp.UnmarshalText([]byte(expStr))
		if err != nil {
			return fmt.Errorf("failed parsing exp string: %v", expStr)
		}
		if !time.Now().Before(exp) {
			return errors.New("exp is in the past")
		}
		return nil
	}
	return errors.New("no exp given")
}

func (validator *Validator) validateLabels(claims map[string]interface{}, required map[string]string) error {
	if labels, ok := claims["labels"].(map[string]interface{}); ok {
		for rKey, rValue := range required {
			if uValue, ok := labels[rKey]; !ok || uValue != rValue {
				return fmt.Errorf("can not validate label %v:%v", rKey, rValue)
			}
		}
		return nil
	}
	return errors.New("no labels given")
}
