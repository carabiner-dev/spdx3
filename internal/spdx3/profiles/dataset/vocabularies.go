package dataset

// ConfidentialityLevelType vocabulary defines valid confidentiality levels using Traffic Light Protocol
type ConfidentialityLevelType string

const (
	// ConfidentialityLevelTypeAmber indicates limited disclosure to participants' organizations
	ConfidentialityLevelTypeAmber ConfidentialityLevelType = "amber"
	// ConfidentialityLevelTypeClear indicates unlimited disclosure
	ConfidentialityLevelTypeClear ConfidentialityLevelType = "clear"
	// ConfidentialityLevelTypeGreen indicates limited disclosure within community
	ConfidentialityLevelTypeGreen ConfidentialityLevelType = "green"
	// ConfidentialityLevelTypeRed indicates no disclosure beyond participants
	ConfidentialityLevelTypeRed ConfidentialityLevelType = "red"
)

// DatasetAvailabilityType vocabulary defines valid dataset availability types
type DatasetAvailabilityType string

const (
	// DatasetAvailabilityTypeClickthrough indicates dataset requires accepting terms through web interface
	DatasetAvailabilityTypeClickthrough DatasetAvailabilityType = "clickthrough"
	// DatasetAvailabilityTypeDirect indicates dataset is directly downloadable without restrictions
	DatasetAvailabilityTypeDirect DatasetAvailabilityType = "direct"
	// DatasetAvailabilityTypeQuery indicates dataset requires query to access portions
	DatasetAvailabilityTypeQuery DatasetAvailabilityType = "query"
	// DatasetAvailabilityTypeRegistration indicates dataset requires registration before access
	DatasetAvailabilityTypeRegistration DatasetAvailabilityType = "registration"
	// DatasetAvailabilityTypeScrapingScript indicates dataset requires script to collect data
	DatasetAvailabilityTypeScrapingScript DatasetAvailabilityType = "scrapingScript"
)

// DatasetType vocabulary defines valid dataset types
type DatasetType string

const (
	// DatasetTypeAudio indicates dataset contains audio data
	DatasetTypeAudio DatasetType = "audio"
	// DatasetTypeCodeSnippet indicates dataset contains code snippets
	DatasetTypeCodeSnippet DatasetType = "codeSnippet"
	// DatasetTypeConfiguration indicates dataset contains configuration data
	DatasetTypeConfiguration DatasetType = "configuration"
	// DatasetTypeGraph indicates dataset is in graph format
	DatasetTypeGraph DatasetType = "graph"
	// DatasetTypeImage indicates dataset contains image data
	DatasetTypeImage DatasetType = "image"
	// DatasetTypeIndex indicates dataset is an index
	DatasetTypeIndex DatasetType = "index"
	// DatasetTypeNoAssertion makes no assertion about dataset type
	DatasetTypeNoAssertion DatasetType = "noAssertion"
	// DatasetTypeNumeric indicates dataset contains numeric data
	DatasetTypeNumeric DatasetType = "numeric"
	// DatasetTypeOther indicates dataset is of another type
	DatasetTypeOther DatasetType = "other"
	// DatasetTypeSensor indicates dataset contains sensor data
	DatasetTypeSensor DatasetType = "sensor"
	// DatasetTypeStructured indicates dataset has structured format
	DatasetTypeStructured DatasetType = "structured"
	// DatasetTypeTable indicates dataset is in tabular format
	DatasetTypeTable DatasetType = "table"
	// DatasetTypeText indicates dataset contains text data
	DatasetTypeText DatasetType = "text"
	// DatasetTypeVideo indicates dataset contains video data
	DatasetTypeVideo DatasetType = "video"
)
