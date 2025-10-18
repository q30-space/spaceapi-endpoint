package models

// SpaceAPI represents the SpaceAPI v15 structure
type SpaceAPI struct {
	APICompatibility []string         `json:"api_compatibility"`
	Space            string           `json:"space"`
	Logo             string           `json:"logo"`
	URL              string           `json:"url"`
	Location         *Location        `json:"location,omitempty"`
	Spacefed         *Spacefed        `json:"spacefed,omitempty"`
	Cam              []string         `json:"cam,omitempty"`
	State            *State           `json:"state,omitempty"`
	Events           []Event          `json:"events,omitempty"`
	Contact          Contact          `json:"contact"`
	Sensors          *Sensors         `json:"sensors,omitempty"`
	Feeds            *Feeds           `json:"feeds,omitempty"`
	Projects         []string         `json:"projects,omitempty"`
	Links            []Link           `json:"links,omitempty"`
	MembershipPlans  []MembershipPlan `json:"membership_plans,omitempty"`
	LinkedSpaces     []LinkedSpace    `json:"linked_spaces,omitempty"`
}

type Location struct {
	Address     string  `json:"address,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Timezone    string  `json:"timezone,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
	Hint        string  `json:"hint,omitempty"`
	Areas       []Area  `json:"areas,omitempty"`
}

type Area struct {
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	SquareMeters float64 `json:"square_meters"`
}

type Spacefed struct {
	Spacenet  bool `json:"spacenet"`
	Spacesaml bool `json:"spacesaml"`
}

type State struct {
	Open          *bool  `json:"open,omitempty"`
	Lastchange    int64  `json:"lastchange,omitempty"`
	TriggerPerson string `json:"trigger_person,omitempty"`
	Message       string `json:"message,omitempty"`
	Icon          *Icon  `json:"icon,omitempty"`
}

type Icon struct {
	Open   string `json:"open"`
	Closed string `json:"closed"`
}

type Event struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Extra     string `json:"extra,omitempty"`
}

type Contact struct {
	Phone      string      `json:"phone,omitempty"`
	Sip        string      `json:"sip,omitempty"`
	Keymasters []Keymaster `json:"keymasters,omitempty"`
	IRC        string      `json:"irc,omitempty"`
	Twitter    string      `json:"twitter,omitempty"`
	Mastodon   string      `json:"mastodon,omitempty"`
	Facebook   string      `json:"facebook,omitempty"`
	Identica   string      `json:"identica,omitempty"`
	Foursquare string      `json:"foursquare,omitempty"`
	Email      string      `json:"email,omitempty"`
	ML         string      `json:"ml,omitempty"`
	XMPP       string      `json:"xmpp,omitempty"`
	IssueMail  string      `json:"issue_mail,omitempty"`
}

type Keymaster struct {
	Name     string `json:"name,omitempty"`
	IRCNick  string `json:"irc_nick,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
	XMPP     string `json:"xmpp,omitempty"`
	Mastodon string `json:"mastodon,omitempty"`
	Matrix   string `json:"matrix,omitempty"`
}

type Sensors struct {
	Temperature        []SensorValue `json:"temperature,omitempty"`
	DoorLocked         []SensorValue `json:"door_locked,omitempty"`
	Barometer          []SensorValue `json:"barometer,omitempty"`
	Radiation          []SensorValue `json:"radiation,omitempty"`
	Humidity           []SensorValue `json:"humidity,omitempty"`
	BeverageSupply     []SensorValue `json:"beverage_supply,omitempty"`
	PowerConsumption   []SensorValue `json:"power_consumption,omitempty"`
	Wind               []SensorValue `json:"wind,omitempty"`
	NetworkConnections []SensorValue `json:"network_connections,omitempty"`
	AccountBalance     []SensorValue `json:"account_balance,omitempty"`
	TotalMemberCount   []SensorValue `json:"total_member_count,omitempty"`
	PeopleNowPresent   []SensorValue `json:"people_now_present,omitempty"`
	NetworkTraffic     []SensorValue `json:"network_traffic,omitempty"`
}

type SensorValue struct {
	Value       interface{} `json:"value"`
	Unit        string      `json:"unit,omitempty"`
	Location    string      `json:"location,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Lastchange  int64       `json:"lastchange,omitempty"`
}

type Feeds struct {
	Blog     *Feed `json:"blog,omitempty"`
	Wiki     *Feed `json:"wiki,omitempty"`
	Calendar *Feed `json:"calendar,omitempty"`
	Flickr   *Feed `json:"flickr,omitempty"`
}

type Feed struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url"`
}

type Link struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

type MembershipPlan struct {
	Name            string  `json:"name"`
	Value           float64 `json:"value"`
	Currency        string  `json:"currency"`
	BillingInterval string  `json:"billing_interval"`
	Description     string  `json:"description,omitempty"`
}

type LinkedSpace struct {
	Endpoint string `json:"endpoint,omitempty"`
	Website  string `json:"website,omitempty"`
}

// Helper function to create a bool pointer
func BoolPtr(b bool) *bool {
	return &b
}
