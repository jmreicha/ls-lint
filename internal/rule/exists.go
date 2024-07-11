package rule

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
)

type Exists struct {
	name      string
	exclusive bool
	min       uint16
	max       uint16
	count     uint16
	*sync.RWMutex
}

func (rule *Exists) Init() Rule {
	rule.name = "exists"
	rule.exclusive = true
	rule.count = 0
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Exists) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

// 0 = regex pattern
func (rule *Exists) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	// exists
	if len(params) == 0 {
		rule.min = 1
		rule.max = math.MaxInt16
		return nil
	}

	// exists:
	if params[0] == "" {
		return fmt.Errorf("exists value is empty")
	}

	// exists:1
	var split = strings.Split(params[0], "-")
	if len(split) == 1 {
		var value int64
		var err error

		if value, err = strconv.ParseInt(params[0], 10, 16); err != nil {
			return err.(*strconv.NumError).Err
		}

		rule.min = uint16(value)
		rule.max = uint16(value)
		return nil
	}

	// exists:1-4
	var minValue int64
	var maxValue int64
	var err error

	if minValue, err = strconv.ParseInt(split[0], 10, 16); err != nil {
		return err.(*strconv.NumError).Err
	}

	if maxValue, err = strconv.ParseInt(split[1], 10, 16); err != nil {
		return err.(*strconv.NumError).Err
	}

	rule.min = uint16(minValue)
	rule.max = uint16(maxValue)
	return nil
}

func (rule *Exists) GetParameters() []string {
	if rule.getMin() == rule.getMax() {
		return []string{fmt.Sprintf("%d", rule.getMin())}
	}

	return []string{fmt.Sprintf("%d-%d", rule.getMin(), rule.getMax())}
}

func (rule *Exists) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

func (rule *Exists) Validate(value string, fail bool) (bool, error) {
	if !fail {
		rule.incrementCount()
		return true, nil
	}

	return rule.getCount() >= rule.getMin() && rule.getCount() <= rule.getMax(), nil
}

func (rule *Exists) getMin() uint16 {
	rule.RLock()
	defer rule.RUnlock()

	return rule.min
}

func (rule *Exists) getMax() uint16 {
	rule.RLock()
	defer rule.RUnlock()

	return rule.max
}

func (rule *Exists) getCount() uint16 {
	rule.RLock()
	defer rule.RUnlock()

	return rule.count
}

func (rule *Exists) incrementCount() {
	rule.Lock()
	defer rule.Unlock()

	rule.count++
}

func (rule *Exists) GetErrorMessage() string {
	if rule.getMin() == rule.getMax() {
		return fmt.Sprintf("%s:%d", rule.GetName(), rule.getMin())
	}

	return fmt.Sprintf("%s:%d-%d", rule.GetName(), rule.getMin(), rule.getMax())
}
