package audio

type TriggerKey  uint8
type ResourceKey uint8

// Most loosely typed mess ever. Very practical.
type AutomationPanel struct {
	triggers []func() error
	resources []any
}

func NewAutomationPanel() *AutomationPanel {
	return &AutomationPanel{
		triggers: make([]func() error, 0, 4),
		resources: make([]any, 0, 4),
	}
}

func (self *AutomationPanel) RegisterTrigger(trigger func() error) TriggerKey {
	key := TriggerKey(len(self.triggers))
	self.triggers = append(self.triggers, trigger)
	return key
}

func (self *AutomationPanel) Trigger(key TriggerKey) error {
	return self.triggers[key]()
}

func (self *AutomationPanel) StoreResource(resource any) ResourceKey {
	key := ResourceKey(len(self.resources))
	self.resources = append(self.resources, resource)
	return key
}

func (self *AutomationPanel) GetResource(key ResourceKey) any {
	return self.resources[key]
}


