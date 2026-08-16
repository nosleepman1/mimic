package domain

// FieldType type de champ
type FieldType string

const (
	TypeString   FieldType = "string"
	TypeInteger  FieldType = "integer"
	TypeFloat    FieldType = "float"
	TypeBoolean  FieldType = "boolean"
	TypeArray    FieldType = "array"
	TypeObject   FieldType = "object"
	TypeUUID     FieldType = "uuid"
	TypeEmail    FieldType = "email"
	TypeName     FieldType = "name"
	TypePhone    FieldType = "phone"
	TypePhoneSN  FieldType = "phone_sn"
	TypeXOF      FieldType = "currency_xof"
	TypeDateISO  FieldType = "date_iso"
	TypeAvatar   FieldType = "avatar_url"
	TypeSentence FieldType = "sentence"
)

// Schema modélise la structure de données d'un corps de requête ou réponse
type Schema struct {
	Type        FieldType          `json:"type" yaml:"type"`
	Format      string             `json:"format,omitempty" yaml:"format,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty" yaml:"items,omitempty"`
	Required    []string           `json:"required,omitempty" yaml:"required,omitempty"`
	Example     interface{}        `json:"example,omitempty" yaml:"example,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty" yaml:"enum,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
}
