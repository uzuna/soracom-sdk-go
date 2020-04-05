package soracom

import (
	"fmt"
	"net/http"
	"strconv"
)

// EventDateTimeConst is the time value that can be specified when execute a rule, enable, or disable
type EventDateTimeConst string

func (x EventDateTimeConst) String() string {
	return string(x)
}

const (
	// EventDateTimeImmediately is immediately
	EventDateTimeImmediately EventDateTimeConst = "IMMEDIATELY"

	// EventDateTimeAfterOneDay is one day (24 hours) later
	EventDateTimeAfterOneDay = "AFTER_ONE_DAY"

	// EventDateTimeBeginningOfNextDay is ...
	EventDateTimeBeginningOfNextDay = "BEGINNING_OF_NEXT_DAY"

	// EventDateTimeBeginningOfNextMonth is ...
	EventDateTimeBeginningOfNextMonth = "BEGINNING_OF_NEXT_MONTH"

	// EventDateTimeNever is ...
	EventDateTimeNever = "NEVER"
)

// EventStatus is status of EventHandler
type EventStatus string

// EventStatus is one of active or inactive
const (
	EventStatusActive   EventStatus = "active"
	EventStatusInactive EventStatus = "inactive"
)

func buildRuleConfig(evrtype EventHandlerRuleType, datetimeConst EventDateTimeConst, prop Properties) RuleConfig {
	prop["inactiveTimeoutDateConst"] = datetimeConst.String()
	return RuleConfig{
		Type:       evrtype,
		Properties: prop,
	}
}
func buildActionConfig(acttype EventHandlerActionType, datetimeConst EventDateTimeConst, prop Properties) ActionConfig {
	prop["executionDateTimeConst"] = datetimeConst.String()
	return ActionConfig{
		Type:       acttype,
		Properties: prop,
	}
}

func RuleDailyTraffic(mib uint64, datetimeConst EventDateTimeConst) RuleConfig {
	prop := Properties{
		"limitTotalTrafficMegaByte": strconv.FormatUint(mib, 10),
	}
	return buildRuleConfig(EventHandlerRuleTypeDailyTraffic, datetimeConst, prop)
}

func RuleMonthlyTraffic(mib uint64, datetimeConst EventDateTimeConst) RuleConfig {
	prop := Properties{
		"limitTotalTrafficMegaByte": strconv.FormatUint(mib, 10),
	}
	return buildRuleConfig(EventHandlerRuleTypeMonthlyTraffic, datetimeConst, prop)
}

func ActionActivate(datetimeConst EventDateTimeConst) ActionConfig {
	return buildActionConfig(EventHandlerActionTypeActivate, datetimeConst, Properties{})
}

func ActionDeactivate(datetimeConst EventDateTimeConst) ActionConfig {
	return buildActionConfig(EventHandlerActionTypeDeactivate, datetimeConst, Properties{})
}

type ActionWebhookProperty struct {
	URL         string
	Method      string
	ContentType string
	Body        string
}

func (p ActionWebhookProperty) Verify() error {

	switch p.Method {
	case http.MethodPost, http.MethodPut:
	default:
		if p.Body != "" {
			return fmt.Errorf("%s method does not use body field [%s]", p.Method, p.Body)
		}
	}
	return nil
}
func (p ActionWebhookProperty) toProperty() Properties {
	prop := Properties{
		"url":         p.URL,
		"httpMethod":  p.Method,
		"contentType": p.ContentType,
	}
	switch p.Method {
	case http.MethodPost, http.MethodPut:
		prop["body"] = p.Body
	}
	return prop
}

func ActionWebHook(datetimeConst EventDateTimeConst, hookprop ActionWebhookProperty) ActionConfig {
	return buildActionConfig(EventHandlerActionTypeExecuteWebRequest, datetimeConst, hookprop.toProperty())
}

func ActionChangeSpeed(datetimeConst EventDateTimeConst, s SpeedClass) ActionConfig {
	prop := Properties{"speedClass": s.String()}
	return buildActionConfig(EventHandlerActionTypeChangeSpeedClass, datetimeConst, prop)
}

type ActionSendEmailProperty struct {
	To      string
	Title   string
	Message string
}

func (p ActionSendEmailProperty) toProperty() Properties {
	return Properties{
		"to":      p.To,
		"title":   p.Title,
		"message": p.Message,
	}
}
func ActionSendEmail(datetimeConst EventDateTimeConst, mailprop ActionSendEmailProperty) ActionConfig {
	return buildActionConfig(EventHandlerActionTypeExecuteWebRequest, datetimeConst, mailprop.toProperty())
}
